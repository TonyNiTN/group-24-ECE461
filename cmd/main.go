package main

import (
	"fmt"
	"group-24-ECE461/internal/logger"
	"group-24-ECE461/internal/parser"
	"os"

	"group-24-ECE461/internal/config"
)

func main() {

	cfg := config.NewConfig()

	if err := cfg.CheckToken(); err != nil {
		fmt.Println(err.Error())
		return
	}

	logger, err := logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	logger.Info("Starting Application")
	fmt.Println("Starting Application")
	argsWithOutProg := os.Args[1:]
	parser.ParseArguments(argsWithOutProg, logger)
}
