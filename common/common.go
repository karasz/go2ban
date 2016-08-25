package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

const GoSock = "/var/run/go2ban/socket"

func SameLog(file string, sum string) bool {
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
