package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

// SettingService handles site settings operations.
type SettingService struct {
	DB *gorm.DB
}

// NewSettingService creates a new SettingService.
func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{DB: db}
}

// DefaultSettings returns the default site settings.
func DefaultSettings() map[string]string {
	return map[string]string{
		"site_title":       "vBlog",
		"site_description": "A lightweight blog for geeks",
		"site_url":         "",
		"posts_per_page":   "10",
		"theme":            "default",
		"footer_text":      "Powered by vBlog Core",
		"grpc_api_key":     "",
		"grpc_port":        "50051",
	}
}

// Get returns the value for a given setting key, or empty string if not found.
func (s *SettingService) Get(key string) (string, error) {
	var setting model.Setting
	if err := s.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return setting.Value, nil
}

// GetAll returns all settings as a key-value map.
func (s *SettingService) GetAll() (map[string]string, error) {
	var settings []model.Setting
	if err := s.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

// Set upserts a single setting key-value pair.
func (s *SettingService) Set(key, value string) error {
	return s.DB.Save(&model.Setting{Key: key, Value: value}).Error
}

// Save upserts all key-value pairs in the settings map.
func (s *SettingService) Save(settings map[string]string) error {
	for key, value := range settings {
		setting := model.Setting{Key: key, Value: value}
		if err := s.DB.Save(&setting).Error; err != nil {
			return err
		}
	}
	return nil
}

// Reset restores settings to defaults.
func (s *SettingService) Reset() error {
	return s.Save(DefaultSettings())
}
