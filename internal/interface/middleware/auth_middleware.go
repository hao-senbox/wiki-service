package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"wiki-service/pkg/gateway"
	libs_constant "wiki-service/pkg/libs/constant"
	libs_helper "wiki-service/pkg/libs/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Secured(userGw gateway.UserGateway) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")

		// app language header
		appLanguage := libs_helper.ParseAppLanguage(c.Get("X-App-Language"), 1)
		c.Set("X-App-Language", strconv.Itoa(int(appLanguage)))
		c.Locals(libs_constant.AppLanguage.String(), appLanguage)
		ctx := context.WithValue(c.Context(), libs_constant.AppLanguage, appLanguage)

		// --- Authorization ---
		if len(authorizationHeader) == 0 {
			return c.SendStatus(http.StatusForbidden)
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]

		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// --- UserID ---
			if userId, ok := claims[libs_constant.UserID.String()].(string); ok {
				c.Locals(libs_constant.UserID.String(), userId)
				ctx = context.WithValue(ctx, libs_constant.UserID, userId)
			}

			// --- UserName ---
			if userName, ok := claims[libs_constant.UserName.String()].(string); ok {
				c.Locals(libs_constant.UserName.String(), userName)
				ctx = context.WithValue(ctx, libs_constant.UserName, userName)
			}

			// --- Roles ---
			if userRoles, ok := claims[libs_constant.UserRoles.String()].(string); ok {
				c.Locals(libs_constant.UserRoles.String(), userRoles)
				ctx = context.WithValue(ctx, libs_constant.UserRoles, userRoles)
			}
		}

		// --- Token ---
		c.Locals(libs_constant.Token.String(), tokenString)
		ctx = context.WithValue(ctx, libs_constant.Token, tokenString)

		// --- Call user-service to get current user ---
		currentUser, err := userGw.GetCurrentUser(
			context.WithValue(ctx, libs_constant.CurrentUserKey, tokenString),
		)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// --- Set currentUser vào context ---
		c.Locals(string(libs_constant.CurrentUserKey), currentUser)
		ctx = context.WithValue(ctx, libs_constant.CurrentUserKey, currentUser)

		// Store the context for use in handlers
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rolesAny := c.Locals(libs_constant.UserRoles.String())
		if rolesAny == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Roles not found"})
		}

		rolesStr, ok := rolesAny.(string)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid roles format"})
		}

		// ví dụ roles: "SuperAdmin, Teacher"
		roles := strings.Split(rolesStr, ",")
		isAdmin := false
		for _, role := range roles {
			if strings.TrimSpace(role) == "SuperAdmin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}

		return c.Next()
	}
}
