# victorops-go
An unofficial Go library for the VictorOps on-call pager rotation service

## Current Status
The [VictorOps API](http://victorops.force.com/knowledgebase/articles/Integration/Alert-Ingestion-API-Documentation/?l=en_US&fs=Search&pn=1) is very limited and thus, this library is also limited.  It only supports the sending of events.

## Usage Example

```go
package main

import (
	"github.com/chrissnell/victorops-go"
	"log"
	"time"
)

func main() {

	vo := victorops.NewClient("YOUR_API_KEY_GOES_HERE")

  // Only RoutingKey and MessageType are required.  The rest are all optional.
  // See the Event struct for all possible options.
  // If you don't set an EntityID, one will be set for you and returned in the response.  You can re-use
  // this EntityID to send updated events (i.e. changed status/MessageType) in subsequent requests.
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
```

## Author
(c) 2015 Christopher Snell  -  http://output.chrissnell.com
