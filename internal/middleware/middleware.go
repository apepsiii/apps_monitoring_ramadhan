package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
)

func JWTMiddleware(userRepo *repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte("your-secret-key"), nil
			})

			if err != nil || !token.Valid {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			claims := token.Claims.(jwt.MapClaims)
			userID := int(claims["user_id"].(float64))

			user, err := userRepo.GetByID(userID)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Check if user is admin
		// This is a simplified check, you should implement proper role checking
		return next(c)
	}
}

func CacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		return next(c)
	}
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		
		echo.New().Logger.Infof("%s %s %s %v",
			c.Request().Method,
			c.Request().URL.Path,
			c.Response().Status,
			time.Since(start),
		)
		return err
	}
}
