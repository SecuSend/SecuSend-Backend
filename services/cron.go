package services

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"

	"secusend/configs"
)

func StartCronJobs() {
	c := cron.New()
	_, err := c.AddFunc("0 4 * * *", CleanUpExpiredNotes)
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}
	c.Start()
}

func CleanUpExpiredNotes() {
	noteCollection := configs.GetCollection(configs.DB, "note")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"expireat": bson.M{"$lt": time.Now().UTC()}}
	result, err := noteCollection.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("Error deleting expired notes: %v", err)
		return
	}
	log.Printf("Cleaned up %v expired notes", result.DeletedCount)
}
