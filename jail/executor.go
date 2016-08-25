package jail

type Executor struct {
	banaction   execfn
	unbanaction execfn
}

type execfn func(string) error

func NewExecutor(commands map[string]execfn) *Executor {
	return &Executor{banaction: commands["banaction"],
		unmabaction: commands["unvanaction"],
	}
}
