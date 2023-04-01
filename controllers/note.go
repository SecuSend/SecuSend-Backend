package controllers

import (
	"context"
	"log"
	"net/http"
	"secusend/configs"
	"secusend/models"
	"secusend/responses"
	"secusend/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection = configs.GetCollection(configs.DB, "note")

func PostNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := struct {
			Key  string `json:"key"`
			Data string `json:"data"`
		}{}

		//validate the request body
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.GenericResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		//Encrypt the text:
		encrypted, err := services.Encrypt(body.Key, body.Data) //todo key 32bit
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		log.Println(encrypted)

		newNote := models.Note{
			Id:                primitive.NewObjectID(),
			Key:               body.Key,
			Data:              encrypted,
			PasswordProtected: true,
			CreatedAt:         time.Now(),
		}

		result, err := noteCollection.InsertOne(ctx, newNote)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		return c.Status(http.StatusCreated).JSON(responses.GenericResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})

		// c.Status(fiber.StatusOK)
		// return c.JSON(fiber.Map{"test": test})
	}
}

func GetNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		id := c.Query("id")
		var note models.Note
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(id)

		err := noteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&note)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(http.StatusNotFound).JSON(responses.GenericResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": err.Error()}})
			} else {
				return c.Status(http.StatusInternalServerError).JSON(responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
			}
		}

		//Decrypt the text:
		decrypted, err := services.Decrypt(note.Key, note.Data)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		return c.Status(http.StatusOK).JSON(responses.GenericResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": decrypted}})
	}
}
