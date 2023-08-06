package api

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/regmarmcem/echo-session-demo/controller"
	"github.com/regmarmcem/echo-session-demo/service"
	"gorm.io/gorm"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRouter(db *gorm.DB, store sessions.Store) *echo.Echo {
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}
	e.Renderer = renderer
	e.Static("/assets", "./static/assets")
	e.Static("/css", "./static/css")
	e.Static("/js", "./static/js")

	e.Use(session.Middleware(store))

	uSer := service.NewUserService(db)
	uCon := controller.NewUserController(uSer)

	e.GET("/", sessionHandler)
	e.GET("/home", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	}, CheckSignin)
	e.GET("/signout", uCon.GetSignout, CheckSignin)
	e.GET("/signup", uCon.GetSignup, CheckSignout)
	e.POST("/signup", uCon.PostSignup, CheckSignout)
	e.GET("/signin", uCon.GetSignin, CheckSignout)
	e.POST("/signin", uCon.PostSignin, CheckSignout)
	return e
}

func sessionHandler(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["foo"] = "bar"
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func CheckSignin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Println(err)
			log.Println("CheckSignin() failed")
			c.Redirect(http.StatusSeeOther, "/signin")
		}

		user, ok := sess.Values["user"].(string)

		if !ok || user == "" {
			c.Redirect(http.StatusSeeOther, "/signin")
		}
		return next(c)
	}
}

func CheckSignout(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		user, ok := sess.Values["user"].(string)

		if ok && user != "" {
			c.Redirect(http.StatusSeeOther, "/home")
		}
		return next(c)
	}
}
