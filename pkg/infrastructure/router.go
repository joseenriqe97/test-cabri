package infrastructure

import "github.com/labstack/echo/v4"

func NewRouter(e *echo.Echo, appCtrl AppController) *echo.Echo {
	apiRouter := e.Group("/api")

	apiRouter.GET("/user", appCtrl.User.Created)

	return e

}
