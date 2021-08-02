package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
)

type Result struct {
	Discover int    `json:"discover"`
	Offer    int    `json:"offer"`
	Request  int    `json:"request"`
	Ack      int    `json:"ack"`
	Message  string `json:"message"`
}

var (
	iface = flag.String("i", "eth0", "Interface to configure via DHCP")
)

func main() {
	flag.Parse()
	log.Printf("Starting DHCP client on interface %s", *iface)

	client := client4.NewClient()

	var result Result
	conversation, err := client.Exchange(*iface)

	// Summary() prints a verbose representation of the exchanged packets.
	for _, packet := range conversation {
		switch mt := packet.MessageType(); mt {
		case dhcpv4.MessageTypeDiscover:
			result.Discover = 1
		case dhcpv4.MessageTypeOffer:
			result.Offer = 1
		case dhcpv4.MessageTypeRequest:
			result.Request = 1
		case dhcpv4.MessageTypeAck:
			result.Ack = 1
		}
	}

	if err != nil {
		result.Message = fmt.Sprintf("%s", err)
	} else {
		result.Message = "None"
	}
	postResult(result)
}

func postResult(result Result) {
	b, err := json.Marshal(result)
	log.Print(string(b))

	if err != nil {
		log.Fatal(err)
	}

	url := "http://harvest.soracom.io"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)
}
