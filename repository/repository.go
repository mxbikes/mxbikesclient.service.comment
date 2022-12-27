package repository

import (
	"github.com/mxbikes/mxbikesclient.service.comment/models"
	"gorm.io/gorm"
)

type ModRepository interface {
	SearchByModID(modID string) ([]*models.Comment, error)
	Save(comment *models.Comment) error
	Delete(id string) error
	Migrate() error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(c *gorm.DB) *postgresRepository {
	return &postgresRepository{db: c}
}

func (p *postgresRepository) SearchByModID(modID string) ([]*models.Comment, error) {
	var l []*models.Comment
	err := p.db.Where(`mod_id = ?`, modID).Find(&l).Error
	return l, err
}

func (p *postgresRepository) Save(comment *models.Comment) error {
	return p.db.Save(comment).Error
}

func (p *postgresRepository) Delete(id string) error {
	return p.db.Delete(&models.Comment{ID: id}).Error
}

func (p *postgresRepository) Migrate() error {
	return p.db.AutoMigrate(&models.Comment{})
}
