package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"net/http"
	"os"
)

type Monit struct {
	// Have to specify where to find the series  since
	// the field of this struct doesn't match the xml tag
	Services []Service `xml:"service"`
}

type Service struct {
	ServiceType int    `xml:"type,attr"`
	Name        string `xml:"name"`
	Status      int    `xml:"status"`
	Uptime      int32  `xml:"uptime"`
	Memory      Stat   `xml:"memory"`
	CPU         Stat   `xml:"cpu"`
	Monitor     int    `xml:"monitor"`
}

type Stat struct {
	Percent float32 `xml:"percent"`
}

func (s Service) String() string {
	return fmt.Sprintf("%s\tStatus: %d\tUptime: %d\tMemory: %.2f\tCPU: %.2f", s.Name, s.Status, s.Uptime, s.Memory.Percent, s.Memory.Percent)
}

func main() {
	response, err := http.Get("http://admin:monit@192.168.1.26:2812/_status?format=xml")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	//contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	/*
			xmlFile, err := os.Open("monit.xml")
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer xmlFile.Close()
			//contents, _ := ioutil.ReadAll(xmlFile)
		    //var monit Monit
		    //xml.Unmarshal(contents, &monit)

	*/
	var monit Monit
	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&monit)

	for _, service := range monit.Services {
		fmt.Printf("%s\n", service)
	}
}
