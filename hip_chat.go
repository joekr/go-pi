package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"time"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

type Presence struct {
	Status	string `json:"status"`
	Show	string `json:"show"`
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MentionName string `json:"mention_name"`
	Presence    Presence `json:"presence"`
}

func Get(status chan<- string, user string) (rep string, err error){

	url := fmt.Sprintf("https://www.hipchat.com/v2/user/%s?auth_token=%s", user, os.Getenv("HIP_CHAT_TOKEN"))
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("get")
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read all")
		log.Fatal(err)
		return "", err
	}

	var data User
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("json")
		log.Fatal(err)
		return "", err
	}

	status <- data.Presence.Show

	return data.Presence.Show, nil

}

func resetLeds(robot *gobot.Robot) {
	robot.Device("dnd").(*gpio.LedDriver).Off()
	robot.Device("away").(*gpio.LedDriver).Off()
	robot.Device("chat").(*gpio.LedDriver).Off()
}

func turnOn(robot *gobot.Robot, device string){
	robot.Device(device).(*gpio.LedDriver).On()
}

func setStatus(robot *gobot.Robot, status string) {

	resetLeds(robot)
	fmt.Println(status)
	switch status {
	case "chat":
		turnOn(robot,"chat")
	case "xa":
		turnOn(robot,"away")
	case "away":
		turnOn(robot,"away")
	case "dnd":
		turnOn(robot,"dnd")
	}
	
}

func gobotFunc(status <-chan string) {

	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	
	dnd := gpio.NewLedDriver(r, "dnd", "16")
	away := gpio.NewLedDriver(r, "away", "18")
	chat := gpio.NewLedDriver(r, "chat", "22")
	
	work := func() {
		resetLeds(gbot.Robot("hipChatBot"))

		for true {

			updatedStatus := <-status

			gobot.After(1*time.Second, func() {
				setStatus(gbot.Robot("hipChatBot"),updatedStatus)
			})
		}

	}

	robot := gobot.NewRobot("hipChatBot",
		[]gobot.Connection{r},
		[]gobot.Device{dnd,away,chat},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

func statusFetcher(statusChannel chan<- string, user string){

	Get(statusChannel, user)	

	c := time.Tick(20 * time.Second)
	for now := range c {
		fmt.Println(now)
		Get(statusChannel, user)
	}

}

func main() { 

	user := os.Args[1]

	statusChannel := make(chan string)

	go statusFetcher(statusChannel, user)

	gobotFunc(statusChannel)
}
