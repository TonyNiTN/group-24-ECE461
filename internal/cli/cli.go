package cli

import (
	"fmt"
)

func Install() {
	fmt.Println("Install....")
}

func Score(urlLinks []string) {
	fmt.Println("Scoring.....")
	fmt.Println(urlLinks)
}

func Build() {
	fmt.Println("Building......")
}

func Test() {
	fmt.Println("Testing.....")
}
