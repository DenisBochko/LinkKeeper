package main

import (
	tg "LinkKeeper/TGinter"
	db "LinkKeeper/database"
	"context"
	"fmt"
	//"runtime/trace"
	//"sync"
	//"time"
)

func main() {
	ctxForDB, cancelDB := context.WithCancel(context.Background())
	ctxForTG, cancelTG := context.WithCancel(context.Background())
	is_End := false
	//wg := sync.WaitGroup{}

	sChan := make(chan db.Field, 100)
	gChan := make(chan db.Field, 100)
	dChan := make(chan db.Field, 100)
	doChan := make(chan db.Field, 100)
	rChan := make(chan []db.Field, 100)
	database := db.DataBase{
		ConnStr:    "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable",
		DriverName: "postgres",
	}

	TGinter := tg.TGinter{
		OK: true,
	}
	//wg.Add(1)
	go func(ctx context.Context, saveChan, getChan, deleteChan, deleteOfItemChan <-chan db.Field, receive chan<- []db.Field) {
		database.Start(ctx, saveChan, getChan, deleteChan, deleteOfItemChan, receive)
	}(ctxForDB, sChan, gChan, dChan, doChan, rChan)

	go func(ctx context.Context, saveChan, getChan, deleteChan, deleteOfItemChan chan<- db.Field, receiveChan <-chan []db.Field) {
		TGinter.Start(ctx, saveChan, getChan, deleteChan, deleteOfItemChan, receiveChan)
	}(ctxForTG, sChan, gChan, dChan, doChan, rChan)

	fmt.Scan(&is_End)
	cancelTG()
	cancelDB()
}

func Printer(inChan <-chan []db.Field) {
	for fields := range inChan {
		for _, field := range fields {
			fmt.Printf("ID: %d, UserID: %s, UserURL: %s\n", field.ID, field.UserID, field.UserURL)
		}
	}
}
