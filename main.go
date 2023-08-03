package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	help := "Please specify a subcommand:\n\t" +
		"- 'ls' to fetch all your TODOs."

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(1)
	}

	var cmd func([]string) error
	switch os.Args[1] {
	case "ls":
		cmd = ls

	default:
		fmt.Println(help)
		os.Exit(1)
	}

	if err := cmd(os.Args[2:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Item struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Todo struct {
	Items []Item `json:"items"`
}

func ls(args []string) error {
	path := "todo.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("opening file[%s]: %w", path, err)
	}

	var todo Todo
	if err := json.Unmarshal(content, &todo); err != nil {
		return fmt.Errorf("unmarshaling content[%s]: %w", content, err)
	}

	fmt.Println("Content: ")
	for _, i := range todo.Items {
		fmt.Println("Item:", i.Name, "-", "Message:", i.Message)
	}

	return nil
}
