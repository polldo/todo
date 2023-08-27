# Todo CLI Tool
A CLI tool written in Go for managing your daily tasks.
Use this tool to quickly add, remove, update, and list tasks.

## Introduction
`todo` allows you to maintain a simple todo list.
It leverages a JSON file for storage, giving you the flexibility to manually edit or inspect your tasks without using the tool.
By default, the tool operates on a file named 'todo.json' in the current directory.
But don't worry, this is configurable!

## Priority Levels
Tasks have three priority levels: `high`, `mid`, and `low`. Assign or change a task's priority to influence the order in which they are displayed.


## Basic commands:

### Add
Add a task with a specified name and description.
``` bash
todo add -n <name> -m "<message>"
```

Optionally, set the task's priority (`high`, `mid`, or `low`).
``` bash
todo add -n <name> -m "<message>" -p <priority>
```

### Remove
Remove a task using its name.
``` bash
todo rm -n <name>
```

### List
Display all tasks, orderd by priority.
``` bash
todo ls
```

Filter tasks by a keyword.
``` bash
todo ls -s <keyword>
```

Include also completed tasks in the list.
``` bash
todo ls -a
```

### Update
Modify the message description of a task.
``` bash
todo update -n <name> -m <new-message>
```

Update the priority of a task (to one of `high`, `mid`, `low`).
``` bash
todo update -n <name> -p <priority>
```

Update both message and priority.
``` bash
todo update -n <name> -m <message> -p <priority>
```

### Done
Mark a task as completed.
``` bash
todo done -n <name>
```

### Specify a Different Directory
By default `todo` looks for 'todo.json' in the current directory. 
However, you can point it to another directory with a flag or with an environment variable.
``` bash
todo ls -d <directory_path>
```

``` bash
export TODO_DIR=<directory_path>
todo ls
```

