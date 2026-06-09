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

	// Search for everything in path
	paths := os.Getenv("PATH")
	executableDirs := strings.SplitSeq(paths, string(os.PathListSeparator))

	// Iterate over all the directories to check if
	// the command is found in any dir and is an executable
	for dirPath := range executableDirs {
		isExecutable, execCheckErr := checkFileIsExecutable(dirPath, args)
		if execCheckErr != nil {
			continue
		}

		if isExecutable {
			fmt.Printf("%s is %s/%s\n", args, dirPath, args)
			return
		}
	}

	fmt.Printf("%s: not found\n", args)
}

func checkFileIsExecutable(dirPath string, progName string) (bool, error) {
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return false, err
	}

	if !dirInfo.IsDir() {
		return false, fmt.Errorf("%s is not a directory", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Name() == progName && !entry.IsDir() {
			fileInfo, infoErr := entry.Info()
			if infoErr != nil {
				return false, infoErr
			}

			isExecutable := fileInfo.Mode().Perm()&0111 != 0
			return isExecutable, nil
		}
	}

	return false, nil
}
