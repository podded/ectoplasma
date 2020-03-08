package podgoo

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/podded/bouncer/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	PodGoo struct {
		BoundHost string
		BoundPort int

		dbClient *mongo.Client
		redis    *redis.Client
		client   *client.BouncerClient
	}
)

func NewPodGoo(bouncerAddress string, maxTimeout time.Duration, descriptor string) (goop *PodGoo) {
	sludge := &PodGoo{}

	bc, version, err := client.NewBouncer(bouncerAddress, maxTimeout, descriptor)
	if err != nil {
		log.Fatalf("Failed to connect to bouncer....: %s\n", err)
	}
	log.Printf("Connected to bouncer. version %s\n", version)

	sludge.client = bc

	// TODO Make the db connection configurable
	clientOptions := options.Client().ApplyURI("mongodb://" + "localhost" + ":" + "27017")
	cl, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}

	// Check the connection
	err = cl.Ping(context.TODO(), nil)
	if err != nil {
		return
	}

	sludge.dbClient = cl

	rclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rclient.Ping().Result()
	if err != nil || pong != "PONG" {
		log.Fatalf("Failed to connect to redis: %s - %s\n", pong, err)
	}
	sludge.redis = rclient

	return sludge
}
