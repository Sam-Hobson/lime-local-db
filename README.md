# lime-local-db

Limedb is a local wrapper for an sqlite database. It is designed to make interacting
with the database easier, to eventually provide similar functionality to Notion or Obsidian,
but with more flexibility.


## Commands


Create default configuration file at ~/.limerc:
```sh
limedb setup
```


Create a new database:
```sh
limedb new-db [DB-NAME] "[<col1> <col2> <col3>]"
# Eg:
limedb new-db Dogs "[Name Gender Breed]"
```


Remove a database:
```sh
limedb rm-db [DB-NAME]
# Eg:
limedb rm-db Dogs
```

## Features to add
- Operations
    - alter
    - insert
    - Command templates
    - Sub document creation (eg. Adding a html file to a row)
    - with configuration flag
    - (Maybe) Dynamic Cobra subcommands.
- Command templates/aliases
- Password authentication for sqlite db.
- HTMX front end interface
- lua scripts to execute on db hooks
