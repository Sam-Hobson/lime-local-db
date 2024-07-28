# lime-local-db

Limedb is a local wrapper for a sqlite database. It is designed to make interacting
with the database easier, to eventually provide similar functionality to Notion or Obsidian,
but with more flexibility.


## Commands


Create a new database:
```sh
limedb [DB-name] [Key flags][Not null]:[Column type]:[Column name]{[Default value]}
# Eg:
limedb new-db petdb P:TEXT:name{default} N:TEXT:gender{F} N::breed{Dog}
limedb new-db petdb P:TEXT:name{default} N:TEXT:gender{F} N:INT:age :REAL:height_cm
```

Remove a database:
```sh
limedb rm-db [DB-NAME]
# Eg:
limedb rm-db petdb
```


## Flags

Provide comma separated list of configuration options:
```sh
limedb [Any command] --with-config key:value
# Eg:
limedb [Any command] --with-config softDeletion:false
limedb rm-db petdb --confirm --with-config softDeletion:false,limedbHome:/etc/limedb/
```

## Features to add
- Proper support for foreign keys
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
