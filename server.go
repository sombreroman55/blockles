package main

import (
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/sombreroman55/blockles/game"
)

type GameInfo struct {
	Title      string
	Id         string
	PlayerName string
}

type Templates struct {
	templates *template.Template
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("site/pages/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func writeCookie(c echo.Context, key string, val string) error {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = val
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return nil
}

func handleHomeGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func handleNewSoloGet(c echo.Context) error {
	return c.Render(http.StatusOK, "newsolo", nil)
}

func handleNewSoloPost(c echo.Context) error {
	title := c.FormValue("gameTitle")
	playerName := c.FormValue("playerName")
	ghost := c.FormValue("showGhostPiece") == "on"
	next := c.FormValue("showNextQueue") == "on"
	hold := c.FormValue("allowHold") == "on"

	gmo := game.GameOptions{
		ShowGhostPiece: ghost,
		ShowNextQueue:  next,
		EnableHolding:  hold,
	}

	gameId := game.NewGame(title, 1, gmo)
	gameInstance := game.GetGame(gameId)
	go gameInstance.Run()
	playerId := game.CreateNewPlayer(playerName, gameInstance)

	writeCookie(c, "gameId", gameId.String())
	writeCookie(c, "playerId", playerId.String())
	gameUrl := "/solo?id=" + gameId.String()
	c.Response().Header().Set("HX-Redirect", gameUrl)
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil
}

func handleSoloGameGet(c echo.Context) error {
	gameId := c.QueryParam("id")
	gameInstance := game.GetGame(uuid.MustParse(gameId))
	playerCookie, err := c.Cookie("playerId")
	if err != nil {
		return err
	}
	playerId := uuid.MustParse(playerCookie.Value)
	player := game.GetPlayer(playerId)

	gameInfo := GameInfo{
		Title:      gameInstance.Name,
		Id:         gameId,
		PlayerName: player.Name,
	}
	return c.Render(http.StatusOK, "sologame", gameInfo)
}

func handleSoloGameWebsocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error(err)
		return err
	}

	gameId := c.QueryParam("id")
	gameInstance := game.GetGame(uuid.MustParse(gameId))

	playerCookie, err := c.Cookie("playerId")
	if err != nil {
		// TODO: Return unauthorized error
		panic("We don't have a player cookie")
	}

	playerId := uuid.MustParse(playerCookie.Value)
	if !game.PlayerExists(playerId) {
		// TODO: Return unauthorized error
		panic("This player doesn't exist")
	}

	player := game.GetPlayer(playerId)
	player.AttachWebsocket(conn)
	gameInstance.AddNewPlayer(player)

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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "site/static")
	e.GET("/", handleHomeGet)

	e.GET("/newsolo", handleNewSoloGet)
	e.POST("/newsolo", handleNewSoloPost)

	e.GET("/solo", handleSoloGameGet, gameExists)
	e.GET("/solows", handleSoloGameWebsocket)

	e.Renderer = newTemplates()
	game.InitGameManager()
	e.Logger.Fatal(e.Start(":8000"))
}
