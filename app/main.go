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
		cmdArgs := parseArgs(cmd)
		prog, args := cmdArgs[0], cmdArgs[1:]
		switch prog {
		case "exit":
			os.Exit(0)
		case "echo":
			execEcho(args)
		case "type":
			execType(args)
		case "pwd":
			execPwd()
		case "cd":
			execChangeDir(args)
		default:
			execProg(prog, args)
		}
	}
}

func execEcho(args []string) {
	output := strings.Join(args, " ")
	fmt.Println(output)
}

func execType(args []string) {
	if len(args) < 1 {
		return
	}

	for _, arg := range args {
		builtinCommands := []string{"exit", "echo", "type", "pwd", "cd"}
		if slices.Contains(builtinCommands, arg) {
			fmt.Printf("%s is a shell builtin\n", arg)
			continue
		}

		// Check if args is executable
		execPath := getFullExecPath(arg)
		if execPath != "" {
			fmt.Printf("%s is %s\n", arg, execPath)
			continue
		}

		fmt.Printf("%s: not found\n", arg)
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

func execChangeDir(args []string) {
	dirPath := strings.Join(args, " ")
	if strings.HasPrefix(dirPath, "~") {
		homeDir := os.Getenv("HOME")
		dirPath = strings.Replace(dirPath, "~", homeDir, 1)
	}

	err := os.Chdir(dirPath)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dirPath)
	}
}

func execProg(prog string, args []string) {
	progPath := getFullExecPath(prog)
	if progPath == "" {
		fmt.Printf("%s: command not found\n", prog)
		return
	}

	cmd := exec.Command(prog, args...)

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
