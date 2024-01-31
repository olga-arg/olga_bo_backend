package db

import (
	"commons/domain"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type CategoryStorage struct {
	Db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryStorage {
	db.AutoMigrate(&domain.Category{})
	return &CategoryStorage{
		Db: db,
	}
}

func (r *CategoryStorage) Save(category *domain.Category) error {
	db := r.Db.AutoMigrate(&domain.Category{})
	err := db.Save(category).Error
	if err != nil {
		fmt.Println("Error saving category: ", err)
		return err
	}
	return nil
}

func (r *CategoryStorage) GetCategories(companyId string) (domain.Categories, error) {
	var categories domain.Categories
	err := r.Db.Where("company_id = ?", companyId).Find(&categories).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No categories found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting categories:", err)
		return nil, err
	}

	return categories, nil
}
