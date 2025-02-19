package application

import "github.com/labstack/echo/v4"

type UserInterface interface {
	Created(c echo.Context) error
}

type userApplication struct{}

func NewUserApplication() UserInterface {
	return &userApplication{}
}

func (u *userApplication) Created(c echo.Context) error {

	return c.JSON(200, "PING 200!")
}
