package main

import (
	db "LinkKeeper/database"
	"context"
	"fmt"
	"sync"
	"time"
)

type Field struct {
	ID      int
	UserID  string
	UserURL string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	sChan := make(chan db.Field, 10)
	gChan := make(chan db.Field, 10)
	dChan := make(chan db.Field, 10)
	rChan := make(chan []db.Field, 10)
	database := db.DataBase{
		ConnStr:    "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable",
		DriverName: "postgres",
	}
	wg.Add(1)
	go func(ctx context.Context, saveChan, getChan, deleteChan <-chan db.Field, receive chan<- []db.Field) {
		database.Start(ctx, sChan, gChan, dChan, rChan)
	}(ctx, sChan, gChan, dChan, rChan)

	go func(inChan <-chan []db.Field) {
		defer wg.Done()
		Printer(rChan)
	}(rChan)
	// sChan <- db.Field{
	// 	ID:      0,
	// 	UserID:  "666",
	// 	UserURL: "https://misha",
	// }

	gChan <- db.Field{
		ID:      0,
		UserID:  "666",
		UserURL: "",
	}

	dChan <- db.Field{
		ID:      0,
		UserID:  "666",
		UserURL: "",
	}

	gChan <- db.Field{
		ID:      0,
		UserID:  "666",
		UserURL: "",
	}
	

	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
	
}

func Printer(inChan <-chan []db.Field) {
	for fields := range inChan {
		for _, field := range fields {
			fmt.Printf("ID: %d, UserID: %s, UserURL: %s\n", field.ID, field.UserID, field.UserURL)
		}
	}
}
