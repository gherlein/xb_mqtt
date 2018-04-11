package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gherlein/xb"
	. "github.com/gherlein/xbevents"
	"math"
)

var (
	xbe         *XBevent
	debugraw    bool   = false
	debugjoy    bool   = true
	debugbutton bool   = true
	broker      string = "tcp://rpisoar:1883"
	raw         string = "xb/1/joy-xy"
	angles      string = "xb/1/joy-vector"
	buttons     string = "xb/1/buttons"
	qos         int    = 0
	xmult       int16  = 1
	ymult       int16  = -1
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
			if debugraw {
				fmt.Printf("%s  x: %d   y: %d\n", xbe.Name, xbe.X, xbe.Y)
			}
			xbe.X *= xmult
			xbe.Y *= ymult
			var p float64 = 0.0
			var msg string = "*|*|*"
			if xbe.Code == LJOYX {
				p = float64(xbe.X) / float64(math.MaxInt16)
				msg = fmt.Sprintf("L|X|%f", p)
			}
			if xbe.Code == LJOYY {
				p = float64(xbe.Y) / float64(math.MaxInt16)
				msg = fmt.Sprintf("L|Y|%f", p)
			}
			if xbe.Code == RJOYX {
				p = float64(xbe.X) / float64(math.MaxInt16)
				msg = fmt.Sprintf("R|X|%f", p)
			}
			if xbe.Code == RJOYY {
				p = float64(xbe.Y) / float64(math.MaxInt16)
				msg = fmt.Sprintf("R|Y|%f", p)
			}
			/* this will set the PWM width - only here for testing
			 v := float64(0.05) * p
			v = 0.15 + v
			msg := fmt.Sprintf("%d=%f", m1pin, v)
			*/
			if debugjoy {
				fmt.Printf("%s\n", msg)
			}
			token := client.Publish(raw, byte(qos), false, msg)
			token.Wait()
		} else {
			var msg string = "*|*|*"
			msg = fmt.Sprintf("%s", xbe.Name)
			if debugbutton {
				fmt.Println(msg)
			}
			token := client.Publish(raw, byte(qos), false, msg)
			token.Wait()
			if xbe.Code == A_DOWN {

			}
			if xbe.Code == A_UP {

			}

		}

	}

}
