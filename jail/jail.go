package jail

import (
	"fmt"
	"github.com/karasz/go2ban/common"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"regexp"
	//	"sync"
	"time"
)

type configJail struct {
	Name        string
	LogFile     string
	Regexp      []string
	MaxFail     int
	BanTime     int
	ActionBan   string
	ActionUnBan string
	Enabled     bool
}

type Jail struct {
	Name        string
	LogFile     string
	Regexp      []*regexp.Regexp
	MaxFail     int
	BanTime     int
	ActionBan   string
	ActionUnBan string
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

func (j *Jail) executeBan(jj *jailee) {
	fmt.Println("Banning IP ", jj.ip)
	fmt.Println(j.ActionBan)

	timer := time.NewTimer(time.Duration(j.BanTime) * time.Minute)
	<-timer.C
	j.executeUnBan(jj)
}

func (j *Jail) executeUnBan(jj *jailee) {
	fmt.Println("UnBanning IP ", jj.ip)
	fmt.Println(j.ActionUnBan)

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
		BanTime:     config.BanTime,
		ActionBan:   config.ActionBan,
		ActionUnBan: config.ActionUnBan,
		Enabled:     config.Enabled,
		jailees:     make([]*jailee, 0),
	}
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
					j.Add(q["HOST"])
					j.check(q["HOST"])
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
