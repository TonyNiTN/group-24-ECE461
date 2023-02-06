package parser

import (
	"bufio"
	"fmt"
	"group-24-ECE461/internal/cli"
	"os"
	"strings"
)

func ParseArguments(argsWithProg []string) (e error) {
	if strings.Contains(argsWithProg[0], "/") {
		file, err := os.Open(argsWithProg[0])
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
		cli.Score(lines)
		return nil
	}

	if argsWithProg[0] == "test" {
		cli.Test()
	}

	if argsWithProg[0] == "build" {
		cli.Build()
	}

	if argsWithProg[0] == "install" {
		cli.Install()
	}

	return nil
}
