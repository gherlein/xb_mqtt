package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gherlein/xb"
	. "github.com/gherlein/xbevents"
	"math"
)

var (
	support_xy      bool   = true
	support_vector  bool   = false
	support_buttons bool   = true
	debugraw        bool   = false
	debugvector     bool   = false
	debugjoy        bool   = false
	debugbutton     bool   = false
	broker          string = "tcp://rpisoar:1883"
	xy              string = "xb/1/joy-xy"
	vector          string = "xb/1/joy-vector"
	buttons         string = "xb/1/buttons"
	qos             int    = 0
	xmult           int16  = 1
	ymult           int16  = -1
	client          MQTT.Client
)

func init() {
	xb.Open()
}

func main() {

	var xbe *XBevent
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	client = MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for {
		xbe = xb.GetEvent()
		if xbe == nil {
			continue
		}
		xbe.X *= xmult
		xbe.Y *= ymult
		if debugraw {
			fmt.Printf("%s  X: %d   Y: %d\n", xbe.Name, xbe.X, xbe.Y)
		}

		if xbe.Code == LJOYX || xbe.Code == LJOYY || xbe.Code == RJOYX || xbe.Code == RJOYY {
			send_joystick(xbe)
			if support_vector {
				send_vector(xbe)
			}

		} else {
			if support_buttons {
				send_button(xbe)
			}
		}
	}
}

func send_button(xbe *XBevent) {
	var msg string = "*"
	msg = fmt.Sprintf("%s", xbe.Name)
	if debugbutton {
		fmt.Println(msg)
	}
	token := client.Publish(buttons, byte(qos), false, msg)
	token.Wait()
}

func send_joystick(xbe *XBevent) {
	var msg string = "*|*|*"
	if xbe.Code == LJOYX {
		msg = fmt.Sprintf("L|X|%d|%d", xbe.X, xbe.Y)
	}
	if xbe.Code == LJOYY {
		msg = fmt.Sprintf("L|Y|%d|%d", xbe.X, xbe.Y)
	}
	if xbe.Code == RJOYX {
		msg = fmt.Sprintf("R|X|%d|%d", xbe.X, xbe.Y)
	}
	if xbe.Code == RJOYY {
		msg = fmt.Sprintf("R|Y|%d|%d", xbe.X, xbe.Y)
	}
	/* this will set the PWM width - only here for testing
	 v := float64(0.05) * p
	v = 0.15 + v
	msg := fmt.Sprintf("%d=%f", m1pin, v)
	*/
	if support_xy {
		if debugjoy {
			fmt.Printf("%s\n", msg)
		}
		token := client.Publish(xy, byte(qos), false, msg)
		token.Wait()
	}
}

func send_vector(xbe *XBevent) {
	var msg string = "*|*|*|*"
	var r, theta float64

	/*
		if xbe.X > 0 {
			lang = math.Atan2(ply, plx) * (180 / math.Pi)
			rang = math.Atan2(pry, prx) * (180 / math.Pi)
		}
		if xbe.X < 0 && xbe.Y >= 0 {
			lang = (math.Atan2(ply, plx) + math.Pi) * (180 / math.Pi)
			rang = (math.Atan2(pry, prx) + math.Pi) * (180 / math.Pi)
		}
		if xbe.X < 0 && xbe.Y > 0 {
			lang = (math.Atan2(ply, plx) - math.Pi) * (180 / math.Pi)
			rang = (math.Atan2(pry, prx) - math.Pi) * (180 / math.Pi)
		}
		if xbe.X == 0 && xbe.Y > 0 {
			lang = (math.Atan2(ply, plx) - math.Pi) * (180 / math.Pi)
			rang = (math.Atan2(pry, prx) - math.Pi) * (180 / math.Pi)
		}
		if xbe.X == 0 && xbe.Y > 0 {
			lang = (math.Pi / 2) * (180 / math.Pi)
			rang = (math.Pi / 2) * (180 / math.Pi)
		}
		if xbe.X == 0 && xbe.Y < 0 {
			lang = -1 * (math.Pi / 2) * (180 / math.Pi)
			rang = -1 * (math.Pi / 2) * (180 / math.Pi)
		}
	*/
	fx := float64(xbe.X)
	fy := float64(xbe.Y)
	rad := math.Atan2(float64(xbe.Y), float64(xbe.X))
	theta = rad * (180 / math.Pi)

	r = math.Sqrt((fx * fx) + (fy * fy))
	if xbe.Code == LJOYX || xbe.Code == LJOYY {
		msg = fmt.Sprintf("LJOY|A|%d|%d||%f|%f", xbe.X, xbe.Y, r, theta)
	}
	if xbe.Code == RJOYX || xbe.Code == RJOYY {
		msg = fmt.Sprintf("RJOY|A|%d|%d||%f|%f", xbe.X, xbe.Y, r, theta)
	}
	if debugvector {
		fmt.Println(msg)
	}
	token := client.Publish(vector, byte(qos), false, msg)
	token.Wait()
}
