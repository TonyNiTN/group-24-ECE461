package main

import (
	"fmt"
	"group-24-ECE461/internal/error"
	"group-24-ECE461/internal/logger"
)

func main() {
	logger, err := logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	infoError := error.NewRequestError("GraphQL", "Error", 400)
	logger.Info(infoError.Error())

	logger.Debug("debug")

	logger.Debug("Starting Application")
	fmt.Println("hello world!")
}
