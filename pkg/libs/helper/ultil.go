package libs_helper

import (
	"context"
	"strconv"
	"strings"

	libs_constant "wiki-service/pkg/libs/constant"
)

func ParseAppLanguage(header string, defaultVal uint) uint {
	header = strings.TrimSpace(strings.Trim(header, "\""))
	if val, err := strconv.Atoi(header); err == nil {
		return uint(val)
	}
	return defaultVal
}

func GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)

	if lang, ok := ctx.Value(libs_constant.AppLanguage).(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	return headers
}

func GetAppLanguage(ctx context.Context, defaultVal uint) uint {
	if lang, ok := ctx.Value(libs_constant.AppLanguage).(uint); ok {
		return lang
	}
	return defaultVal
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(libs_constant.UserID).(string); ok {
		return userID
	}
	return ""
}
