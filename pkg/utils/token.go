package utils

import (
	"strconv"
	"strings"
	"todo/pkg/errors"
)

// ParseUserIDFromToken 从token中解析用户ID
func ParseUserIDFromToken(token string) (uint, error) {
	// 这里简化处理，实际应该使用JWT解析
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, errors.ErrInvalidToken
	}

	userID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, errors.ErrInvalidToken
	}

	return uint(userID), nil
}
