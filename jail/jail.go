package jail

import (
	"fmt"
	"github.com/karasz/go2ban/common"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type configJail struct {
	Name        string
	LogFile     string
	TimeFormat  string
	Regexp      []string
	MaxFail     int
	BanTime     int
	FindTime    int
	ActionBan   string
	ActionUnBan string
	ActionSetup string
	Enabled     bool
}

type Jail struct {
	Name        string
	LogFile     string
	TimeFormat  string
	Regexp      []*regexp.Regexp
	MaxFail     int
	BanTime     int
	FindTime    int
	ActionBan   string
	ActionUnBan string
	ActionSetup string
	Enabled     bool
	logreader   *logReader
	jailees     []*jailee
}

type jailee struct {
	ip        string
	failcount int
}

func (j *Jail) getJailee(ip string) (int, *jailee, bool) {
	for i, ja := range j.jailees {
		if ja.ip == ip {
			return i, ja, true
		}
	}
	return -1, nil, false
}

func (j *Jail) Add(ip string) {
	if _, jj, ok := j.getJailee(ip); !ok {
		j.jinit(ip)
	} else {
		jj.failcount++
	}
}

func (j *Jail) jinit(ip string) {
	ja := jailee{failcount: 1, ip: ip}
	j.jailees = append(j.jailees, &ja)
}

func (j *Jail) check(ip string) {
	if _, jj, ok := j.getJailee(ip); ok {
		if jj.failcount == j.MaxFail {
			go j.executeBan(jj)
		}
	}
}

func (j *Jail) checkFind(toCheck string) bool {
	to, _ := time.Parse(j.TimeFormat, toCheck)
	if to.Year() == 0 {
		nowY := time.Now().Year()
		to = to.AddDate(nowY, 0, 0)
	}
	if time.Since(to) > time.Duration(j.FindTime)*time.Minute {
		return false
	}
	return true
}

func (j *Jail) executeBan(jj *jailee) {
	fmt.Println(j.ActionBan)

	timer := time.NewTimer(time.Duration(j.BanTime) * time.Minute)
	<-timer.C
	j.executeUnBan(jj)
}

func (j *Jail) executeUnBan(jj *jailee) {
	fmt.Println(j.ActionUnBan)

}

func (j *Jail) executeSetup() {
	fmt.Println(j.ActionSetup)
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
	j := Jail{
		Name:        common.Basename(jailfile),
		logreader:   newLogReader(config.LogFile),
		LogFile:     config.LogFile,
		TimeFormat:  config.TimeFormat,
		Regexp:      rg,
		MaxFail:     config.MaxFail,
		BanTime:     config.BanTime,
		FindTime:    config.FindTime,
		ActionBan:   config.ActionBan,
		ActionUnBan: config.ActionUnBan,
		ActionSetup: config.ActionSetup,
		Enabled:     config.Enabled,
		jailees:     make([]*jailee, 0),
	}
	j.executeSetup()
	return &j
}

func (j *Jail) Run() {
	if j.Enabled {
		for {
			j.logreader.readLine()
			select {
			case <-j.logreader.errors:
				j.logreader.reset()
			case z := <-j.logreader.lines:
				if q, ok := j.matchLine(z); ok {
					if j.checkFind(q["DATETIME"]) {
						j.Add(q["HOST"])
						j.check(q["HOST"])
					}
				}
			}
		}
	}
}

func (j *Jail) matchLine(line string) (map[string]string, bool) {
	result := make(map[string]string)
	for _, z := range j.Regexp {
		match := z.FindStringSubmatch(line)
		if match != nil {
			for i, name := range z.SubexpNames() {
				if i == 0 || name == "" {
					continue
				}
				result[name] = match[i]
			}
			return result, true
		}
	}
	return result, false
}
