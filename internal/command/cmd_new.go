package command

type newDbCmd struct {
	err error
}

func (ps *newDbCmd) error() error {
	return ps.err
}

func (ps *newDbCmd) process(key string) argProcessor {
	return ps
}

func newNewDbCmd() *newDbCmd {
	return &newDbCmd{}
}
