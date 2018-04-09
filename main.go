package main

import (
	"fmt"
	"github.com/gherlein/xb"
	. "github.com/gherlein/xbevents"
)

var (
	xbe   *XBevent
	debug bool = true
)

func init() {
	xb.Open()
}

func main() {

	for {
		xbe := xb.GetEvent()
		if xbe == nil {
			continue
		}
		if debug {
			fmt.Println(xbe.Name)
		}

	}

}
