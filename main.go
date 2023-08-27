package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func main() {
	help := "Please specify a subcommand:\n\t" +
		"- 'ls' to fetch tasks.\n\t" +
		"- 'add' to add a new task.\n\t" +
		"- 'rm' to delete a task.\n\t" +
		"- 'update' to update a task.\n\t" +
		"- 'done' to mark a task as complete.\n\t"

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
	default:
		fmt.Println(help)
		os.Exit(1)
	}

	if err := cmd(os.Args[2:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func path(dir string) string {
	// Default to current directory.
	const def = "todo.json"

	// Env var has precedence.
	if p, ok := os.LookupEnv("TODO_DIR"); ok {
		return filepath.Join(p, def)
	}

	// If not specified use the provided dir.
	if dir != "" {
		return filepath.Join(dir, def)
	}

	return def
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

func fromFile(path string) (Todo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Todo{}, fmt.Errorf("opening file[%s]: %w", path, err)
	}

	var todo Todo
	if err := json.Unmarshal(content, &todo); err != nil {
		return Todo{}, fmt.Errorf("unmarshaling content[%s]: %w", content, err)
	}

	return todo, nil
}

func toFile(path string, todo Todo) error {
	b, err := json.MarshalIndent(todo, "", "    ")
	if err != nil {
		return fmt.Errorf("marshaling new content[%v]: %w", todo, err)
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("writing new content: %w", err)
	}

	return nil
}

func ls(args []string) error {
	flag := flag.NewFlagSet("todo ls", flag.ExitOnError)

	var all bool
	flag.BoolVar(&all, "a", false, "Show also completed TODOs.")

	var search string
	flag.StringVar(&search, "s", "", "Search for tasks containing specific string.")

	var dir string
	flag.StringVar(&dir, "d", "", "Path of the directory containing the 'todo.json' to read/write.")

	flag.Parse(args)

	fpath := path(dir)
	todo, err := fromFile(fpath)
	if err != nil {
		return fmt.Errorf("read todo from file: %w", err)
	}

	filter := func(items []Item, search string) []Item {
		res := make([]Item, 0, len(items))
		for _, it := range items {
			if !strings.Contains(it.Name, search) && !strings.Contains(it.Message, search) {
				continue
			}

			res = append(res, Item{
				Name:    it.Name,
				Message: it.Message,
			})
		}
		return res
	}

	print := func(items []Item, head string, colors []color.Attribute) {
		max := 0
		for _, it := range items {
			if len(it.Name) > max {
				max = len(it.Name)
			}
		}

		fmt.Println(head)
		for i, it := range items {
			c := colors[i%len(colors)]
			tab := strings.Repeat(" ", max+4-len(it.Name))
			color.New(c).Printf("%s:%s%s\n", it.Name, tab, it.Message)
		}
		fmt.Println()
	}

	todo.High = filter(todo.High, search)
	if len(todo.High) > 0 {
		print(todo.High, "High:", []color.Attribute{color.FgRed, color.FgHiMagenta})
	}

	todo.Mid = filter(todo.Mid, search)
	if len(todo.Mid) > 0 {
		print(todo.Mid, "Mid:", []color.Attribute{color.FgGreen, color.FgHiGreen})
	}

	todo.Low = filter(todo.Low, search)
	if len(todo.Low) > 0 {
		print(todo.Low, "Low:", []color.Attribute{color.FgBlue, color.FgHiBlue})
	}

	todo.Done = filter(todo.Done, search)
	if all && len(todo.Done) > 0 {
		print(todo.Done, "Done:", []color.Attribute{color.Underline})
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

	var prio string
	flag.StringVar(&prio, "p", "low", "Priority of the item to add (high - mid - low).")

	var dir string
	flag.StringVar(&dir, "d", "", "Path of the directory containing the 'todo.json' to read/write.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
	}

	if msg == "" {
		return fmt.Errorf("message not passed")
	}

	if prio != "" && prio != "high" && prio != "mid" && prio != "low" {
		return fmt.Errorf("invalid priority, choose one of: 'high', 'mid' or 'low'")
	}

	fpath := path(dir)
	todo, err := fromFile(fpath)
	if err != nil {
		return fmt.Errorf("read todo from file: %w", err)
	}

	if contains(todo.Low, name) || contains(todo.Mid, name) || contains(todo.High, name) {
		return errors.New("this task exists already")
	}

	item := Item{
		Name:    name,
		Message: msg,
	}

	switch prio {
	case "high":
		todo.High = append(todo.High, item)
	case "mid":
		todo.Mid = append(todo.Mid, item)
	case "low":
		todo.Low = append(todo.Low, item)
	default:
		return fmt.Errorf("priority[%s] not valid", prio)
	}

	if err := toFile(fpath, todo); err != nil {
		return fmt.Errorf("write todo to file: %w", err)
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

	var dir string
	flag.StringVar(&dir, "d", "", "Path of the directory containing the 'todo.json' to read/write.")

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

	fpath := path(dir)
	todo, err := fromFile(fpath)
	if err != nil {
		return fmt.Errorf("read todo from file: %w", err)
	}

	updateMsg := func(items []Item, name string, newMsg string) {
		for i, item := range items {
			if item.Name == name {
				items[i].Message = newMsg
			}
		}
	}

	if msg != "" {
		updateMsg(todo.Low, name, msg)
		updateMsg(todo.Mid, name, msg)
		updateMsg(todo.High, name, msg)
	}

	if pri != "" {
		var item Item
		if it, items, found := pop(todo.Low, name); found {
			item = it
			todo.Low = items
		}
		if it, items, found := pop(todo.Mid, name); found {
			item = it
			todo.Mid = items
		}
		if it, items, found := pop(todo.High, name); found {
			item = it
			todo.High = items
		}

		if item.Name == "" {
			return errors.New("not found")
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
	}

	if err := toFile(fpath, todo); err != nil {
		return fmt.Errorf("write todo to file: %w", err)
	}

	return nil
}

func pop(items []Item, name string) (Item, []Item, bool) {
	item := Item{}
	res := make([]Item, 0, len(items))
	found := false

	for _, it := range items {
		if it.Name == name {
			item = it
			found = true
			continue
		}

		res = append(res, Item{
			Name:    it.Name,
			Message: it.Message,
		})
	}
	return item, res, found
}

func done(args []string) error {
	flag := flag.NewFlagSet("todo done", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to complete.")

	var dir string
	flag.StringVar(&dir, "d", "", "Path of the directory containing the 'todo.json' to read/write.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
	}

	fpath := path(dir)
	todo, err := fromFile(fpath)
	if err != nil {
		return fmt.Errorf("read todo from file: %w", err)
	}

	var item Item
	if it, items, found := pop(todo.Low, name); found {
		item = it
		todo.Low = items
	}
	if it, items, found := pop(todo.Mid, name); found {
		item = it
		todo.Mid = items
	}
	if it, items, found := pop(todo.High, name); found {
		item = it
		todo.High = items
	}

	if item.Name == "" {
		return errors.New("not found")
	}

	todo.Done = append(todo.Done, item)

	if err := toFile(fpath, todo); err != nil {
		return fmt.Errorf("write todo to file: %w", err)
	}

	return nil
}

func rm(args []string) error {
	flag := flag.NewFlagSet("todo rm", flag.ExitOnError)

	var name string
	flag.StringVar(&name, "n", "", "Name of the item to remove.")

	var dir string
	flag.StringVar(&dir, "d", "", "Path of the directory containing the 'todo.json' to read/write.")

	flag.Parse(args)

	if name == "" {
		return fmt.Errorf("name not passed")
	}

	fpath := path(dir)
	todo, err := fromFile(fpath)
	if err != nil {
		return fmt.Errorf("read todo from file: %w", err)
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

	if err := toFile(fpath, todo); err != nil {
		return fmt.Errorf("write todo to file: %w", err)
	}

	return nil
}
