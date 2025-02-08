package handlers

import (
	"github.com/gin-gonic/gin"
	"todo/pkg/response"
)

// HealthCheck 健康检查处理器
// @Summary 服务健康检查
// @Description 用于监控系统运行状态的健康检查接口
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]string "返回ok表示服务正常"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	response.Success(c, map[string]string{
		"status": "ok",
	})
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string `json:"status"`     // 服务状态
	Timestamp string `json:"timestamp"`  // 时间戳
	Version   string `json:"version"`    // 服务版本
}