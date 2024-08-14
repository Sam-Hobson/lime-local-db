package types

type TriggerType byte

func (c TriggerType) String() string {
	return triggerTypeName[c]
}

const (
	BeforeInsert    = TriggerType(0)
	AfterInsert     = TriggerType(1)
	BeforeUpdate    = TriggerType(2)
	AfterUpdate     = TriggerType(3)
	BeforeDelete    = TriggerType(4)
	AfterDelete     = TriggerType(5)
	InsteadOfInsert = TriggerType(6)
	InsteadOfDelete = TriggerType(7)
	InsteadOfUpdate = TriggerType(8)
)

var triggerTypeName = [...]string{
	"BEFORE INSERT",
	"AFTER INSERT",
	"BEFORE UPDATE",
	"AFTER UPDATE",
	"BEFORE DELETE",
	"AFTER DELETE",
	"INSTEAD OF INSERT",
	"INSTEAD OF DELETE",
	"INSTEAD OF UPDATE",
}
