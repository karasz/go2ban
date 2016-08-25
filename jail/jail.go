package jail

type Jail struct {
	LogReader *LogReader
	Executor  *Executor
	Regex     string
	Enabled   bool
}

func NewJail(filename string, regex string, commands map[string]execfn, enable bool) *Jail {
	return &Jail{
		Logreader: NewLogReader(filename, 0),
		Regex:     regex,
		Executor:  NewExecutor(commands),
		Enabled:   enabled,
	}
}
