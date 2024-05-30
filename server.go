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

func handleNewSoloPost(c echo.Context) error {
	title := c.FormValue("gameTitle")
	ghost := c.FormValue("showGhostPiece") == "on"
	next := c.FormValue("showNextQueue") == "on"
	hold := c.FormValue("allowHold") == "on"
	crot := c.FormValue("classicRotation") == "on"
	clock := c.FormValue("classicLockdown") == "on"
	debug := c.FormValue("debug") == "on"

	log.Printf("Options for game %s\n", title)
	log.Println("---------------------------")
	log.Printf("ghost: %v\n", ghost)
	log.Printf("next:  %v\n", next)
	log.Printf("hold:  %v\n", hold)
	log.Printf("crot:  %v\n", crot)
	log.Printf("clock: %v\n", clock)
	log.Printf("debug: %v\n", debug)

	// TODO: Hook up to game logic here to kick off solo game
	return c.Render(http.StatusOK, "newsolorecv", title)
}

func serveBlockles() {
	log.Println("Serving Blockles site")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "site/static")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/newsolo", func(c echo.Context) error {
		return c.Render(http.StatusOK, "newsolo", nil)
	})
	e.POST("/newsolo", handleNewSoloPost)

	e.Renderer = newTemplates()
	e.Logger.Fatal(e.Start(":8000"))
}
