# lime-local-db

Limedb is a local wrapper for an sqlite database. It is designed to make interacting
with the database easier, to eventually provide similar functionality to Notion or Obsidian,
but with more flexibility.


## Commands

Create a new database
```sh
limedb new-db [DB_NAME]
```

## Features to add
- Operations
    - new-db
    - alter
    - insert
    - delete
- HTMX front end interface
- lua scripts to execute on db hooks
