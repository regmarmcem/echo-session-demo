package controller

import (
	"html"
	"log"
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

func (ctr *UserController) PostSignup(c echo.Context) error {
	email := html.EscapeString(c.FormValue("email"))
	password := html.EscapeString(c.FormValue("password"))

	if email == "" || password == "" {
		log.Println("Invalid Email or Password")
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	_, err := ctr.s.Signup(email, password)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	return c.Redirect(http.StatusSeeOther, "/home")
}
