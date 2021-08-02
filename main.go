package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
)

type Result struct {
	Discover bool `json:"discover"`
	Offer    bool `json:"offer"`
	Request  bool `json:"request"`
	Ack      bool `json:"ack"`
}

var (
	iface = flag.String("i", "eth0", "Interface to configure via DHCP")
)

func main() {
	flag.Parse()
	log.Printf("Starting DHCP client on interface %s", *iface)

	client := client4.NewClient()

	conversation, err := client.Exchange(*iface)

	var result Result
	// Summary() prints a verbose representation of the exchanged packets.
	for _, packet := range conversation {
		switch mt := packet.MessageType(); mt {
		case dhcpv4.MessageTypeDiscover:
			result.Discover = true
		case dhcpv4.MessageTypeOffer:
			result.Offer = true
		case dhcpv4.MessageTypeRequest:
			result.Request = true
		case dhcpv4.MessageTypeAck:
			result.Ack = true
		}
	}

	if err != nil {
		log.Fatal(err)
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
