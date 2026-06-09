package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
			checkCmdType(args)
		case "pwd":
			execPwd()
		case "cd":
			execChangeDir(args)
		default:
			checkAndRunCmd(prog, args)
		}
	}
}

func execChangeDir(args string) {
	dirPath := args
	if strings.HasPrefix(dirPath, "~") {
		homeDir := os.Getenv("HOME")
		dirPath = strings.Replace(dirPath, "~", homeDir, 1)
	}

	err := os.Chdir(dirPath)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dirPath)
	}
}

func execPwd() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error finding current dir: %s\n", err)
		return
	}

	fmt.Println(currDir)
}

func checkAndRunCmd(prog string, args string) {
	progPath := getFullExecPath(prog)
	if progPath == "" {
		fmt.Printf("%s: command not found\n", prog)
		return
	}

	arguments := strings.Fields(args)
	cmd := exec.Command(prog, arguments...)

	// Redirect command output directly to the
	// terminal standard streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run() starts the process and blocks until completion
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running %s: %s\n", prog, err)
	}
}

func checkCmdType(args string) {
	builtinCommands := []string{"exit", "echo", "type", "pwd"}
	if slices.Contains(builtinCommands, args) {
		fmt.Printf("%s is a shell builtin\n", args)
		return
	}

	// Check if args is executable
	execPath := getFullExecPath(args)
	if execPath != "" {
		fmt.Printf("%s is %s\n", args, execPath)
		return
	}

	fmt.Printf("%s: not found\n", args)
}

func getFullExecPath(progNameOrArgs string) string {
	// Search for everything in path
	paths := os.Getenv("PATH")
	executableDirs := strings.SplitSeq(paths, string(os.PathListSeparator))

	// Iterate over all the directories to check if
	// the command is found in any dir and is an executable
	for dirPath := range executableDirs {
		isExecutable, err := checkFileIsExecutable(dirPath, progNameOrArgs)
		if err != nil {
			continue
		}

		if isExecutable {
			return fmt.Sprintf("%s%c%s", dirPath, os.PathSeparator, progNameOrArgs)
		}
	}

	return ""
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
