package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//Structures for monit xml files
type Server struct {
	Name     string
	Url      string    `xml:"url"`
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
	return fmt.Sprintf("%s\tStatus: %s\tUptime: %s\tMemory: %.1f%%\tCPU: %.1f%%", s.Name, GetStatus(s), TransformUptime(s), s.Memory.Percent, s.CPU.Percent)
}

//Transform uptime (seconds) to string "dayshoursminutes" -- part of Service structure
func TransformUptime(s Service) string {
	uptime := s.Uptime
	var minutes int32 = uptime % 3600 / 60
	var hours int32 = uptime % 86400 / 3600
	var days int32 = uptime / 86400
	return fmt.Sprintf("%dd%dh%dm", days, hours, minutes)
}

//Transform status
func GetStatus(s Service) string {
	var status string = "Failure"
	if s.Status == 0 {
		status = "Running"
	}
	return status
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
		log.Fatalf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yamlFile.Unmarshal: %v", err)
		os.Exit(1)
	}

	return c
}

// Struct for HTML page

type Data struct {
	Servers []Server
	Date    string
}

// create html page using go templates
func createPage(data Data, filename string) {

	t := template.New("servertemplate").Funcs(template.FuncMap{
		"transformUptime": TransformUptime,
		"getStatus":       GetStatus,
	})

	t = template.Must(t.ParseFiles("layout.html"))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("createPage.newFile: %s", err)
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	// Exécution de la fusion et injection dans le flux de sortie
	// La variable p sera réprésentée par le "." dans le layout
	// Exemple {{.}} == p
	//var buff bytes.Buffer

	err2 := t.ExecuteTemplate(w, "layout", data)

	if err2 != nil {
		log.Fatalf("createPage.Template: %s", err2)
	}
	w.Flush()

}

func main() {

	var configuration Conf
	var serversList []Server
	configuration.getConf()

	fmt.Printf("\n** Godillot v0.6**\n")
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05")
	log.Printf("Fetching data...")

	for _, serverConf := range configuration.Servers {
		response, err := http.Get(serverConf.Url + "/_status?format=xml")
		if err != nil {
			log.Fatalf("main.GetUrl: %v", err)
			continue
		}
		defer response.Body.Close()
		if err != nil {
			log.Fatalf("main.Body: %v", err)
			continue
		}

		var server Server
		decoder := xml.NewDecoder(response.Body)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&server)
		if err != nil {
			log.Fatalf("main.Unmarshal: %v", err)
			continue
		}
		server.Name = serverConf.ServerName
		server.Url = serverConf.Url
		serversList = append(serversList, server)

	}

	for _, server := range serversList {
		fmt.Printf("** %s **\n", server.Name)
		for _, service := range server.Services {
			fmt.Printf("%s\n", service)
		}
	}
	log.Printf("Creating html page...")
	htmlData := Data{serversList, newFormat}
	createPage(htmlData, configuration.OutputFile)
	log.Printf("Finish!")

}
