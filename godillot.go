package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Structures for monit xml files
type Server struct {
	Name     string
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
	return fmt.Sprintf("%s\tStatus: %d\tUptime: %s\tMemory: %.1f%%\tCPU: %.1f%%", s.Name, s.Status, transformUptime(s.Uptime), s.Memory.Percent, s.CPU.Percent)
}

//Transform uptime (seconds) to string "dayshoursminutes"
func transformUptime(uptime int32) string {
	var minutes int32 = uptime % 3600 / 60
	var hours int32 = uptime % 86400 / 3600
	var days int32 = uptime / 86400
	return fmt.Sprintf("%dd%dh%dm", days, hours, minutes)
}

//Struct file for configuration file
type Conf struct {
	Servers    []ServerConf `yaml:"servers"`
	OutputFile string       `yaml:"outputfile"`
}

type ServerConf struct {
	Url        string `yaml:"url"`
	ServerName string `yaml:"server"`
}

//Get conf from yaml file
func (c *Conf) getConf() *Conf {

	yamlFile, err := ioutil.ReadFile("godillot.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yamlFile.Unmarshal: %v", err)
	}

	return c
}

// Struct for HTML page

type Data struct {
	Servers []Server
}

// create html page using go templates
func createPage(data Data) {
	t := template.New("server template")

	t = template.Must(t.ParseFiles("layout.html", "server.html"))

	// Exécution de la fusion et injection dans le flux de sortie
	// La variable p sera réprésentée par le "." dans le layout
	// Exemple {{.}} == p
	err := t.ExecuteTemplate(os.Stdout, "layout", data)

	if err != nil {
		log.Fatalf("createPage.Template: %s", err)
	}
}

func main() {

	var configuration Conf
	var serversList []Server
	configuration.getConf()

	fmt.Printf("\n** Godillot v0.5**\n")

	for _, serverConf := range configuration.Servers {
		response, err := http.Get(serverConf.Url)
		if err != nil {
			log.Fatalf("main.GetUrl: %v", err)
			os.Exit(1)
		}
		defer response.Body.Close()
		if err != nil {
			log.Fatalf("main.Body: %v", err)
			os.Exit(1)
		}

		var server Server
		decoder := xml.NewDecoder(response.Body)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&server)
		if err != nil {
			log.Fatalf("main.Unmarshal: %v", err)
			os.Exit(1)
		}
		server.Name = serverConf.ServerName

		serversList = append(serversList, server)

	}

	for _, server := range serversList {
		fmt.Printf("\n** %s **\n", server.Name)
		for _, service := range server.Services {
			fmt.Printf("%s\n", service)
		}
	}
	htmlData := Data{serversList}
	createPage(htmlData)

}
