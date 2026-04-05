package service

import (
	"context"
	"errors"
	"time"

	"github.com/bsonger/devflow-app-service/pkg/model"
	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var ApplicationService = NewApplicationService()

var ErrManifestNotForApplication = errors.New("manifest does not belong to application")
var ErrProjectReferenceNotFound = errors.New("project reference not found")

type applicationService struct{}

func NewApplicationService() *applicationService {
	return &applicationService{}
}

// Create 创建 Application
func (s *applicationService) Create(ctx context.Context, app *model.Application) (uuid.UUID, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "create_application"),
	)

	if err := s.syncProjectReference(ctx, app); err != nil {
		log.Error("resolve project reference failed", zap.Error(err))
		return uuid.Nil, err
	}
	doc, err := applicationToDoc(app)
	if err != nil {
		log.Error("prepare application doc failed", zap.Error(err))
		return uuid.Nil, err
	}

	if err := mongo.Repo.Create(ctx, doc); err != nil {
		log.Error("create application failed", zap.Error(err))
		return uuid.Nil, err
	}
	app.ID = bridgeObjectIDToUUID(doc.ID)

	log.Info("application created", zap.String("application_id", app.GetID().String()))
	return app.GetID(), nil
}

// Get 根据 ID 查询 Application
func (s *applicationService) Get(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return nil, err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "get_application"),
		zap.String("application_id", id.String()),
	)

	doc := &applicationDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, oid); err != nil {
		log.Error("get application failed", zap.Error(err))
		return nil, err
	}
	if doc.DeletedAt != nil {
		log.Warn("application already deleted")
		return nil, mongoDriver.ErrNoDocuments
	}

	app := applicationFromDoc(doc)
	log.Debug("application fetched", zap.String("application_name", app.Name))
	return &app, nil
}

// Update 更新 Application
func (s *applicationService) Update(ctx context.Context, app *model.Application) error {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "update_application"),
		zap.String("application_id", app.GetID().String()),
	)
	appOID, err := bridgeUUIDToObjectID(app.GetID())
	if err != nil {
		return err
	}

	currentDoc := &applicationDoc{}
	if err := mongo.Repo.FindByID(ctx, currentDoc, appOID); err != nil {
		log.Error("load application failed", zap.Error(err))
		return err
	}
	if currentDoc.DeletedAt != nil {
		log.Warn("update skipped for deleted application")
		return mongoDriver.ErrNoDocuments
	}

	current := applicationFromDoc(currentDoc)
	app.CreatedAt = current.CreatedAt
	app.DeletedAt = current.DeletedAt
	app.WithUpdateDefault()

	if err := s.syncProjectReference(ctx, app); err != nil {
		log.Error("resolve project reference failed", zap.Error(err))
		return err
	}

	doc, err := applicationToDoc(app)
	if err != nil {
		return err
	}
	doc.ID = appOID

	if err := mongo.Repo.Update(ctx, doc); err != nil {
		log.Error("update application failed", zap.Error(err))
		return err
	}

	log.Debug("application updated", zap.String("application_name", app.Name))
	return nil
}

// Delete 删除 Application
func (s *applicationService) Delete(ctx context.Context, id uuid.UUID) error {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "delete_application"),
		zap.String("application_id", id.String()),
	)

	now := time.Now()
	update := primitive.M{
		"$set": primitive.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	if err := mongo.Repo.UpdateByID(ctx, &applicationDoc{}, oid, update); err != nil {
		log.Error("delete application failed", zap.Error(err))
		return err
	}

	log.Info("application deleted")
	return nil
}

// UpdateActiveManifest updates the application active manifest reference.
func (s *applicationService) UpdateActiveManifest(ctx context.Context, appID, manifestID uuid.UUID) error {
	appOID, err := bridgeUUIDToObjectID(appID)
	if err != nil {
		return err
	}
	manifestOID, err := bridgeUUIDToObjectID(manifestID)
	if err != nil {
		return err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "update_application_active_manifest"),
		zap.String("application_id", appID.String()),
		zap.String("manifest_id", manifestID.String()),
	)

	doc := &applicationDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, appOID); err != nil {
		log.Error("get application failed", zap.Error(err))
		return err
	}
	if doc.DeletedAt != nil {
		log.Warn("application already deleted")
		return mongoDriver.ErrNoDocuments
	}

	manifest := &manifestDoc{}
	if err := mongo.Repo.FindByID(ctx, manifest, manifestOID); err != nil {
		log.Error("get manifest failed", zap.Error(err))
		return err
	}
	if manifest.DeletedAt != nil {
		log.Warn("manifest already deleted")
		return mongoDriver.ErrNoDocuments
	}
	if manifest.ApplicationID != appOID {
		log.Warn("manifest does not belong to application")
		return ErrManifestNotForApplication
	}

	update := primitive.M{
		"$set": primitive.M{
			"active_manifest_id":   manifestOID,
			"active_manifest_name": manifest.Name,
			"updated_at":           time.Now(),
		},
	}

	if err := mongo.Repo.UpdateByID(ctx, &applicationDoc{}, appOID, update); err != nil {
		log.Error("update active manifest failed", zap.Error(err))
		return err
	}

	log.Info("active manifest updated", zap.String("active_manifest_name", manifest.Name))
	return nil
}

// List 查询 Application 列表
func (s *applicationService) List(ctx context.Context, filter primitive.M) ([]model.Application, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "list_applications"),
		zap.Any("filter", filter),
	)

	var docs []applicationDoc
	if err := mongo.Repo.List(ctx, &applicationDoc{}, filter, &docs); err != nil {
		log.Error("list applications failed", zap.Error(err))
		return nil, err
	}

	apps := make([]model.Application, 0, len(docs))
	for i := range docs {
		apps = append(apps, applicationFromDoc(&docs[i]))
	}

	log.Debug("applications listed", zap.Int("count", len(apps)))
	return apps, nil
}

func (s *applicationService) syncProjectReference(ctx context.Context, app *model.Application) error {
	if app.ProjectID == uuid.Nil {
		return nil
	}

	_, err := ProjectService.Get(ctx, app.ProjectID)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return ErrProjectReferenceNotFound
		}
		return err
	}
	return nil
}
