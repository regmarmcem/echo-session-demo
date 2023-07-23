package api

import (
	"html/template"
	"io"
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

func NewRouter(db *gorm.DB) *echo.Echo {
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}
	e.Renderer = renderer
	e.Static("/assets", "./static/assets")
	e.Static("/css", "./static/css")
	e.Static("/js", "./static/js")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	uSer := service.NewUserService(db)
	uCon := controller.NewUserController(uSer)

	e.GET("/", sessionHandler)
	e.GET("/home", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	})
	e.GET("/signup", uCon.GetSignup)
	e.POST("/signup", uCon.PostSignup)
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
