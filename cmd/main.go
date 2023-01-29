package main

import (
	"fmt"
	"group-24-ECE461/internal"
);

func main() {
	logger, err := utils.InitLogger()
	if err != nil {
		fmt.Println("Could not init logger", err)
		return
	}
	defer logger.Sync()

	logger.Info("Starting Application")
	fmt.Println("hello world!");
}