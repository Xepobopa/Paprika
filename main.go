package main

import (
	"Paprika/app"
	"Paprika/utils"
	"log"
)

func main() {
	logFile := "logs.log"
	customWriter, err := app.LogToFile(logFile)
	if err != nil {
		log.Fatalf("Panic while trying to create custom writer (logs): %v\n", err)
	}

	cfg, err := utils.NewConfig(".env")
	if err != nil {
		log.Fatalf("Panic while trying to create config: %v", err)
	}

	if err := app.Run(cfg, customWriter); err != nil {
		log.Fatalf("Panic in the main app: %v", err)
	}
}
