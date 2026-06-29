package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

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

func parseArgs(input string) []string {
	var args []string
	var builder strings.Builder

	inSingleQuotes := false
	inDoubleQuotes := false
	besideBackSlash := false
	hasContent := false

	for i, r := range input {
		if besideBackSlash {
			builder.WriteRune(r)
			besideBackSlash = false
			hasContent = true
			continue
		}

		switch {
		case r == '\\':
			if !inSingleQuotes && !inDoubleQuotes {
				besideBackSlash = true
				hasContent = true
			} else if inDoubleQuotes {
				if input[i+1] == '"' || input[i+1] == '\\' {
					besideBackSlash = true
					hasContent = true
				} else {
					builder.WriteRune(r)
					hasContent = true
				}
			} else {
				builder.WriteRune(r)
				hasContent = true
			}
		case r == '\'':
			if !inDoubleQuotes {
				// Toggle single quote state
				inSingleQuotes = !inSingleQuotes
				hasContent = true
			} else {
				// Inside double quotes, take everything literally
				builder.WriteRune(r)
				hasContent = true
			}
		case r == '"':
			if !inSingleQuotes {
				// Toggle double quote state
				inDoubleQuotes = !inDoubleQuotes
				hasContent = true
			} else {
				// Inside single quotes, take everything literally
				builder.WriteRune(r)
				hasContent = true
			}
		case inSingleQuotes || inDoubleQuotes:
			// Inside quotes, take everything literally
			builder.WriteRune(r)
			hasContent = true
		case unicode.IsSpace(r):
			// Outside quotes and hit whitespace, finalize content if present
			if hasContent {
				args = append(args, builder.String())
				builder.Reset()
				hasContent = false
			}
		default:
			// Outside quotes, and hit normal text, append to builder
			builder.WriteRune(r)
			hasContent = true
		}
	}

	// Flush remaining content
	if hasContent {
		args = append(args, builder.String())
	}

	return args
}
