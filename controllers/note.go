package controllers

import (
	"context"
	"log"
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

func CreatetNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := struct {
			Password *string `json:"password"`
			Data     string  `json:"data"`
		}{}

		//validate the request body
		if err := c.BodyParser(&body); err != nil {
			return responses.BadRequestResponse(c, err.Error())
		}

		var data string
		var passwordProtected bool

		// Check if password is not null or empty
		if body.Password != nil && *body.Password != "" {
			passwordProtected = true

			//Encrypt the text with the password:
			encrypted, err := services.Encrypt(*body.Password, body.Data) //todo key 32bit
			if err != nil {
				log.Println(err)
				return responses.InternalServerErrorResponse(c, err.Error())
			}
			data = encrypted
		} else {
			passwordProtected = false
			data = body.Data
		}

		newNote := models.Note{
			Id:                primitive.NewObjectID(),
			Key:               "test", //todo
			Data:              data,
			PasswordProtected: passwordProtected,
			CreatedAt:         time.Now(),
		}

		result, err := noteCollection.InsertOne(ctx, newNote)
		if err != nil {
			return responses.InternalServerErrorResponse(c, err.Error())
		}

		return responses.OKResponse(c, &fiber.Map{"result": result})

		// c.Status(fiber.StatusOK)
		// return c.JSON(fiber.Map{"test": test})
	}
}

func GetNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := struct {
			Id       string  `json:"id"`
			Password *string `json:"password"`
		}{}

		//validate the request body
		if err := c.BodyParser(&body); err != nil {
			return responses.BadRequestResponse(c, err.Error())
		}

		var note models.Note

		objId, _ := primitive.ObjectIDFromHex(body.Id)
		err := noteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&note)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return responses.NotFoundResponse(c, err.Error())
			} else {
				return responses.InternalServerErrorResponse(c, err.Error())
			}
		}

		var data string

		if note.PasswordProtected == true {
			if body.Password != nil && *body.Password != "" {
				return responses.UnauthorizedResponse(c, err.Error())
			}

			//Decrypt the text:
			decrypted, err := services.Decrypt(*body.Password, note.Data)
			if err != nil {
				log.Println(err)
				return responses.InternalServerErrorResponse(c, err.Error())
			}
			data = decrypted
		} else {
			data = note.Data
		}

		return responses.CreatedResponse(c, &fiber.Map{"data": data})
	}
}
