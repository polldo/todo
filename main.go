package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	help := "Please specify a subcommand:\n\t" +
		"- 'ls' to fetch all your TODOs.\n\t" +
		"- 'add' to add a new TODO.\n\t" +

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(1)
	}

	var cmd func([]string) error
	switch os.Args[1] {
	case "ls":
		cmd = ls

	case "add":
		cmd = add

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

func add(args []string) error {
	flag := flag.NewFlagSet("todo add", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to add.")

	var msg string
	flag.StringVar(&msg, "m", "", "Message of the item to add.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
	}

	if msg == "" {
		return fmt.Errorf("message not passed")
	}

	path := "todo.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("opening file[%s]: %w", path, err)
	}

	var todo Todo
	if err := json.Unmarshal(content, &todo); err != nil {
		return fmt.Errorf("unmarshaling content[%s]: %w", content, err)
	}

	todo.Items = append(todo.Items, Item{
		Name:    name,
		Message: msg,
	})

	b, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}
