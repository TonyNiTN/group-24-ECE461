package main

import (
	"fmt"
	"group-24-ECE461/internal/logger"
	"group-24-ECE461/internal/parser"
	"os"

	"group-24-ECE461/internal/config"
	"group-24-ECE461/internal/error"
)

func main() {

	cfg := config.NewConfig()
	if err := cfg.CheckToken(); err != nil {
		fmt.Println(error.NewGeneralError("cfg.CheckToken", err.Error()).Error())
		os.Exit(1)
	}

	logger, err := logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Application")
	// fmt.Println("Starting Application")
	argsWithOutProg := os.Args[1:]
	parser.ParseArguments(argsWithOutProg, logger)
	os.Exit(0)
}
