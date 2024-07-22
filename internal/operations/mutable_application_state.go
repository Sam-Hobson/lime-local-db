package operations

type applicationState struct {
	selectedDatabase string
	selectedColumn   string
}

func (mas *applicationState) setSelectedDatabase(db string) {
	mas.selectedDatabase = db
}

func (mas *applicationState) getSelectedDatabase() string {
	return mas.selectedDatabase
}

func (mas *applicationState) setSelectedColumn(db string) {
	mas.selectedColumn = db
}

func (mas *applicationState) getSelectedColumn() string {
	return mas.selectedColumn
}

var mutableApplicationState = &applicationState{}
