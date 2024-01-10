package handle

import (
	"context"
	"encoding/json"
	"log"
	"myapp/config"
	"myapp/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var projectCollection *mongo.Collection = configs.GetCollection(configs.DB, "Projects")
var chatCollection *mongo.Collection = configs.GetCollection(configs.DB, "Chats")
var chatroomCollection *mongo.Collection = configs.GetCollection(configs.DB, "Chatrooms")
var validate = validator.New()

func AddChatRoom(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	body := new(models.Chatroom)
	if err := c.BodyParser(body); err != nil {
		return err
	}
	doc := body
	doc.CreateDate = time.Now()

	result, err := chatroomCollection.InsertOne(ctx, doc)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.JSON(models.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result.InsertedID.(primitive.ObjectID).Hex()}})
}

func GetChat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	chatId := c.Params("chatId")
	getDb, err := chatCollection.Find(ctx, bson.M{"chatroomid": chatId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	var results []models.Chat
	if err = getDb.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	return c.JSON(models.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": results}})
}

func GetChatRoom(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	projectId := c.Params("projectId")
	getDb, err := chatroomCollection.Find(ctx, bson.M{"projectid": projectId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	var results []models.ChatroomView
	if err = getDb.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	for i := range results {
		count, err := chatCollection.CountDocuments(context.Background(), bson.M{"chatroomid": results[i].UUID, "isread": false})
		if err != nil {
			log.Fatal(err)
		}
		results[i].CountRead = count
	}

	return c.JSON(models.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": results}})
}

func AddChat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	body := new(models.Chat)
	if err := c.BodyParser(body); err != nil {
		return err
	}
	doc := body
	result, err := chatCollection.InsertOne(ctx, doc)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.JSON(models.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result.InsertedID.(primitive.ObjectID).Hex()}})
}

func SendMessage(c *websocket.Conn) {
	curr_hub := new(Hub)
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		id := c.Params("id")

		curr_hub = getCurrHub(id)
		(*curr_hub).register <- c

		curr_hub.current <- c
		(*curr_hub).broadcast <- append([]byte{}, msg...)

		if mt, msg, err = c.ReadMessage(); err != nil {
			c.Close()
			break
		}
		data := &models.Chat{}
		if err := json.Unmarshal(msg, &data); err != nil {
			c.Close()
			break
		}

		if data.ChatroomID != "" {
			_, errs := chatCollection.InsertOne(context.Background(), data)
			if errs != nil {
				log.Println(errs)
			}
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			c.Close()
			break
		}

	}

}

func UpdateChat(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := c.Params("chatId")
	log.Println(id)
	filter := bson.D{{Key: "chatroomid", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isread", Value: true}}}}
	result, err := chatCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.JSON(models.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})

}
