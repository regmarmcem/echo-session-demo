package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/regmarmcem/echo-session-demo/service"
)

type UserController struct {
	s *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{s: s}
}

func (ctr *UserController) GetSignup(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.html", nil)
}
