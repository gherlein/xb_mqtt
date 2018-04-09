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
		if xbe.Code == LJOYX || xbe.Code == LJOYY || xbe.Code == RJOYX || xbe.Code == RJOYY {
			if debug {
				fmt.Printf("%s  x: %d   y: %d\n", xbe.Name, xbe.X, xbe.Y)
			}

		} else {

			if debug {
				fmt.Println(xbe.Name)
			}
		}

	}

}
