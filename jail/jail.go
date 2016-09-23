package jail

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/naoina/toml"
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
	Testing     bool
}

type Jail struct {
	name        string
	logFile     string
	whitelist   []string
	timeFormat  string
	regexp      []*regexp.Regexp
	maxFail     int
	banTime     int
	findTime    int
	actionBan   string
	actionUnBan string
	actionSetup string
	enabled     bool
	testing     bool
	logreader   *logReader
	jailees     []*jailee
	Cells       map[string]time.Time
}

type jailee struct {
	q         map[string]string
	ip        string
	failcount int
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
		name:        basename(jailfile),
		logreader:   newLogReader(config.LogFile),
		logFile:     config.LogFile,
		whitelist:   readwhite(strings.TrimSuffix(jailfile, ".g2b") + ".whitelist"),
		timeFormat:  config.TimeFormat,
		regexp:      rg,
		maxFail:     config.MaxFail,
		banTime:     config.BanTime,
		findTime:    config.FindTime,
		actionBan:   config.ActionBan,
		actionUnBan: config.ActionUnBan,
		actionSetup: config.ActionSetup,
		enabled:     config.Enabled,
		testing:     config.Testing,
		jailees:     make([]*jailee, 0),
		Cells:       make(map[string]time.Time),
	}
	j.executeSetup()
	return &j
}

func (j *Jail) getJailee(ip string) (int, *jailee, bool) {
	for i, ja := range j.jailees {
		if ja.ip == ip {
			return i, ja, true
		}
	}
	return -1, nil, false
}

func (j *Jail) add(q map[string]string) {
	ip := q["HOST"]
	if !isWhite(ip, j.whitelist) {
		if _, jj, ok := j.getJailee(ip); !ok {
			j.jinit(q)
		} else {
			jj.failcount++
		}
	}
}

func (j *Jail) jinit(q map[string]string) {
	ja := jailee{failcount: 1, ip: q["HOST"], q: q}
	j.jailees = append(j.jailees, &ja)
	j.Cells[q["HOST"]] = time.Now()
}

func (j *Jail) check(ip string) {
	if _, jj, ok := j.getJailee(ip); ok {
		if jj.failcount == j.maxFail {
			go j.executeBan(jj)
		}
	}
}

func (j *Jail) remove(jj *jailee) bool {
	// we have a slice and concurent access
	// so we cannot remove it hence we do a soft delete
	delete(j.Cells, jj.ip)
	jj.failcount = 0
	if jj.failcount == 0 {
		return true
	}
	return false
}

func (j *Jail) checkFind(toCheck string) bool {
	to, _ := time.Parse(j.timeFormat, toCheck)
	if to.Year() == 0 {
		nowY := time.Now().Year()
		to = to.AddDate(nowY, 0, 0)
	}
	if time.Since(to) > time.Duration(j.findTime)*time.Minute && !j.testing {
		return false
	}
	return true
}

func (j *Jail) executeSetup() {
	cmd := j.parseCommand(j.actionSetup, nil)

	if j.testing {
		prettyprint(cmd)
	} else {

		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
	}
}

func (j *Jail) executeBan(jj *jailee) {
	cmd := j.parseCommand(j.actionBan, jj)

	if j.testing {
		prettyprint(cmd)
	} else {
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
	}

	timer := time.NewTimer(time.Duration(j.banTime) * time.Minute)
	<-timer.C
	j.executeUnBan(jj)
}

func (j *Jail) executeUnBan(jj *jailee) {
	cmd := j.parseCommand(j.actionUnBan, jj)

	if j.testing {
		prettyprint(cmd)
	} else {
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
	}

	if ok := j.remove(jj); !ok {
		fmt.Println("cannot remove jailee")
	}
}

func (j *Jail) parseCommand(cmd string, jj *jailee) *exec.Cmd {
	bin := strings.Fields(cmd)[0]
	args := strings.Fields(cmd)[1:]
	if jj != nil {
		for i, k := range args {
			if strings.HasPrefix(k, "<") && strings.HasSuffix(k, ">") {
				s := strings.TrimPrefix(strings.TrimSuffix(k, ">"), "<")
				args[i] = jj.q[s]
			}
		}
	}
	c := exec.Command(bin, strings.Join(args, " "))
	c.Stdout = os.Stdout
	return c
}

func (j *Jail) Run() {
	if j.enabled {
		for {
			j.logreader.readLine()
			select {
			case <-j.logreader.errors:
				j.logreader.reset()
			case z := <-j.logreader.lines:
				if q, ok := j.matchLine(z); ok {
					if j.checkFind(q["DATETIME"]) {
						j.add(q)
						j.check(q["HOST"])
					}
				}
			}
		}
	}
}

func (j *Jail) matchLine(line string) (map[string]string, bool) {
	result := make(map[string]string)
	for _, z := range j.regexp {
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

func sameLog(file string, sum string) bool {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return false
	}

	r := bufio.NewReader(f)
	line, _, er := r.ReadLine()
	if er != nil {
		fmt.Println(er)
		return false
	}

	hash := md5.Sum(line)
	strHash := hex.EncodeToString(hash[:])

	if strHash == sum {
		return true
	}
	return false
}

func basename(s string) string {
	base := path.Base(s)
	n := strings.LastIndexByte(base, '.')
	if n >= 0 {
		return base[:n]
	}
	return base
}

func prettyprint(c *exec.Cmd) {
	var s string = "In testing mode. If not, would have executed:\n"
	s += c.Path
	s += " "
	for i, w := range c.Args {
		if i != 0 {
			s += w
			s += " "
		}
	}
	fmt.Println(s)
}

func readwhite(wf string) []string {
	result := []string{}

	whiteFile, _ := os.Open(wf)
	defer whiteFile.Close()

	scanner := bufio.NewScanner(whiteFile)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}

func isWhite(ip string, whitelist []string) bool {
	for _, k := range whitelist {
		if ip == k {
			return true
		}
	}
	return false
}
