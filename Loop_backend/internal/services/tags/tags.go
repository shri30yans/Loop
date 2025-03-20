package services

import (
    "Loop_backend/internal/models"
tags "Loop_backend/internal/repositories/tags"
)

type TagService interface {
    // Tag management
    CreateTag(tag *models.Tag) error
    GetTagByID(id string) (*models.Tag, error)
    GetTagByName(name string) (*models.Tag, error)
    UpdateTag(tag *models.Tag) error
    
    // Tag relationships
    CreateTagRelationship(tag1, tag2 string, strength float64) error
    GetRelatedTags(tagName string, minStrength float64) ([]*models.TagRelationship, error)
    
    // Project tags
    GetProjectTags(projectID string) ([]string, error)
    GetTagsByProject(projectID string) ([]*models.Tag, error)
    
    // User expertise
    SetUserTagExpertise(userID, tagName, level string, years int) error
    GetUserExpertise(userID string) (map[string]string, error)
    GetTagExperts(tagName string) ([]models.User, error)
}

type tagService struct {
    tagRepo    tags.TagRepository
}

func NewTagService(tagRepo tags.TagRepository) TagService {
    return &tagService{
        tagRepo: tagRepo,
    }
}

func (s *tagService) CreateTag(tag *models.Tag) error {
    return s.tagRepo.CreateTag(tag)
}

func (s *tagService) GetTagByID(id string) (*models.Tag, error) {
    return s.tagRepo.GetTagByID(id)
}

func (s *tagService) GetTagByName(name string) (*models.Tag, error) {
    return s.tagRepo.GetTagByName(name)
}

func (s *tagService) UpdateTag(tag *models.Tag) error {
    return s.tagRepo.UpdateTag(tag)
}

func (s *tagService) CreateTagRelationship(tag1, tag2 string, strength float64) error {
    return nil // Placeholder for relationship logic
}

func (s *tagService) GetRelatedTags(tagName string, minStrength float64) ([]*models.TagRelationship, error) {
    return nil, nil // Placeholder for related tags logic
}

func (s *tagService) GetProjectTags(projectID string) ([]string, error) {
    return nil, nil // Placeholder for project tags logic
}

func (s *tagService) GetTagsByProject(projectID string) ([]*models.Tag, error) {
    return nil, nil // Placeholder for tags by project logic
}

func (s *tagService) SetUserTagExpertise(userID, tagName, level string, years int) error {
    return nil // Placeholder for user expertise logic
}

func (s *tagService) GetUserExpertise(userID string) (map[string]string, error) {
    return nil, nil // Placeholder for user expertise logic
}

func (s *tagService) GetTagExperts(tagName string) ([]models.User, error) {
    return nil, nil // Placeholder for tag experts logic
}
