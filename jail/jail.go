package jail

import (
	"fmt"
	"github.com/karasz/go2ban/common"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"regexp"
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
	Regexp      []*regexp.Regexp
	MaxFail     int
	TimeVal     int
	ActionBan   string
	ActionUnBan string
	Enabled     bool
	logreader   *logReader
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

	rg := make([]*regexp.Regexp, 0)
	for _, v := range config.Regexp {
		rr := regexp.MustCompile(v)
		rg = append(rg, rr)

	}
	return &Jail{
		Name:        common.Basename(jailfile),
		logreader:   newLogReader(config.LogFile),
		LogFile:     config.LogFile,
		Regexp:      rg,
		MaxFail:     config.MaxFail,
		TimeVal:     config.TimeVal,
		ActionBan:   config.ActionBan,
		ActionUnBan: config.ActionUnBan,
		Enabled:     config.Enabled,
	}
}

func (j *Jail) Run() {
	if j.Enabled {
	loop:
		for {
			j.logreader.readLine()
			select {
			case _ = <-j.logreader.errors:
				break loop
			case z := <-j.logreader.lines:
				if j.matchLine(z) {
					fmt.Println(z)
				}
			}
		}
	}
}

func (j *Jail) matchLine(line string) bool {
	for _, z := range j.Regexp {
		if z.MatchString(line) {
			return true
		}
	}
	return false
}
