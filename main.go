package main

import (

	"github.com/gofiber/fiber/v2"

	// "myapp/routes"
	"myapp/config"
	"myapp/routes"
    "myapp/handle"

	"github.com/gofiber/websocket/v2"
	"github.com/gofiber/template/html/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    
)

func main() {

    go handle.HubRunner()
    engine := html.New("./html", ".html")
    app := fiber.New(fiber.Config{
        Views: engine,
    })
 
    app.Use(cors.New())
	configs.ConnectDB()

	api := app.Group("/api")

    
	app.Use("/ws", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            c.Locals("allowed", true)
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })
    
    app.Get("/" , func(c *fiber.Ctx) error {
        return c.Render("index", fiber.Map{"Title": "Chat Service"})
    })
    
    routes.ChatRoutes(api , app)
  


    app.Listen(":3000")
}