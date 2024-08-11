package masterdatabase

import (
	"database/sql"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"
	dbutil "github.com/sam-hobson/internal/database/util"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/types"
	"github.com/sam-hobson/internal/util"
)

const (
	masterDatabaseName     = ".master"
	masterDatabaseFileName = ".master.db"
)

var tablesColumns = []*types.Column{
	{
		Name:       "name",
		DataType:   types.ColumnTextDataType,
		NotNull:    true,
		PrimaryKey: true,
	},
	{
		Name:     "created",
		DataType: types.ColumnTextDataType,
		NotNull:  true,
	},
	{
		Name:     "softdeleted",
		DataType: types.ColumnIntDataType,
	},
	{
		Name:     "harddeleted",
		DataType: types.ColumnIntDataType,
	},
}

func QueryTables(where *sqlbuilder.WhereClause, columnNames ...string) (*sql.Rows, error) {
    util.Log("860b9c05").Info("Querying master database.", "Where", where)

	// This will do nothing if the table already exists
	createTablesTable()

	db, err := openMasterDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()


    sb := sqlbuilder.NewSelectBuilder().Select(columnNames...).From("tables").AddWhereClause(where)
    selStr, args := sb.Build()

    util.Log("e80c1d0a").Info("Querying master database with SQL.", "SQL", selStr, "Args", args)

	res, err := db.Query(selStr, args...)
	if err != nil {
		util.Log("10be1def").Warn("Could not query master database.")
		return nil, err
	}

    return res, nil
}

func AddNewDatabaseRecord(data map[string]string) error {
	util.Log("ba44622e").Info("Adding new database record to master database.", "Data", data)

	// This will do nothing if the table already exists
	createTablesTable()

	db, err := openMasterDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	insertStr, args := dbutil.InsertIntoTableSql("tables", data)

	util.Log("e1d6d757").Info("Inserting with SQL Command.", "SQL", insertStr, "Args", args)

	if _, err = db.Exec(insertStr, args...); err != nil {
		util.Log("e298f78c").Error("Failed executing insert table command.", "SQL", insertStr, "Args", args)
		return err
	}

	util.Log("3cc23989").Info("Successfully inserted database record into master database.", "Data", data)

	return nil
}

func RemoveDatabaseRecord(databaseName string, softDelete bool) error {
	util.Log("0b3a756a").Info("Setting database to deleted.", "Database name", databaseName)

	// This will do nothing if the table already exists
	createTablesTable()

	db, err := openMasterDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	cond := sqlbuilder.NewCond()
	where := sqlbuilder.NewWhereClause()
	where.AddWhereExpr(cond.Args, cond.Equal("name", databaseName))

	ub := sqlbuilder.NewUpdateBuilder().Update("tables")
	if softDelete {
		ub.Set(ub.Assign("softdeleted", 1))
	} else {
		ub.Set(ub.Assign("harddeleted", 1))
	}

	ub.AddWhereClause(where)
	sql, args := ub.Build()
	util.Log("8fc6aa55").Info("Setting master record for database to deleted.", "SQL", sql, "Args", args)

	res, err := db.Exec(sql, args...)
	if err != nil {
		util.Log("8f3a39af").Warn("Failed setting database record in master database to deleted.", "SQL", sql, "Args", args)
		return err
	}
	if util.PanicIfErr(res.RowsAffected()) != 1 {
		util.Log("4f37a3d4").Warn("Failed setting database record in master database to deleted.", "SQL", sql, "Args", args)
		return errors.Errorf("Failed removing database record from master database")
	}

	util.Log("377df20f").Info("Successfully set master database to deleted.", "Database name", databaseName)

	return nil
}

func createTablesTable() error {
	util.Log("91f36c46").Info("Creating tables table if doesn't exist in master database.")

	db, err := openMasterDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	createTableStr, args := dbutil.CreateTableSql("tables", tablesColumns)

	util.Log("543593c8").Info("Creating table with SQL command.", "SQL", createTableStr, "Args", args)

	if _, err = db.Exec(createTableStr, args...); err != nil {
		util.Log("9c945589").Error("Failed executing create table command.", "SQL", createTableStr, "Args", args)
		return err
	}

	return nil
}

func openMasterDatabase() (*sql.DB, error) {
	if db, err := sql.Open("sqlite3", MasterDatabaseFilePath()); err != nil {
		util.Log("facf46fb").Warn("Could not open master database.")
		return nil, err
	} else {
		return db, nil
	}
}

func MasterDatabaseFilePath() string {
	return filepath.Join(state.ApplicationState().GetLimedbHome(), masterDatabaseFileName)
}
