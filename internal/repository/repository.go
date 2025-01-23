package repository

import (
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepo{db: db}
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepo{db: db}
}
