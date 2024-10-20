package parameters

import (
	"fmt"
	"handson/bot/dm/loadtest"
	"os"
	"strconv"
	"strings"
)

func ParseParams(args []string) *loadtest.Params {
	p := &loadtest.Params{}
	for _, v := range args {
		list := strings.Split(v, "=")
		if len(list) != 2 {
			continue
		}
		name := list[0]
		value := list[1]
		switch name {
		case "host":
			p.Host = value
		case "bots":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Howmany = n
		case "protocol":
			p.Protocol = strings.ToUpper(value)
		case "size":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Size = n
		case "interval":
			n, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Invalid value for bots param")
				os.Exit(1)
				return nil
			}
			p.Interval = int64(n)
		}
	}
	return p
}
