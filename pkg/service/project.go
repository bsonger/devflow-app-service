package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-app-service/pkg/model"
	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var ProjectService = NewProjectService()

type projectService struct{}

func NewProjectService() *projectService {
	return &projectService{}
}

func (s *projectService) Create(ctx context.Context, project *model.Project) (uuid.UUID, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "create_project"),
	)

	project.ApplyDefaults()
	doc, err := projectToDoc(project)
	if err != nil {
		log.Error("prepare project doc failed", zap.Error(err))
		return uuid.Nil, err
	}
	if err := mongo.Repo.Create(ctx, doc); err != nil {
		log.Error("create project failed", zap.Error(err))
		return uuid.Nil, err
	}
	project.ID = bridgeObjectIDToUUID(doc.ID)

	log.Info("project created", zap.String("project_id", project.GetID().String()), zap.String("project_key", project.Key))
	return project.GetID(), nil
}

func (s *projectService) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return nil, err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "get_project"),
		zap.String("project_id", id.String()),
	)

	doc := &projectDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, oid); err != nil {
		log.Error("get project failed", zap.Error(err))
		return nil, err
	}
	if doc.DeletedAt != nil {
		log.Warn("project already deleted")
		return nil, mongoDriver.ErrNoDocuments
	}

	project := projectFromDoc(doc)
	log.Debug("project fetched", zap.String("project_key", project.Key))
	return &project, nil
}

func (s *projectService) Update(ctx context.Context, project *model.Project) error {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "update_project"),
		zap.String("project_id", project.GetID().String()),
	)
	projectOID, err := bridgeUUIDToObjectID(project.GetID())
	if err != nil {
		return err
	}

	currentDoc := &projectDoc{}
	if err := mongo.Repo.FindByID(ctx, currentDoc, projectOID); err != nil {
		log.Error("load project failed", zap.Error(err))
		return err
	}
	if currentDoc.DeletedAt != nil {
		log.Warn("update skipped for deleted project")
		return mongoDriver.ErrNoDocuments
	}

	current := projectFromDoc(currentDoc)
	project.CreatedAt = current.CreatedAt
	project.DeletedAt = current.DeletedAt
	project.WithUpdateDefault()
	project.ApplyDefaults()

	doc, err := projectToDoc(project)
	if err != nil {
		return err
	}
	doc.ID = projectOID

	if err := mongo.Repo.Update(ctx, doc); err != nil {
		log.Error("update project failed", zap.Error(err))
		return err
	}

	if current.Name != project.Name {
		if err := s.syncApplicationProjectNames(ctx, project.GetID(), project.Name); err != nil {
			log.Error("sync project name to applications failed", zap.Error(err))
			return err
		}
	}

	log.Info("project updated", zap.String("project_key", project.Key))
	return nil
}

func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "delete_project"),
		zap.String("project_id", id.String()),
	)

	now := time.Now()
	update := primitive.M{
		"$set": primitive.M{
			"deleted_at": now,
			"updated_at": now,
			"status":     model.ProjectArchived,
		},
	}

	if err := mongo.Repo.UpdateByID(ctx, &projectDoc{}, oid, update); err != nil {
		log.Error("delete project failed", zap.Error(err))
		return err
	}

	log.Info("project deleted")
	return nil
}

func (s *projectService) List(ctx context.Context, filter primitive.M) ([]model.Project, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "list_projects"),
		zap.Any("filter", filter),
	)

	var docs []projectDoc
	if err := mongo.Repo.List(ctx, &projectDoc{}, filter, &docs); err != nil {
		log.Error("list projects failed", zap.Error(err))
		return nil, err
	}

	projects := make([]model.Project, 0, len(docs))
	for i := range docs {
		projects = append(projects, projectFromDoc(&docs[i]))
	}

	log.Debug("projects listed", zap.Int("count", len(projects)))
	return projects, nil
}

func (s *projectService) ListApplications(ctx context.Context, projectID uuid.UUID) ([]model.Application, error) {
	if _, err := s.Get(ctx, projectID); err != nil {
		return nil, err
	}
	projectOID, err := bridgeUUIDToObjectID(projectID)
	if err != nil {
		return nil, err
	}

	filter := primitive.M{
		"deleted_at": primitive.M{"$exists": false},
		"project_id": projectOID,
	}

	return ApplicationService.List(ctx, filter)
}

func (s *projectService) syncApplicationProjectNames(ctx context.Context, projectID uuid.UUID, projectName string) error {
	projectOID, err := bridgeUUIDToObjectID(projectID)
	if err != nil {
		return err
	}
	return mongo.Repo.UpdateMany(ctx, &applicationDoc{}, bson.M{"project_id": projectOID}, bson.M{
		"$set": bson.M{
			"project_name": projectName,
			"updated_at":   time.Now(),
		},
	})
}
