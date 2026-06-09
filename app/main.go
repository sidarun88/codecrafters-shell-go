package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
		prog, args, _ := strings.Cut(cmd, " ")
		switch prog {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(args)
		case "type":
			typeCmdOutput(args)
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}

func typeCmdOutput(args string) {
	builtinCommands := []string{"exit", "echo", "type"}
	if slices.Contains(builtinCommands, args) {
		fmt.Printf("%s is a shell builtin\n", args)
		return
	}

	fmt.Printf("%s: not found\n", args)
}
