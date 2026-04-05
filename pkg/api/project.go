package api

import (
	"net/http"

	"github.com/bsonger/devflow-app-service/pkg/model"
	"github.com/bsonger/devflow-app-service/pkg/service"
	"github.com/bsonger/devflow-service-common/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ProjectRouteApi = NewProjectHandler()

type ProjectHandler struct{}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

// Create
// @Summary 创建项目
// @Description 创建一个新的项目
// @Tags Project
// @Accept json
// @Produce json
// @Param data body model.Project true "Project Data"
// @Success 200 {object} httpx.CreateResponse
// @Router /api/v1/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var project *model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.WithCreateDefault()
	project.ApplyDefaults()

	id, err := service.ProjectService.Create(c.Request.Context(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpx.CreateResponse{ID: id.String()})
}

// Get
// @Summary 获取项目
// @Tags Project
// @Param id path string true "Project ID"
// @Success 200 {object} model.Project
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	project, err := service.ProjectService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Update
// @Summary 更新项目
// @Tags Project
// @Param id path string true "Project ID"
// @Param data body model.Project true "Project Data"
// @Success 200 {object} map[string]string
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.SetID(id)
	project.ApplyDefaults()

	if err := service.ProjectService.Update(c.Request.Context(), &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Delete
// @Summary 删除项目
// @Tags Project
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := service.ProjectService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// List
// @Summary 获取项目列表
// @Tags Project
// @Success 200 {array} model.Project
// @Router /api/v1/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	filter := service.ProjectListFilter{
		IncludeDeleted: httpx.IncludeDeleted(c),
		Name:           c.Query("name"),
		Key:            c.Query("key"),
		Namespace:      c.Query("namespace"),
		Owner:          c.Query("owner"),
	}

	projects, err := service.ProjectService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := len(projects)
	projects = httpx.PaginateSlice(projects, paging)
	httpx.SetPaginationHeaders(c, total, paging)

	c.JSON(http.StatusOK, projects)
}

// ListApplications
// @Summary 获取项目下的应用列表
// @Tags Project
// @Param id path string true "Project ID"
// @Success 200 {array} model.Application
// @Router /api/v1/projects/{id}/applications [get]
func (h *ProjectHandler) ListApplications(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	applications, err := service.ProjectService.ListApplications(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paging, err := httpx.ParsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := len(applications)
	applications = httpx.PaginateSlice(applications, paging)
	httpx.SetPaginationHeaders(c, total, paging)

	c.JSON(http.StatusOK, applications)
}
