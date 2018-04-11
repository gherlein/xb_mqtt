package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gherlein/xb"
	. "github.com/gherlein/xbevents"
	"math"
)

var (
	xbe    *XBevent
	debug  bool   = true
	topic  string = "pi-blaster-mqtt/text"
	broker string = "tcp://rpisoar:1883"
	qos    int    = 0
	m1pin  int    = 4
)

func init() {
	xb.Open()
}

func main() {

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
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

			p := float64(xbe.Y) / float64(math.MaxInt16)
			v := float64(0.05) * p
			v = 0.15 + v
			msg := fmt.Sprintf("%d=%f", m1pin, v)
			fmt.Printf("%s\n", msg)

			token := client.Publish(topic, byte(qos), false, msg)
			token.Wait()
		} else {

			if xbe.Code == A_DOWN {

			}
			if xbe.Code == A_UP {

			}

			if debug {
				fmt.Println(xbe.Name)
			}
		}

	}

}
