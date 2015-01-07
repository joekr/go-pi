package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func work(robot *gobot.Robot) {
	gobot.On(robot.Device("button").(*gpio.ButtonDriver).Event("push"), func(data interface{}) {
		robot.Device("led").(*gpio.LedDriver).On()
	})
	gobot.On(robot.Device("button").(*gpio.ButtonDriver).Event("release"), func(data interface{}) {
		robot.Device("led").(*gpio.LedDriver).Off()
	})

}

func goFunc(){
	gbot := gobot.NewGobot()
	
	r := raspi.NewRaspiAdaptor("raspi")
	led := gpio.NewLedDriver(r, "led", "37")
	button := gpio.NewButtonDriver(r, "button", "36")

	work := func() {
		work(gbot.Robot("buttonBot"))
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{led, button},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

func main() {

	goFunc()

}
