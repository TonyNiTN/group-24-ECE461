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

	newError := error.NewGraphQLError("query failed", "query.getUser")

	logger.Info(newError.Error())

	logger.Info("Starting Application")
	fmt.Println("hello world!")
}
