# lime-local-db 🍋‍🟩

Limedb is a local wrapper for a sqlite database. It is designed to make interacting
with the database easier, to eventually provide similar functionality to Notion or Obsidian,
but with more flexibility.


## Limedb table of contents
- [Getting started](#getting-started)
- [Usage](#usage)
- [Commands](#commands)
- [Global flags](#global-flags)
- [Configuration](#configuration)
- [Features in progress](#features-in-progress)


## Getting started


## Usage


## Commands

#### Create a new database:
```sh
limedb db new [DB-name] [Key flags][Not null]:[Column type]:[Column name]{[Default value]}
# Eg:
limedb db new pets P:TEXT:name{default} N:TEXT:gender{F} N::breed{Dog}
limedb db new pets P:TEXT:name{default} N:TEXT:gender{F} N:INT:age :REAL:height_cm
```

#### Remove a database:
By default, this will soft delete a database.
```sh
limedb db rm [DB-NAME]
# Eg:
limedb db rm pets
```

#### Add an entry to a database:
```sh
limedb add-entry [column names and values]...
# Eg:
limedb --db pets add-entry name{Woofy} age{5} gender{M} breed{Beagle}
```

#### Remove entries from a database:
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

#### Create a backup of a database:
This will create a backup of the database in its current state.
```sh
limedb backup new
# Eg:
limedb --db pets backup new
limedb --db pets backup new -m "Backup before risky operation"
```

#### List all backups of a database:
```sh
limedb backup ls
# Eg:
limedb --db pets backup ls
```

#### Remove a backup of a database:
```sh
limedb backup rm [Backup id]
# A backup id can be obtained through the `limedb backup ls` command. Eg:
limedb --db pets backup ls
limedb --db pets backup rm 1
```

#### Restore to a database backup:
This will restore the selected database to its state at the time of the backup.
```sh
limedb backup restore [Backup id]
# A backup id can be obtained through the `limedb backup ls` command. Eg:
limedb --db pets backup ls
limedb --db pets backup restore 1
```


## Global flags

| Flag                | Action                                                                        |
| :------------------ | :---------------------------------------------------------------------------- |
| `--db`              | Provide the database that operations will operate on.                         |
| `--with-config`     | Provide comma separated list of configuration options.                        |

#### Example
```sh
# Use --db to select the `pets` database.
limedb --db pets add-entry name{Woofy} age{5} gender{M}

# --with-config to set the `soft_deletion` and `limedb_home` options in the current run of limedb.
limedb db rm pets --confirm --with-config soft_deletion:false,limedb_home:/etc/limedb/
```


## Configuration

| Config option       | Action                                                      | Valid args                           |
| :------------------ | :-----------------------------------------------------------|:------------------------------------ |
| `limedb_home`       | Set the home directory limedb. (Default `$HOME/.limedb/`).  |                                      |
| `log_mode`          | Set how logs should be reported. (Default `file`).          | `file`, `stdout`, `stderr`           |
| `log_level`         | The level of logging. (Default `info`).                     | `info`, `warn`, `debug`, `error`     |
| `soft_deletion`     | Soft delete databases. (Default `true`).                    | `true`, `false`                      |
| `default_db`        | Default db used for operations, (Default None).             |                                      |


## Features in progress
- Proper support for foreign keys
- Add option in where clauses to input raw sqlite
- Operations
    - alter
    - restore soft deleted db and update master db
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
