package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bsonger/devflow-app-service/pkg/model"
	"github.com/bsonger/devflow-app-service/pkg/store"
	"github.com/bsonger/devflow-service-common/loggingx"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ProjectService = NewProjectService()

type ProjectListFilter struct {
	IncludeDeleted bool
	Name           string
	Key            string
	Namespace      string
	Owner          string
}

type projectService struct{}

func NewProjectService() *projectService {
	return &projectService{}
}

func (s *projectService) Create(ctx context.Context, project *model.Project) (uuid.UUID, error) {
	log := loggingx.LoggerWithContext(ctx).With(zap.String("operation", "create_project"))

	project.ApplyDefaults()
	labels, err := marshalLabels(project.Labels)
	if err != nil {
		log.Error("marshal project labels failed", zap.Error(err))
		return uuid.Nil, err
	}

	_, err = store.DB().ExecContext(ctx, `
		insert into projects (
			id, key, name, description, namespace, owner, labels, created_at, updated_at, deleted_at
		) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, project.ID, project.Key, project.Name, project.Description, project.Namespace, project.Owner, labels, project.CreatedAt, project.UpdatedAt, project.DeletedAt)
	if err != nil {
		log.Error("create project failed", zap.Error(err))
		return uuid.Nil, err
	}

	log.Info("project created", zap.String("project_id", project.GetID().String()), zap.String("project_key", project.Key))
	return project.GetID(), nil
}

func (s *projectService) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "get_project"),
		zap.String("project_id", id.String()),
	)

	project, err := scanProject(store.DB().QueryRowContext(ctx, `
		select id, key, name, description, namespace, owner, labels, created_at, updated_at, deleted_at
		from projects
		where id = $1 and deleted_at is null
	`, id))
	if err != nil {
		log.Error("get project failed", zap.Error(err))
		return nil, err
	}

	log.Debug("project fetched", zap.String("project_key", project.Key))
	return project, nil
}

func (s *projectService) Update(ctx context.Context, project *model.Project) error {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "update_project"),
		zap.String("project_id", project.GetID().String()),
	)

	current, err := s.Get(ctx, project.GetID())
	if err != nil {
		log.Error("load project failed", zap.Error(err))
		return err
	}

	project.CreatedAt = current.CreatedAt
	project.DeletedAt = current.DeletedAt
	project.WithUpdateDefault()
	project.ApplyDefaults()

	labels, err := marshalLabels(project.Labels)
	if err != nil {
		return err
	}

	result, err := store.DB().ExecContext(ctx, `
		update projects
		set key=$2, name=$3, description=$4, namespace=$5, owner=$6, labels=$7, updated_at=$8, deleted_at=$9
		where id = $1 and deleted_at is null
	`, project.ID, project.Key, project.Name, project.Description, project.Namespace, project.Owner, labels, project.UpdatedAt, project.DeletedAt)
	if err != nil {
		log.Error("update project failed", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Info("project updated", zap.String("project_key", project.Key))
	return nil
}

func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "delete_project"),
		zap.String("project_id", id.String()),
	)

	now := time.Now()
	result, err := store.DB().ExecContext(ctx, `
		update projects
		set deleted_at=$2, updated_at=$2
		where id = $1 and deleted_at is null
	`, id, now)
	if err != nil {
		log.Error("delete project failed", zap.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Info("project deleted")
	return nil
}

func (s *projectService) List(ctx context.Context, filter ProjectListFilter) ([]model.Project, error) {
	log := loggingx.LoggerWithContext(ctx).With(
		zap.String("operation", "list_projects"),
		zap.Any("filter", filter),
	)

	query := `
		select id, key, name, description, namespace, owner, labels, created_at, updated_at, deleted_at
		from projects
	`
	clauses := make([]string, 0, 5)
	args := make([]any, 0, 5)

	if !filter.IncludeDeleted {
		clauses = append(clauses, "deleted_at is null")
	}
	if filter.Name != "" {
		args = append(args, filter.Name)
		clauses = append(clauses, placeholderClause("name", len(args)))
	}
	if filter.Key != "" {
		args = append(args, filter.Key)
		clauses = append(clauses, placeholderClause("key", len(args)))
	}
	if filter.Namespace != "" {
		args = append(args, filter.Namespace)
		clauses = append(clauses, placeholderClause("namespace", len(args)))
	}
	if filter.Owner != "" {
		args = append(args, filter.Owner)
		clauses = append(clauses, placeholderClause("owner", len(args)))
	}
	if len(clauses) > 0 {
		query += " where " + strings.Join(clauses, " and ")
	}
	query += " order by created_at desc"

	rows, err := store.DB().QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("list projects failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Debug("projects listed", zap.Int("count", len(projects)))
	return projects, nil
}

func (s *projectService) ListApplications(ctx context.Context, projectID uuid.UUID) ([]model.Application, error) {
	if _, err := s.Get(ctx, projectID); err != nil {
		return nil, err
	}

	return ApplicationService.List(ctx, ApplicationListFilter{ProjectID: &projectID})
}

func scanProject(scanner interface {
	Scan(dest ...any) error
}) (*model.Project, error) {
	var (
		project     model.Project
		labelsBytes []byte
		deletedAt   sql.NullTime
	)

	if err := scanner.Scan(
		&project.ID,
		&project.Key,
		&project.Name,
		&project.Description,
		&project.Namespace,
		&project.Owner,
		&labelsBytes,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
	); err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		project.DeletedAt = &deletedAt.Time
	}
	if len(labelsBytes) > 0 {
		if err := json.Unmarshal(labelsBytes, &project.Labels); err != nil {
			return nil, err
		}
	}

	return &project, nil
}

func marshalLabels(labels map[string]string) ([]byte, error) {
	if labels == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(labels)
}

func placeholderClause(column string, position int) string {
	return column + " = $" + strconv.Itoa(position)
}
