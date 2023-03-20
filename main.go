package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pusher/pusher-http-go"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/api/SendMessages", SendMessagePost)

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	e.GET("/", Home)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

type ReqMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type ResMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Handler
func SendMessagePost(c echo.Context) error {

	var message ReqMessage
	err := c.Bind(&message)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// config pusher
	pusherClient := pusher.Client{
		AppID:   "1483564",
		Key:     "f074590e163d64f9e82c",
		Secret:  "c02da4f12d6c3f8ff806",
		Cluster: "ap1",
		Secure:  true,
	}

	resp := ResMessage{
		Username: message.Username,
		Message:  message.Message,
	}
	fmt.Println(resp)

	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}

	// send pusher
	fmt.Println(string(b))
	pusherClient.Trigger("chat", "message", string(b))

	return c.String(http.StatusOK, string(b))
}
