package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		cmd = strings.TrimSpace(cmd)
		if cmd == "exit" {
			os.Exit(0)
		}

		if strings.HasPrefix(cmd, "echo ") {
			fmt.Println(cmd[5:])
			continue
		}

		fmt.Printf("%s: command not found\n", cmd)
	}
}
