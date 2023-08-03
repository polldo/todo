# Project
Simple todo cli

## Basic commands:

### Add a new item
todo add -m "Something todo here" -n "some"

### Remove a item
todo rm -n "some"

### Fetch all items
todo ls


## Storage
Where do we want to store our items?
- filesystem: \*json

Let's start by using a `todo.json` file placed in the same directory the script is executed.

