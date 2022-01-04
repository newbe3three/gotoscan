package main

import (
	"bufio"
	"fmt"
	"gotoscan"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		Info    gotoscan.ArgsInfo
		file    io.ReadCloser
		err     error
		hosts   []string
		results []string
		output  string
	)
	t := time.Now()
	Info.Flag()
	if Info.Host != "" {
		fmt.Println("Get the host to be tested is: " + Info.Host)
		output = t.Format("2006-01-02-") + Info.Host[strings.Index(Info.Host, "/")+2:] + ".txt"
		file = ioutil.NopCloser(strings.NewReader(Info.Host))
	} else {
		fmt.Println("Get the hosts file to be tested is: " + Info.Hosts)
		output = t.Format("2006-01-02-") + Info.Hosts + ".txt"
		file, err = os.Open(Info.Hosts)
		if err != nil {
			log.Fatalf("err ./Cannot open the hosts file %s: %s\n", Info.Hosts, err)
		}
	}
	defer file.Close()

	//hosts := make(chan string, 10)
	scan := bufio.NewScanner(file)

	for scan.Scan() {
		log.Printf("Load host: %s\n", scan.Text())
		hosts = append(hosts, scan.Text())
	}
	//fmt.Println(hosts, cap(hosts))
	if Info.CmsJson != "" {
		cmsfile, err := ioutil.ReadFile(Info.CmsJson)
		if err != nil {
			log.Fatalf("err ./Error opening json file: %s", err)
		}
		cmsList, cmsSortList, err := gotoscan.ParseCmsFeatureFromJson(cmsfile)
		if err != nil {
			log.Fatalf("err ./Error parsing json file: %s", err)
		}
		//fmt.Println(cmsList, cmsSortList)
		log.Println("Successfully parsed the json file: " + Info.CmsJson)
		log.Println("Start testing")
		results = gotoscan.HostWorker(hosts, cmsList, cmsSortList)
	}

	outfile, err := os.Create(output)
	if err != nil {
		log.Fatalf("err ./Error creating output file: %s", err)
		return
	}
	defer outfile.Close()
	for _, result := range results {
		fmt.Println(result)
		outfile.WriteString(result + "\n")
	}

	tim := time.Since(t)
	fmt.Printf("Takes: %s", tim)
}
