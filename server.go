package main

import (
	"io"
	"html/template"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("site/pages/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func serveBlockles() {
	log.Println("Serving Blockles site")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "site/static")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	e.Renderer = newTemplates()
	e.Logger.Fatal(e.Start(":8000"))
}
