.PHONY: dev build test clean deps swagger

# 变量定义
APP_NAME=todo-api
MAIN_FILE=cmd/server/main.go
BUILD_DIR=build

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