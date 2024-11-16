package main

import (
	tg "LinkKeeper/TGinter"
	ai "LinkKeeper/analyzer"
	db "LinkKeeper/database"
	"context"
	"fmt"
) 

func main() {
	ctxForDB, cancelDB := context.WithCancel(context.Background())
	ctxForTG, cancelTG := context.WithCancel(context.Background())
	ctxForAI, cancelAI := context.WithCancel(context.Background())
	is_End := false

	sChan := make(chan db.Field, 100)
	gChan := make(chan db.Field, 100)
	dChan := make(chan db.Field, 100)
	doChan := make(chan db.Field, 100)
	rChan := make(chan []db.Field, 100)
	sAiChan := make(chan ai.Field, 100)
	gAiChan := make(chan ai.Field, 100)

	database := db.DataBase{
		ConnStr:    "user=postgres dbname=LinkKeeper password=postgres host=localhost sslmode=disable",
		DriverName: "postgres",
	}

	TGinter := tg.TGinter{
		OK: true,
	}

	Analyzer := ai.Analyzer{
		OK: true,
	}

	go func(ctx context.Context, saveChan, getChan, deleteChan, deleteOfItemChan <-chan db.Field, receive chan<- []db.Field) {
		database.Start(ctx, saveChan, getChan, deleteChan, deleteOfItemChan, receive)
	}(ctxForDB, sChan, gChan, dChan, doChan, rChan)

	go func(ctx context.Context, saveChan, getChan, deleteChan, deleteOfItemChan chan<- db.Field, receiveChan <-chan []db.Field, sendAiChan chan<- ai.Field, getAiChan <-chan ai.Field) {
		TGinter.Start(ctx, saveChan, getChan, deleteChan, deleteOfItemChan, receiveChan, sendAiChan, getAiChan)
	}(ctxForTG, sChan, gChan, dChan, doChan, rChan, sAiChan, gAiChan)

	go func(ctx context.Context, inputChan <-chan ai.Field, outputChan chan<- ai.Field) {
		Analyzer.Start(ctx, inputChan, outputChan)
	}(ctxForAI, sAiChan, gAiChan)

	fmt.Scan(&is_End)
	cancelTG()
	cancelDB()
	cancelAI()
}
