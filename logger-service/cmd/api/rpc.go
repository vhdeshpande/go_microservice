package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

// RPC server
type RPCServer struct {
}

// Data received for RPC methods
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// write the log to mongo
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// response to send
	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
