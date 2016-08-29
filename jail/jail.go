package jail

import (
	"github.com/karasz/go2ban/common"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
)

type configJail struct {
	Name        string
	LogFile     string
	Regexp      []string
	MaxFail     int
	TimeVal     int
	ActionBan   string
	ActionUnBan string
	Enabled     bool
}

type Jail struct {
	Name        string
	LogFile     string
	Regexp      []string
	MaxFail     int
	TimeVal     int
	ActionBan   string
	ActionUnBan string
	Enabled     bool
	logreader   *LogReader
}

func NewJail(jailfile string) *Jail {
	f, err := os.Open(jailfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config configJail
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}

	return &Jail{
		Name:        common.Basename(jailfile),
		logreader:   NewLogReader(config.LogFile),
		LogFile:     config.LogFile,
		Regexp:      config.Regexp,
		MaxFail:     config.MaxFail,
		TimeVal:     config.TimeVal,
		ActionBan:   config.ActionBan,
		ActionUnBan: config.ActionUnBan,
		Enabled:     config.Enabled,
	}
}
