package infrastructure

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joseenriqe97/test-cabri/config"
	"github.com/joseenriqe97/test-cabri/pkg/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, appCtrl AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//PUBLIC
	publicGroup := e.Group("/public")
	publicGroup.POST("/auth", appCtrl.User.Authenticate)
	publicGroup.POST("/user", appCtrl.User.Created)

	//API
	apiGroup := e.Group("/api", jwtMiddleware)
	apiGroup.GET("/user", appCtrl.User.GetByID)

	return e

}

func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "Token not provided or invalid format")
		}

		jwtToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &application.SignatureOfAuth{}
		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetSecretJWT()), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("id", claims.ID)

		return next(c)
	}
}
