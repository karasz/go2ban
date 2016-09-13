package go2banmain

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/karasz/go2ban"
	"github.com/karasz/go2ban/jail"
	"github.com/naoina/toml"
)

type config struct {
	Jaildir string
}

func Run() {
	var conffile string
	jailfiles := make([]string, 0)
	flag.StringVar(&conffile, "f", "/etc/go2ban/g2b.toml", "specify the config file.  defaults to /etc/go2ban/g2b.toml.")
	flag.Parse()
	f, err := os.Open(conffile)
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config config
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	f.Close()

	d, err := os.Open(config.Jaildir)
	if err != nil {
		panic(err)
	}

	fls, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	d.Close()

	for _, file := range fls {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".g2b" {
				jailfiles = append(jailfiles, config.Jaildir+file.Name())
			}
		}
	}

	for _, j := range jailfiles {
		jail := jail.NewJail(j)
		go2ban.Srv.Js = append(go2ban.Srv.Js, jail)
		go jail.Run()
	}

	go2ban.TrapSignals()

	//for now we just sleep forever
	select {}
}
