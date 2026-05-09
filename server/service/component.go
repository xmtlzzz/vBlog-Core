package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

// ComponentService handles component CRUD operations.
type ComponentService struct {
	DB *gorm.DB
}

// NewComponentService creates a new ComponentService.
func NewComponentService(db *gorm.DB) *ComponentService {
	return &ComponentService{DB: db}
}

// List returns all components.
func (s *ComponentService) List() ([]model.Component, error) {
	var components []model.Component
	err := s.DB.Find(&components).Error
	return components, err
}

// Create creates a new component with auto-set origin "uploaded".
func (s *ComponentService) Create(c *model.Component) error {
	c.Origin = "uploaded"
	return s.DB.Create(c).Error
}

// Update updates an existing component.
func (s *ComponentService) Update(c *model.Component) error {
	return s.DB.Save(c).Error
}

// Delete deletes a component by ID.
func (s *ComponentService) Delete(id uint) error {
	return s.DB.Delete(&model.Component{}, id).Error
}

// Toggle switches a component between active and inactive status.
func (s *ComponentService) Toggle(id uint) error {
	var comp model.Component
	if err := s.DB.First(&comp, id).Error; err != nil {
		return err
	}
	if comp.Status == "active" {
		comp.Status = "inactive"
	} else {
		comp.Status = "active"
	}
	return s.DB.Save(&comp).Error
}
