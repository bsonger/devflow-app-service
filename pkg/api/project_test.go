package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsonger/devflow-app-service/pkg/app"
	"github.com/bsonger/devflow-app-service/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubProjectService struct {
	createFn           func(context.Context, *domain.Project) (uuid.UUID, error)
	getFn              func(context.Context, uuid.UUID) (*domain.Project, error)
	updateFn           func(context.Context, *domain.Project) error
	deleteFn           func(context.Context, uuid.UUID) error
	listFn             func(context.Context, app.ProjectListFilter) ([]domain.Project, error)
	listApplicationsFn func(context.Context, uuid.UUID) ([]domain.Application, error)
}

func (s stubProjectService) Create(ctx context.Context, project *domain.Project) (uuid.UUID, error) {
	return s.createFn(ctx, project)
}
func (s stubProjectService) Get(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.getFn(ctx, id)
}
func (s stubProjectService) Update(ctx context.Context, project *domain.Project) error {
	return s.updateFn(ctx, project)
}
func (s stubProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteFn(ctx, id)
}
func (s stubProjectService) List(ctx context.Context, filter app.ProjectListFilter) ([]domain.Project, error) {
	return s.listFn(ctx, filter)
}
func (s stubProjectService) ListApplications(ctx context.Context, projectID uuid.UUID) ([]domain.Application, error) {
	return s.listApplicationsFn(ctx, projectID)
}

func TestCreateProjectReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ProjectHandler{svc: stubProjectService{createFn: func(_ context.Context, project *domain.Project) (uuid.UUID, error) { return project.GetID(), nil }}}

	r := gin.New()
	r.POST("/api/v1/projects", handler.Create)

	body := bytes.NewBufferString(`{"name":"alpha","description":"platform project","labels":[{"key":"team","value":"platform"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("got %d want %d", rec.Code, http.StatusCreated)
	}

	var payload struct {
		Data domain.Project `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if payload.Data.Name != "alpha" || payload.Data.Description != "platform project" || len(payload.Data.Labels) != 1 {
		t.Fatalf("unexpected payload: %#v", payload.Data)
	}
}

func TestListProjectsReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ProjectHandler{svc: stubProjectService{listFn: func(_ context.Context, filter app.ProjectListFilter) ([]domain.Project, error) {
		if filter.Name != "" {
			t.Fatalf("unexpected filter: %#v", filter)
		}
		return []domain.Project{{Name: "alpha", Description: "platform project"}}, nil
	}}}

	r := gin.New()
	r.GET("/api/v1/projects", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?page=1&page_size=20", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data       []domain.Project `json:"data"`
		Pagination struct {
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if len(payload.Data) != 1 || payload.Pagination.Total != 1 {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestDeleteProjectNotFoundReturnsErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	handler := &ProjectHandler{svc: stubProjectService{deleteFn: func(_ context.Context, _ uuid.UUID) error { return sql.ErrNoRows }}}

	r := gin.New()
	r.DELETE("/api/v1/projects/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/projects/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d want %d", rec.Code, http.StatusNotFound)
	}
}
