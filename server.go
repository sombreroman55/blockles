package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/sombreroman55/blockles/game"
)

type GameInfo struct {
	Title string
	Id    string
}

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

func handleWstestGet(c echo.Context) error {
	return c.Render(200, "wstest", nil)
}

func handleChatGet(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error("upgrade:", err)
		return err
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}

		log.Printf("recv: %s, type: %d", message, mt)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
	return err
}

func handleHomeGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func handleNewSoloGet(c echo.Context) error {
	return c.Render(http.StatusOK, "newsolo", nil)
}

func handleNewSoloPost(c echo.Context) error {
	title := c.FormValue("gameTitle")
	ghost := c.FormValue("showGhostPiece") == "on"
	next := c.FormValue("showNextQueue") == "on"
	hold := c.FormValue("allowHold") == "on"

	log.Printf("Options for game %s\n", title)
	log.Println("---------------------------")
	log.Printf("ghost: %v\n", ghost)
	log.Printf("next:  %v\n", next)
	log.Printf("hold:  %v\n", hold)

	// TODO: Hook up to game logic here to kick off solo game
	gameId := uuid.New()
	gameUrl := "/solo?id=" + gameId.String()
	c.Response().Header().Set("HX-Redirect", gameUrl)
	c.Response().WriteHeader(303)
	return nil
}

func handleNewMultiGet(c echo.Context) error {
	return c.Render(http.StatusOK, "newmulti", nil)
}

func handleNewMultiPost(c echo.Context) error {
	title := c.FormValue("gameTitle")
	ghost := c.FormValue("showGhostPiece") == "on"
	next := c.FormValue("showNextQueue") == "on"
	hold := c.FormValue("allowHold") == "on"

	log.Printf("Options for game %s\n", title)
	log.Println("---------------------------")
	log.Printf("ghost: %v\n", ghost)
	log.Printf("next:  %v\n", next)
	log.Printf("hold:  %v\n", hold)

	// TODO: Hook up to game logic here to kick off solo game
	gameId := uuid.New()
	gameUrl := "/multi?id=" + gameId.String()
	c.Response().Header().Set("HX-Redirect", gameUrl)
	c.Response().WriteHeader(303)
	return nil
}

func handleSoloGameGet(c echo.Context) error {
	gameInfo := GameInfo{
		Title: "test game",
		Id:    c.QueryParam("id"),
	}
	return c.Render(http.StatusOK, "sologame", gameInfo)
}

func handleMultiGameGet(c echo.Context) error {
	return nil
}

func gameExists(next echo.HandlerFunc) echo.HandlerFunc {
	// TODO: Implement this for real
	return func(c echo.Context) error {
		log.Printf("Game id: %s\n", c.QueryParam("id"))
		return next(c)
	}
}

func serveBlockles() {
	log.Println("Serving Blockles site")
	hub := newHub()
	go hub.run()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "site/static")
	e.GET("/", handleHomeGet)

	e.GET("/newsolo", handleNewSoloGet)
	e.POST("/newsolo", handleNewSoloPost)

	e.GET("/solo", handleSoloGameGet, gameExists)
	e.GET("/multi", handleMultiGameGet, gameExists)

	e.GET("/wstest", handleWstestGet)
	e.GET("/ws", serveWs(hub))
	e.GET("/chat", handleChatGet)

	e.GET("/newmulti", handleNewMultiGet)
	e.POST("/newmulti", handleNewMultiPost)

	e.Renderer = newTemplates()
	game.InitGameManager()
	e.Logger.Fatal(e.Start(":8000"))
}
