# lime-local-db

Limedb is a local wrapper for a sqlite database. It is designed to make interacting
with the database easier, to eventually provide similar functionality to Notion or Obsidian,
but with more flexibility.


## Commands


Create a new database:
```sh
limedb [DB-name] [Key flags][Not null]:[Column type]:[Column name]{[Default value]}
# Eg:
limedb new-db pets P:TEXT:name{default} N:TEXT:gender{F} N::breed{Dog}
limedb new-db pets P:TEXT:name{default} N:TEXT:gender{F} N:INT:age :REAL:height_cm
```

Remove a database:
```sh
limedb rm-db [DB-NAME]
# Eg:
limedb rm-db pets
```

Add an entry to a database:
```sh
limedb add-entry [column names and values]...
# Eg:
limedb --db pets add-entry name{Woofy} age{5} gender{M} breed{Beagle}
```

Remove entries from a database:
```sh
limedb rm-entries-all
limedb rm-entries-where [Column name]:[Operation]{[Value]}
# Remove all database entries
limedb --db pets rm-entries-all
# Remove database entries conditionally
limedb --db pets rm-entries-where name:like{W%}
# Adding multiple conditions will logically AND them together
limedb --db pets rm-entries-where "age:>{5}" gender:!={F} height_cm:between{10:30} name:null
```

Create a backup of a database:
```sh
limedb backup new [Table name]
# Eg:
limedb backup new pets
limedb backup new pets -m "Backup before risky operation"
```


List all backups of a database:
```sh
limedb backup ls [Table name]
# Eg:
limedb backup ls pets
```


## Flags

Provide the database that operations will operate on:
```sh
limedb --db [db name] [Any command]
# Eg:
limedb --db pets add-entry name{Woofy} age{5} gender{M}
```

Provide comma separated list of configuration options:
```sh
limedb [Any command] --with-config key:value
# Eg:
limedb [Any command] --with-config soft_deletion:false
limedb rm-db pets --confirm --with-config soft_deletion:false,limedb_home:/etc/limedb/
```

## Features to add
- Proper support for foreign keys
- Add option in where clauses to input raw sqlite
- Operations
    - Backup db
    - Restore db
    - alter
    - List databases and properties
    - Command templates
    - Sub document creation (eg. Adding a html file to a row)
    - (Maybe) Dynamic Cobra subcommands
- More complex deletions, eg extracting a subset of data into a new table
- Add soft deletion for rows
- Command templates/aliases
- Password authentication for sqlite db
- HTMX front end interface
- lua scripts to execute on db hooks
