package jail

type Jail struct {
	LogReader *LogReader
	Enabled   bool
}

func NewJail(filename string) *Jail {
	return &Jail{
		Logreader: NewLogReader(filename),
		Enabled:   enabled,
	}
}
