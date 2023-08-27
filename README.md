# Todo CLI Tool in Go
A CLI tool written in Go for managing your daily tasks. Use this tool to quickly add, remove, update, and list tasks. 

## Introduction
`todo-cli` allows you to maintain a simple todo list. It leverages a JSON file for storage, giving you the flexibility to manually edit or inspect your tasks without using the tool. By default, the tool operates on a file named 'todo.json' in the current directory. But don't worry, this is configurable!

## Priority Levels
Tasks have three priority levels: `high`, `mid`, and `low`. Assign or change a task's priority to influence the order in which they are displayed.


## Basic commands:

### Add
Add a task with a specified name and description.
```
todo add -n <name> -m "<message>"
```

Optionally, set the task's priority (`high`, `mid`, or `low`).
```
todo add -n <name> -m "<message>" -p <priority>
```

### Remove
Remove a task using its name.
```
todo rm -n <name>
```

### List
Display all tasks, orderd by priority.
```
todo ls
```

Filter tasks by a keyword.
```
todo ls -s <keyword>
```

Include also completed tasks in the list.
```
todo ls -a
```

### Update
Modify the message description of a task.
```
todo update -n <name> -m <new-message>
```

Update the priority of a task (to one of `high`, `mid`, `low`).
```
todo update -n <name> -p <priority>
```

Update both message and priority.
```
todo update -n <name> -m <message> -p <priority>
```

### Done
Mark a task as completed.
```
todo done -n <name>
```

### Specify a Different Directory
By default `todo` looks for 'todo.json' in the current directory. 
However, you can point it to another directory with a flag or with an environment variable.
```
todo ls -d <directory_path>
```

```
export TODO_DIR=<directory_path>
todo ls
```

