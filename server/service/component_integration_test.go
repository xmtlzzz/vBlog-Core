package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestComponentService_Create(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewComponentService(db)

	comp := &model.Component{
		Name:        "TestWidget",
		Description: "A test component",
		Version:     "1.0.0",
		Code:        "<div>Test</div>",
	}
	err := svc.Create(comp)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if comp.Origin != "uploaded" {
		t.Errorf("expected origin 'uploaded', got '%s'", comp.Origin)
	}
	if comp.ID == 0 {
		t.Fatal("expected component ID to be set")
	}
	db.Unscoped().Delete(comp)
}

func TestComponentService_List(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewComponentService(db)

	svc.Create(&model.Component{Name: "ListComp1", Category: "test"})
	svc.Create(&model.Component{Name: "ListComp2", Category: "test"})

	comps, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(comps) < 2 {
		t.Errorf("expected >= 2 components, got %d", len(comps))
	}

	db.Unscoped().Where("category = ?", "test").Delete(&model.Component{})
}

func TestComponentService_Toggle(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewComponentService(db)

	comp := &model.Component{Name: "ToggleComp", Status: "active"}
	svc.Create(comp)

	err := svc.Toggle(comp.ID)
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}

	var found model.Component
	db.First(&found, comp.ID)
	if found.Status != "inactive" {
		t.Errorf("expected status 'inactive', got '%s'", found.Status)
	}

	// Toggle back
	svc.Toggle(comp.ID)
	db.First(&found, comp.ID)
	if found.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", found.Status)
	}

	db.Unscoped().Delete(comp)
}

func TestComponentService_Delete(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewComponentService(db)

	comp := &model.Component{Name: "DeleteComp"}
	svc.Create(comp)

	err := svc.Delete(comp.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	var count int64
	db.Model(&model.Component{}).Where("id = ?", comp.ID).Count(&count)
	if count != 0 {
		t.Error("expected component to be deleted")
	}
	db.Unscoped().Delete(comp)
}
