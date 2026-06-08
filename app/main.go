package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("$ ")
	cmd, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("%s: command not found\n", cmd[:len(cmd)-1])
}
