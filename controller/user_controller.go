package controller

import (
	"html"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
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

	user, err := ctr.s.Signup(email, password)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println(err)
		log.Println("failed to Store.New() in PostSignup()")
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	sess.Values["user"] = user.Email
	err = sess.Save(c.Request(), c.Response().Writer)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	return c.Redirect(http.StatusSeeOther, "/home")
}

func (ctr *UserController) GetSignin(c echo.Context) error {
	return c.Render(http.StatusOK, "signin.html", nil)
}

func (ctr *UserController) PostSignin(c echo.Context) error {
	email := html.EscapeString(c.FormValue("email"))
	password := html.EscapeString(c.FormValue("password"))

	if email == "" || password == "" {
		log.Println("Invalid Email or Password")
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	user, err := ctr.s.Signin(email, password)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signup")
	}

	sess.Values["user"] = user.Email
	err = sess.Save(c.Request(), c.Response().Writer)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	return c.Redirect(http.StatusSeeOther, "/home")
}

func (ctr *UserController) GetSignout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	delete(sess.Values, "user")
	if err = sess.Save(c.Request(), c.Response().Writer); err != nil {
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "home.html")
	}
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "home.html")
}
