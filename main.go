package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const nodeModules = "node_modules"

var filesToDel = make([]string, 0)

func main() {
	root := flag.String("root", ".", "Root directory to start")
	flag.Parse()
	if err := run(*root, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(filesToDel) == 0 {
		fmt.Fprintf(os.Stdout, "No %s directories found under '%s', exiting...\n", nodeModules, *root)
		os.Exit(0)
	}
	if err := confirmDel(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer) error {
	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filterOut(path, info) {
				return nil
			}

			filesToDel = append(filesToDel, path)

			return nil
		})
}

func confirmDel(out io.Writer) error {
	fmt.Fprintf(out, "Found %d %s folders:\n%s\n", len(filesToDel), nodeModules, strings.Join(filesToDel, "\n"))
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "Are you sure you wish to delete these directories? (y/n): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.Contains(input, "\r\n") {
		input = strings.TrimSuffix(input, "\r\n")
	} else if strings.Contains(input, "\n") {
		input = strings.TrimSuffix(input, "\n")
	}

	switch input {
	case "y":
		fmt.Fprintf(out, "Deleting directories...\n")
		return delFiles(filesToDel)
	case "n":
		fmt.Fprintf(out, "Leaving the directories intact...\n")
	default:
		fmt.Fprintf(out, "Leaving the directories intact...\n")
	}

	return nil
}

func delFiles(paths []string) error {
	for _, path := range paths {
		if err := delFile(path); err != nil {
			return err
		}
	}

	return nil
}

func delFile(path string) error {
	return os.RemoveAll(path)
}

func filterOut(path string, info os.FileInfo) bool {
	if !info.IsDir() || info.Name() != nodeModules {
		return true
	}
	return false
}
