package lambda

import (
	"context"
	"zeroslope/database"

	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Create(ctx context.Context, db *gorm.DB, entity database.SampleEntity) (*database.SampleEntity, error) {
	result := db.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func Read(ctx context.Context, db *gorm.DB, id int) (*database.SampleEntity, error) {
	var entity database.SampleEntity
	result := db.First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func Update(ctx context.Context, db *gorm.DB, id int, newEntity database.SampleEntity) (*database.SampleEntity, error) {
	var entity database.SampleEntity
	db.First(&entity, id)
	entity.Name = newEntity.Name
	entity.Description = newEntity.Description
	db.Save(&entity)
	return &entity, nil
}

func Delete(ctx context.Context, db *gorm.DB, id int) error {
	var entity database.SampleEntity
	result := db.First(&entity, id)
	if result.Error != nil {
		return result.Error
	}
	db.Delete(&entity)
	return nil
}
