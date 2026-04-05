package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestApplicationContract(t *testing.T) {
	typ := reflect.TypeOf(Application{})

	for _, field := range []string{"ProjectID", "Name", "RepoAddress", "ActiveManifestID", "Labels"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("Application missing field %s", field)
		}
		switch field {
		case "ProjectID":
			if f.Type != reflect.TypeOf(uuid.UUID{}) {
				t.Fatalf("Application.ProjectID type = %v, want uuid.UUID", f.Type)
			}
		case "RepoAddress":
			if got := f.Tag.Get("db"); got != "repo_address" {
				t.Fatalf("Application.RepoAddress db tag = %q, want %q", got, "repo_address")
			}
		}
	}

	for _, removed := range []string{"ProjectName", "ActiveManifestName", "ConfigMaps", "Service", "Internet", "Envs", "Replica", "Type", "Status"} {
		if _, ok := typ.FieldByName(removed); ok {
			t.Fatalf("Application should not expose legacy field %s", removed)
		}
	}
}

func TestServiceResourceContract(t *testing.T) {
	typ := reflect.TypeOf(ServiceResource{})
	for _, field := range []string{"ApplicationID", "Name", "Exposure", "Ports", "Status"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("ServiceResource missing field %s", field)
		}
		if field == "ApplicationID" && f.Type != reflect.TypeOf(uuid.UUID{}) {
			t.Fatalf("ServiceResource.ApplicationID type = %v, want uuid.UUID", f.Type)
		}
	}
}

func TestProjectContractAfterAudit(t *testing.T) {
	typ := reflect.TypeOf(Project{})
	if _, ok := typ.FieldByName("Status"); ok {
		t.Fatal("Project should not expose Status")
	}
}

func TestEnvironmentContract(t *testing.T) {
	typ := reflect.TypeOf(Environment{})
	for _, field := range []string{"Key", "Name", "Cluster", "Namespace", "Labels"} {
		if _, ok := typ.FieldByName(field); !ok {
			t.Fatalf("Environment missing field %s", field)
		}
	}
}

func TestBaseModelWithCreateDefault(t *testing.T) {
	var base BaseModel
	base.WithCreateDefault()

	if base.ID == uuid.Nil {
		t.Fatal("BaseModel.WithCreateDefault should assign a UUID")
	}
	if base.CreatedAt.IsZero() {
		t.Fatal("BaseModel.WithCreateDefault should set CreatedAt")
	}
	if base.UpdatedAt.IsZero() {
		t.Fatal("BaseModel.WithCreateDefault should set UpdatedAt")
	}
}
