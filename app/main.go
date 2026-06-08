package main

import (
	"fmt"
)

func main() {
	var cmd string
	fmt.Print("$ ")
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("%s: command not found\n", cmd)
}
