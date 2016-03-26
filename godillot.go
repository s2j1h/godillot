package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Monit struct {
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
	return fmt.Sprintf("%s\tStatus: %d\tUptime: %s\tMemory: %.2f\tCPU: %.2f", s.Name, s.Status, transformUptime(s.Uptime), s.Memory.Percent, s.Memory.Percent)
}

//Transform uptime (seconds) to string "dayshoursminutes"
func transformUptime(uptime int32) string {
	var minutes int32 = uptime % 3600 / 60
	var hours int32 = uptime % 86400 / 3600
	var days int32 = uptime / 86400
	return fmt.Sprintf("%dd%dh%dm", days, hours, minutes)
}

type Conf struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Url string `yaml:"url"`
}

//Get conf from yaml file
func (c *Conf) getConf() *Conf {

	yamlFile, err := ioutil.ReadFile("godillot.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {

	var configuration Conf
	configuration.getConf()

	for _, server := range configuration.Servers {
		response, err := http.Get(server.Url)
		if err != nil {
			log.Fatalf("main GetUrl: %v", err)
			os.Exit(1)
		}
		defer response.Body.Close()
		if err != nil {
			log.Fatalf("main GetUrl: %v", err)
			os.Exit(1)
		}

		var monit Monit
		decoder := xml.NewDecoder(response.Body)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&monit)

		for _, service := range monit.Services {
			fmt.Printf("%s\n", service)
		}

	}

}
