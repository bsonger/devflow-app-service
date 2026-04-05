package service

import (
	"errors"

	"github.com/bsonger/devflow-app-service/pkg/model"
	commonmodel "github.com/bsonger/devflow-common/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errInvalidUUIDBridge = errors.New("invalid bridged uuid")
	bridgeUUIDPrefix     = [4]byte{'d', 'f', 'l', 'w'}
)

type projectDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	Name        string              `bson:"name"`
	Key         string              `bson:"key"`
	Description string              `bson:"description,omitempty"`
	Namespace   string              `bson:"namespace,omitempty"`
	Owner       string              `bson:"owner,omitempty"`
	Labels      map[string]string   `bson:"labels,omitempty"`
	Status      model.ProjectStatus `bson:"status"`
}

func (projectDoc) CollectionName() string { return "projects" }

type applicationDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	Name               string              `bson:"name"`
	ProjectID          *primitive.ObjectID `bson:"project_id,omitempty"`
	ProjectName        string              `bson:"project_name,omitempty"`
	RepoURL            string              `bson:"repo_url,omitempty"`
	RepoAddress        string              `bson:"repo_address,omitempty"`
	ActiveManifestName string              `bson:"active_manifest_name,omitempty"`
	ActiveManifestID   *primitive.ObjectID `bson:"active_manifest_id,omitempty"`
	Replica            *int32              `bson:"replica,omitempty"`
	Type               model.ReleaseType   `bson:"type"`
	Status             string              `bson:"status"`
}

func (applicationDoc) CollectionName() string { return "applications" }

type manifestDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	ApplicationID primitive.ObjectID `bson:"application_id"`
	Name          string             `bson:"name"`
}

func (manifestDoc) CollectionName() string { return "manifests" }

func bridgeObjectIDToUUID(id primitive.ObjectID) uuid.UUID {
	var raw [16]byte
	copy(raw[:4], bridgeUUIDPrefix[:])
	copy(raw[4:], id[:])
	return uuid.UUID(raw)
}

func bridgeUUIDToObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	raw := [16]byte(id)
	if raw[0] != bridgeUUIDPrefix[0] || raw[1] != bridgeUUIDPrefix[1] || raw[2] != bridgeUUIDPrefix[2] || raw[3] != bridgeUUIDPrefix[3] {
		return primitive.NilObjectID, errInvalidUUIDBridge
	}

	var oid primitive.ObjectID
	copy(oid[:], raw[4:])
	return oid, nil
}

func BridgeUUIDToObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	return bridgeUUIDToObjectID(id)
}

func projectFromDoc(doc *projectDoc) model.Project {
	return model.Project{
		BaseModel: model.BaseModel{
			ID:        bridgeObjectIDToUUID(doc.ID),
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			DeletedAt: doc.DeletedAt,
		},
		Name:        doc.Name,
		Key:         doc.Key,
		Description: doc.Description,
		Namespace:   doc.Namespace,
		Owner:       doc.Owner,
		Labels:      doc.Labels,
		Status:      doc.Status,
	}
}

func projectToDoc(project *model.Project) (*projectDoc, error) {
	id := primitive.NewObjectID()
	if project.ID != uuid.Nil {
		bridgedID, err := bridgeUUIDToObjectID(project.ID)
		if err == nil {
			id = bridgedID
		}
	}

	return &projectDoc{
		BaseModel: commonmodel.BaseModel{
			ID:        id,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
			DeletedAt: project.DeletedAt,
		},
		Name:        project.Name,
		Key:         project.Key,
		Description: project.Description,
		Namespace:   project.Namespace,
		Owner:       project.Owner,
		Labels:      project.Labels,
		Status:      project.Status,
	}, nil
}

func applicationFromDoc(doc *applicationDoc) model.Application {
	app := model.Application{
		BaseModel: model.BaseModel{
			ID:        bridgeObjectIDToUUID(doc.ID),
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			DeletedAt: doc.DeletedAt,
		},
		Name:        doc.Name,
		RepoAddress: doc.RepoAddress,
		Replica:     doc.Replica,
		Type:        doc.Type,
		Status:      doc.Status,
	}
	if app.RepoAddress == "" {
		app.RepoAddress = doc.RepoURL
	}
	if doc.ProjectID != nil && !doc.ProjectID.IsZero() {
		projectID := bridgeObjectIDToUUID(*doc.ProjectID)
		app.ProjectID = projectID
	}
	if doc.ActiveManifestID != nil && !doc.ActiveManifestID.IsZero() {
		manifestID := bridgeObjectIDToUUID(*doc.ActiveManifestID)
		app.ActiveManifestID = &manifestID
	}
	return app
}

func applicationToDoc(app *model.Application) (*applicationDoc, error) {
	id := primitive.NewObjectID()
	if app.ID != uuid.Nil {
		bridgedID, err := bridgeUUIDToObjectID(app.ID)
		if err == nil {
			id = bridgedID
		}
	}
	projectID, err := modelUUIDPtrToObjectIDPtr(&app.ProjectID)
	if err != nil {
		return nil, err
	}
	activeManifestID, err := modelUUIDPtrToObjectIDPtr(app.ActiveManifestID)
	if err != nil {
		return nil, err
	}

	return &applicationDoc{
		BaseModel: commonmodel.BaseModel{
			ID:        id,
			CreatedAt: app.CreatedAt,
			UpdatedAt: app.UpdatedAt,
			DeletedAt: app.DeletedAt,
		},
		Name:             app.Name,
		ProjectID:        projectID,
		RepoURL:          app.RepoAddress,
		RepoAddress:      app.RepoAddress,
		ActiveManifestID: activeManifestID,
		Replica:          app.Replica,
		Type:             app.Type,
		Status:           app.Status,
	}, nil
}

func modelUUIDPtrToObjectIDPtr(id *uuid.UUID) (*primitive.ObjectID, error) {
	if id == nil || *id == uuid.Nil {
		return nil, nil
	}
	oid, err := bridgeUUIDToObjectID(*id)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}
