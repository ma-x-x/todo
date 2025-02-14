package app

import (
	"todo/internal/repository"
	"todo/internal/service"
)

// initServices 初始化所有服务
func (a *App) initServices() error {
	// 初始化仓储层
	repos := repository.NewRepositories(a.db, a.rdb)

	// 初始化服务层
	authService := service.NewAuthService(repos.User, repos.Auth, &a.cfg.JWT)
	todoService := service.NewTodoService(repos.Todo)
	categoryService := service.NewCategoryService(repos.Category)
	reminderService := service.NewReminderService(repos.Reminder, repos.Todo)

	// 创建服务集合
	a.services = service.NewServiceCollection(
		authService,
		todoService,
		categoryService,
		reminderService,
		a.db,
		a.rdb,
	)

	return nil
}
