package routes

import (
	"myapp/handle"
	"github.com/gofiber/websocket/v2"
	"github.com/gofiber/fiber/v2"
)

func ChatRoutes(api fiber.Router , app *fiber.App) {
	api.Post("/addChat", handle.AddChatRoom)
	api.Get("/getChat/:chatId", handle.GetChat)
	api.Get("/getChatRoom/:projectId", handle.GetChatRoom)
	api.Put("/updateChat/:chatId", handle.UpdateChat)

	app.Get("/ws/:id", websocket.New(handle.SendMessage))

}