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
			Password     *string    `json:"password"`
			Data         string     `json:"data"`
			SelfDestruct bool       `json:"selfDestruct"`
			ExpireAt     *time.Time `json:"expireAt"`
		}{}

		//validate the request body
		if err := c.BodyParser(&body); err != nil {
			return responses.BadRequestResponse(c, "Parser error")
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
				return responses.InternalServerErrorResponse(c, "Encrypt error")
			}
			data = encrypted
		} else {
			passwordProtected = false
			data = body.Data
		}

		randomKey, err := services.GenerateUniqueKey()
		if err != nil {
			log.Println(err)
			return responses.InternalServerErrorResponse(c, "Key generation error")
		}

		newNote := models.Note{
			Id:                primitive.NewObjectID(),
			Key:               randomKey,
			Data:              data,
			PasswordProtected: passwordProtected,
			SelfDestruct:      body.SelfDestruct,
			ExpireAt:          body.ExpireAt,
			CreatedAt:         time.Now(),
		}

		result, err := noteCollection.InsertOne(ctx, newNote)
		if err != nil {
			log.Println(err)
			return responses.InternalServerErrorResponse(c, "Insertion error")
		}
		log.Println(result)

		return responses.CreatedResponse(c, &fiber.Map{"key": newNote.Key})
	}
}

func GetNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := struct {
			Key      string  `json:"key"`
			Password *string `json:"password"`
		}{}

		//validate the request body
		if err := c.BodyParser(&body); err != nil {
			return responses.BadRequestResponse(c, "Parser error")
		}

		var note models.Note

		//Get the note
		err := noteCollection.FindOne(ctx, bson.M{"key": body.Key}).Decode(&note)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return responses.NotFoundResponse(c, "Key not found!")
			} else {
				log.Println(err)
				return responses.InternalServerErrorResponse(c, "DB error")
			}
		}

		//Expiry
		if note.ExpireAt != nil && time.Now().After(*note.ExpireAt) {
			return responses.UnauthorizedResponse(c, "Expired!")
		}

		var data string

		//Password aes decrypt
		if note.PasswordProtected == true {
			if body.Password == nil || *body.Password == "" {
				return responses.UnauthorizedResponse(c, "Wrong password!")
			}

			//Decrypt the text:
			decrypted, err := services.Decrypt(*body.Password, note.Data)
			if err != nil {
				log.Println(err)
				return responses.UnauthorizedResponse(c, "Wrong Password!")
			}
			data = decrypted
		} else {
			data = note.Data
		}

		//SelfDestruct
		if note.SelfDestruct == true {
			_, err := noteCollection.DeleteOne(ctx, bson.M{"key": body.Key})
			if err != nil {
				log.Println(err)
				return responses.InternalServerErrorResponse(c, "DB error")
			}
		}

		return responses.OKResponse(c, &fiber.Map{"data": data, "selfDestruct": note.SelfDestruct})
	}
}
