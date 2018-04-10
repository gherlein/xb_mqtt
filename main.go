package main

import (
	"fmt"
	"github.com/gherlein/xb"
	. "github.com/gherlein/xbevents"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

var (
	xbe   *XBevent
	debug bool = true

	network string = "tcp"
	url     string = "rpisoar:1883"
)

func init() {
	xb.Open()
}

func main() {

	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	err := cli.Connect(&client.ConnectOptions{
		Network:  network,
		Address:  url,
		ClientID: []byte("example-client"),
	})
	if err != nil {
		panic(err)
	}

	for {
		xbe := xb.GetEvent()
		if xbe == nil {
			continue
		}

		if xbe.Code == LJOYX || xbe.Code == LJOYY || xbe.Code == RJOYX || xbe.Code == RJOYY {
			if debug {
				fmt.Printf("%s  x: %d   y: %d\n", xbe.Name, xbe.X, xbe.Y)
			}

			msg := fmt.Sprintf("%s|%d|%d", xbe.Name, xbe.X, xbe.Y)
			err = cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte("xb/1/joy"),
				Message:   []byte(msg),
			})
			if err != nil {
				panic(err)
			}

		} else {

			msg := fmt.Sprintf("%s", xbe.Name)
			err = cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte("xb/1/raw"),
				Message:   []byte(msg),
			})
			if err != nil {
				panic(err)
			}
			if debug {
				fmt.Println(xbe.Name)
			}
		}

	}

}
