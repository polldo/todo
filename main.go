package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	help := "Please specify a subcommand:\n\t" +
		"- 'ls' to fetch all your TODOs.\n\t" +
		"- 'add' to add a new TODO.\n\t" +
		"- 'rm' to delete a TODO.\n\t" +
		"- 'update' to update a TODO.\n\t" +
		"- 'done' to mark a TODO as complete.\n\t" +
		"- 'undone' to mark a TODO as not complete.\n\t"

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
	case "rm":
		cmd = rm
	case "update":
		cmd = update
	case "done":
		cmd = done
	case "undone":
		cmd = undone
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
	High []Item `json:"high"`
	Mid  []Item `json:"mid"`
	Low  []Item `json:"low"`
	Done []Item `json:"done"`
}

func ls(args []string) error {
	flag := flag.NewFlagSet("todo ls", flag.ExitOnError)

	var all bool
	flag.BoolVar(&all, "a", false, "Show also completed TODOs.")

	flag.Parse(args)

	path := "todo.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("opening file[%s]: %w", path, err)
	}

	var todo Todo
	if err := json.Unmarshal(content, &todo); err != nil {
		return fmt.Errorf("unmarshaling content[%s]: %w", content, err)
	}

	fmt.Println("High:")
	for _, i := range todo.High {
		fmt.Printf("%s:    %s\n", i.Name, i.Message)
	}
	fmt.Println()

	fmt.Println("Mid:")
	for _, i := range todo.Mid {
		fmt.Printf("%s:    %s\n", i.Name, i.Message)
	}
	fmt.Println()

	fmt.Println("Low:")
	for _, i := range todo.Low {
		fmt.Printf("%s:    %s\n", i.Name, i.Message)
	}
	fmt.Println()

	if all {
		fmt.Println("Done:")
		for _, i := range todo.Done {
			fmt.Printf("%s:    %s\n", i.Name, i.Message)
		}
	}

	return nil
}

func contains(items []Item, name string) bool {
	for _, it := range items {
		if it.Name == name {
			return true
		}
	}
	return false
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

	if contains(todo.Low, name) || contains(todo.Mid, name) || contains(todo.High, name) || contains(todo.Done, name) {
		return errors.New("this task exists already")
	}

	todo.Low = append(todo.Low, Item{
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

func update(args []string) error {
	flag := flag.NewFlagSet("todo update", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to update.")

	var msg string
	flag.StringVar(&msg, "m", "", "Updated message.")

	var pri string
	flag.StringVar(&pri, "p", "", "Updated priority (high - mid - low).")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
	}

	if msg == "" && pri == "" {
		return nil
	}

	if pri != "" && pri != "high" && pri != "mid" && pri != "low" {
		return fmt.Errorf("invalid priority, choose one of: 'high', 'mid' or 'low'")
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

	pop := func(items []Item, name string) (Item, []Item) {
		item := Item{}
		res := make([]Item, 0, len(items))

		for _, it := range items {
			if it.Name == name {
				item = it
				continue
			}

			res = append(res, Item{
				Name:    it.Name,
				Message: it.Message,
			})
		}
		return item, res
	}

	var item Item
	var temp Item
	temp, todo.Low = pop(todo.Low, name)
	if temp.Name != "" {
		item = temp
	}
	temp, todo.Mid = pop(todo.Mid, name)
	if temp.Name != "" {
		item = temp
	}
	temp, todo.High = pop(todo.High, name)
	if temp.Name != "" {
		item = temp
	}

	if item.Name == "" {
		return errors.New("not found")
	}

	// Update task message.
	if msg != "" {
		item.Message = msg
	}

	switch pri {
	case "high":
		todo.High = append(todo.High, item)
	case "mid":
		todo.Mid = append(todo.Mid, item)
	case "low":
		todo.Low = append(todo.Low, item)
	default:
		return fmt.Errorf("priority[%s] not valid", pri)
	}

	b, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}

func done(args []string) error {
	flag := flag.NewFlagSet("todo done", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to complete.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
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

	pop := func(items []Item, name string) (Item, []Item) {
		item := Item{}
		res := make([]Item, 0, len(items))

		for _, it := range items {
			if it.Name == name {
				item = it
				continue
			}

			res = append(res, Item{
				Name:    it.Name,
				Message: it.Message,
			})
		}
		return item, res
	}

	var item Item
	var temp Item
	temp, todo.Low = pop(todo.Low, name)
	if temp.Name != "" {
		item = temp
	}
	temp, todo.Mid = pop(todo.Mid, name)
	if temp.Name != "" {
		item = temp
	}
	temp, todo.High = pop(todo.High, name)
	if temp.Name != "" {
		item = temp
	}

	if item.Name == "" {
		return errors.New("not found")
	}

	todo.Done = append(todo.Done, item)

	b, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}

func undone(args []string) error {
	flag := flag.NewFlagSet("todo undone", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to mark as not completed.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
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

	pop := func(items []Item, name string) (Item, []Item) {
		item := Item{}
		res := make([]Item, 0, len(items))

		for _, it := range items {
			if it.Name == name {
				item = it
				continue
			}

			res = append(res, Item{
				Name:    it.Name,
				Message: it.Message,
			})
		}
		return item, res
	}

	var item Item
	item, todo.Done = pop(todo.Done, name)
	if item.Name == "" {
		return errors.New("not found")
	}

	todo.High = append(todo.High, item)

	b, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}

func rm(args []string) error {
	flag := flag.NewFlagSet("todo rm", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to remove.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
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

	filter := func(items []Item, name string) []Item {
		res := make([]Item, 0, len(items))
		for _, it := range items {
			if it.Name == name {
				continue
			}

			res = append(res, Item{
				Name:    it.Name,
				Message: it.Message,
			})
		}
		return res
	}

	todo.Low = filter(todo.Low, name)
	todo.Mid = filter(todo.Mid, name)
	todo.High = filter(todo.High, name)
	todo.Done = filter(todo.Done, name)

	b, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}
