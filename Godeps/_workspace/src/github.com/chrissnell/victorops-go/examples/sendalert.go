package main

import (
	"github.com/chrissnell/victorops-go"
	"log"
	"time"
)

func main() {

	vo := victorops.NewClient("YOUR_API_KEY_GOES_HERE")

	e := &victorops.Event{
		RoutingKey:        "SomeRoutingKey",
		MessageType:       victorops.Critical,
		EntityID:          "SomeServer/Disk",
		StateMessage:      "Disk space is almost full",
		Timestamp:         time.Now(),
		EntityDisplayName: "SomeServer",
	}

	resp, err := vo.SendAlert(e)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Printf("Response: %+v\n", resp)
}
