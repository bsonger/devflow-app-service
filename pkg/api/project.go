package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/bsonger/devflow-app-service/pkg/model"
	"github.com/bsonger/devflow-app-service/pkg/service"
	"github.com/bsonger/devflow-service-common/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ProjectRouteApi = NewProjectHandler()

type projectService interface {
	Create(ctx context.Context, project *model.Project) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter service.ProjectListFilter) ([]model.Project, error)
	ListApplications(ctx context.Context, projectID uuid.UUID) ([]model.Application, error)
}

type ProjectHandler struct {
	svc projectService
}

type CreateProjectRequest struct {
	Name        string            `json:"name"`
	Key         string            `json:"key"`
	Description string            `json:"description,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type UpdateProjectRequest = CreateProjectRequest

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{svc: service.ProjectService}
}

// Create
// @Summary 创建项目
// @Description 创建一个新的项目
// @Tags Project
// @Accept json
// @Produce json
// @Param data body api.CreateProjectRequest true "Project Data"
// @Success 201 {object} httpx.DataResponse[model.Project]
// @Router /api/v1/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	project := &model.Project{
		Name:        req.Name,
		Key:         req.Key,
		Description: req.Description,
		Namespace:   req.Namespace,
		Owner:       req.Owner,
		Labels:      req.Labels,
	}
	project.WithCreateDefault()
	project.ApplyDefaults()

	_, err := h.svc.Create(c.Request.Context(), project)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteData(c, http.StatusCreated, project)
}

// Get
// @Summary 获取项目
// @Tags Project
// @Param id path string true "Project ID"
// @Success 200 {object} httpx.DataResponse[model.Project]
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	project, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteData(c, http.StatusOK, project)
}

// Update
// @Summary 更新项目
// @Tags Project
// @Param id path string true "Project ID"
// @Param data body api.UpdateProjectRequest true "Project Data"
// @Success 204
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	project := model.Project{
		Name:        req.Name,
		Key:         req.Key,
		Description: req.Description,
		Namespace:   req.Namespace,
		Owner:       req.Owner,
		Labels:      req.Labels,
	}
	project.SetID(id)
	project.ApplyDefaults()

	if err := h.svc.Update(c.Request.Context(), &project); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteNoContent(c)
}

// Delete
// @Summary 删除项目
// @Tags Project
// @Param id path string true "Project ID"
// @Success 204
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	httpx.WriteNoContent(c)
}

// List
// @Summary 获取项目列表
// @Tags Project
// @Success 200 {object} httpx.ListResponse[model.Project]
// @Router /api/v1/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	filter := service.ProjectListFilter{
		IncludeDeleted: httpx.IncludeDeleted(c),
		Name:           c.Query("name"),
		Key:            c.Query("key"),
		Namespace:      c.Query("namespace"),
		Owner:          c.Query("owner"),
	}

	projects, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	total := len(projects)
	projects = httpx.PaginateSlice(projects, paging)
	httpx.WriteList(c, http.StatusOK, projects, paging, total)
}

// ListApplications
// @Summary 获取项目下的应用列表
// @Tags Project
// @Param id path string true "Project ID"
// @Success 200 {object} httpx.ListResponse[model.Application]
// @Router /api/v1/projects/{id}/applications [get]
func (h *ProjectHandler) ListApplications(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", "invalid id", nil)
		return
	}

	applications, err := h.svc.ListApplications(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(c, http.StatusNotFound, "not_found", "not found", nil)
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "internal", err.Error(), nil)
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid_argument", err.Error(), nil)
		return
	}

	total := len(applications)
	applications = httpx.PaginateSlice(applications, paging)
	httpx.WriteList(c, http.StatusOK, applications, paging, total)
}
