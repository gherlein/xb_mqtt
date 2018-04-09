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
)

func init() {
	xb.Open()
}

func main() {

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "iot.eclipse.org:1883",
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

		// Publish a message.
		err = cli.Publish(&client.PublishOptions{
			QoS:       mqtt.QoS0,
			TopicName: []byte("bar/baz"),
			Message:   []byte("testMessage"),
		})
		if err != nil {
			panic(err)
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
