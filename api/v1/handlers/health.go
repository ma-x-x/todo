package handlers

import (
	"github.com/gin-gonic/gin"
)

// Health 健康检查处理器
// @Summary 服务健康检查
// @Description 用于监控系统运行状态的健康检查接口
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]string "返回ok表示服务正常"
// @Router /health [get]
func Health(c *gin.Context) {
	// 返回服务健康状态
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
