.PHONY: dev build test clean deps swagger init-db init-db-with-credentials

# 变量定义
APP_NAME=todo-api
MAIN_FILE=cmd/server/main.go
BUILD_DIR=build

# 数据库配置（可以从环境变量读取）
DB_USER ?= root
DB_PASS ?= root
DB_HOST ?= localhost
DB_PORT ?= 3306

dev:
	@mkdir -p logs tmp
	@air -c .air.toml

# 构建项目
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR) logs tmp
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# 运行测试
test:
	@go test -v ./...

test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# 清理构建文件
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@rm -f logs/*.log
	@rm -f tmp/*

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go get github.com/rs/zerolog
	@go get github.com/redis/go-redis/v9
	@go get gorm.io/gorm@latest
	@go get gorm.io/plugin/dbresolver@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/air-verse/air@latest

swagger:
	@echo "Generating swagger docs..."
	@swag init -g cmd/server/main.go --parseDependency --parseInternal --parseDepth 1

docker-build:
	@docker build -t $(APP_NAME) .

docker-run:
	@docker run -p 8080:8080 $(APP_NAME)

# 添加 init-db 命令
init-db:
	@echo "Initializing database..."
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) < scripts/init.sql
	@echo "Database initialization completed"

# 添加一个更安全的版本，允许指定用户名和密码
init-db-with-credentials:
	@echo "Initializing database..."
	@read -p "Enter MySQL username: " username; \
	read -s -p "Enter MySQL password: " password; \
	echo ""; \
	mysql -u $$username -p$$password < scripts/init.sql
	@echo "Database initialization completed" 