package utils

import (
	"strconv"
	"strings"
	"todo/pkg/errors"
)

// Token 格式说明:
// token的格式为 "prefix:userID"
// 例如: "todo:123" 表示用户ID为123的token
//
// prefix: 固定为"todo"，用于标识token来源
// userID: 用户的唯一标识符，为无符号整数
//
// 注意：当前实现为简化版本，生产环境建议使用JWT等标准token格式
// JWT优势：
// 1. 可以包含更多用户信息
// 2. 支持token过期机制
// 3. 提供签名验证，防止篡改
// 4. 支持跨服务认证

// ParseUserIDFromToken 从token字符串中解析出用户ID
// 参数:
//   - token: 需要解析的token字符串
//
// 返回值:
//   - uint: 解析出的用户ID
//   - error: 解析过程中的错误，可能的错误类型:
//   - ErrInvalidToken: token格式不正确或无法解析userID
func ParseUserIDFromToken(token string) (uint, error) {
	// 将token按":"分割成两部分
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, errors.ErrInvalidToken
	}

	// 将userID部分转换为uint类型
	userID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, errors.ErrInvalidToken
	}

	return uint(userID), nil
}
