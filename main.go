package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Body        string             `json:"body"`
	Translation string             `json:"translation"`
}

var collection *mongo.Collection

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/words", getTodos)
	app.Post("/api/words", createTodos)
	app.Delete("/api/words/:id", deleteTodos)
	app.Delete("/api/words", deleteAll)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func translateWord(word string) []byte {
	url := "https://deep-translate1.p.rapidapi.com/language/translate/v2"

	payload := strings.NewReader(fmt.Sprintf(`{ "q" : "%v", "source": "en", "target": "ru"}`, word))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-rapidapi-key", "2c503f43edmsh89fee0e8f7af91cp149179jsn1be39bb05ae3")
	req.Header.Add("x-rapidapi-host", "deep-translate1.p.rapidapi.com")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return parseWord(body)
}

func parseWord(word []byte) []byte {
	count := 0
	rb := 0
	lb := 0
	for i := len(word) - 1; i > 0; i-- {
		if word[i] == '"' {
			count++
			if count == 1 {
				rb = i
			}
			if count == 2 {
				lb = i
				break
			}
		}
	}
	return word[lb+1 : rb]
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}
func createTodos(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(404).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}
	s := string(translateWord(todo.Body))
	decodedString, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return err
	}
	todo.Translation = decodedString
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(todo)
}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}
	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteAll(c *fiber.Ctx) error {
	_, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
