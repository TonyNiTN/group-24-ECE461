package parser

import (
	"bufio"
	"fmt"
	"group-24-ECE461/internal/cli"
	"os"
	"strings"

	"go.uber.org/zap"
)

func ParseArguments(argsWithProg []string, logger *zap.Logger) (e error) {
	arg := argsWithProg[0]
	if strings.Contains(arg, "\\") {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Println("Error opening URL file:", err)
			return err
		}
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), "\n")
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading URL file:", err)
			return err
		}
		cli.Score(lines, logger)
		return nil
	}

	if argsWithProg[0] == "test" {
		cli.Test(logger)
	}

	if argsWithProg[0] == "build" {
		cli.Build(logger)
	}

	if argsWithProg[0] == "install" {
		cli.Install(logger)
	}

	return nil
}
