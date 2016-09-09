package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/karasz/go2ban/jail"
	"github.com/karasz/go2ban/utils"
	"github.com/naoina/toml"
)

type config struct {
	Jaildir string
}

func init() {
	utils.TrapSignals()
}

func main() {
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
		go jail.Run()
	}

	//for now we just sleep forever
	select {}
}
