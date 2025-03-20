package repositories

import (
    "Loop_backend/internal/models"
)

// TagRepository defines interfaces for tag operations in PostgreSQL
type TagRepository interface {
    // Basic CRUD
    CreateTag(tag *models.Tag) error
    GetTagByID(id string) (*models.Tag, error)
    GetTagByName(name string) (*models.Tag, error)
    UpdateTag(tag *models.Tag) error
    DeleteTag(id int) error

    // Vector operations
    StoreTagVector(tagID int, vector []float64) error
    FindSimilarTags(vector []float64, threshold float64, limit int) ([]*models.Tag, error)
}
