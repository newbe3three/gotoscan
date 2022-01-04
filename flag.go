package gotoscan

import (
	"flag"
	"log"
)

type ArgsInfo struct {
	Host    string
	Hosts   string
	CmsJson string
}

func Banner() {
	banner := `
	________     ___________     __________________     _____    _______   
	/  _____/  ___\__    ___/___ /   _____/\_   ___ \   /  _  \   \      \  
   /   \  ___ /  _ \|    | /  _ \\_____  \ /    \  \/  /  /_\  \  /   |   \ 
   \    \_\  (  <_> )    |(  <_> )        \\     \____/    |    \/    |    \
	\______  /\____/|____| \____/_______  / \______  /\____|__  /\____|__  /
		   \/                           \/         \/         \/         \/ 
			   			GoToScan version:1.0
		   `
	print(banner)
}

func (Info *ArgsInfo) Flag() {
	Banner()
	//可以指定的参数
	flag.StringVar(&Info.Host, "host", "", "Test a host,http://xxxxx")
	flag.StringVar(&Info.Hosts, "hosts", "", "Filename with hosts,One host per line")
	flag.StringVar(&Info.CmsJson, "cmsjson", "cms.json", "Cms fingerprint feature json file, The default is cms.json")
	flag.Parse()
	if Info.Host == "" && Info.Hosts == "" {
		log.Fatalln("err:./no host parameter")
		return
	}

	if Info.Host != "" && Info.Hosts != "" {
		log.Fatalln("err:./only one host parameter")
		return
	}
}
