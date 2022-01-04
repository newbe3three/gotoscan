package gotoscan

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
)

//对多个host的并发操作
func HostWorker(hosts []string, cmslist map[string][]CmsFeature, sortList CmsSortList) []string {
	hostsChan := make(chan string)
	resultChan := make(chan string)
	var resultList []string
	for i := 0; i < len(hosts); i++ {
		go cmsWorker(hostsChan, cmslist, sortList, resultChan)
	}
	for _, host := range hosts {
		hostsChan <- host
	}
	for i := 0; i < len(hosts); i++ {
		result := <-resultChan
		resultList = append(resultList, result)
	}

	close(hostsChan)
	close(resultChan)
	return resultList
}

//对多个cms的并发操作
func cmsWorker(hosts chan string, cmslist map[string][]CmsFeature, sortList CmsSortList, resultChan chan string) {
	for host := range hosts {
		var scanStatus bool = false
		cmsListChan := make(chan map[string][]CmsFeature, 10)

		var wg sync.WaitGroup
		for i := 0; i < cap(cmsListChan); i++ {
			go featureWorker(host, cmsListChan, &wg, &scanStatus, resultChan)
		}
		for _, data := range sortList {
			if !scanStatus {
				wg.Add(1)
				cmsListChan <- map[string][]CmsFeature{data.Name: cmslist[data.Name]}
			} else {
				//wg1.Done()
				return
			}

		}

		wg.Wait()
		resultChan <- fmt.Sprintf("The host: %s has no matching results", host)
		close(cmsListChan)

	}

}

func featureWorker(host string, cmslistchan chan map[string][]CmsFeature, wg *sync.WaitGroup, scanStatus *bool, resultChan chan string) {

	for cmslist := range cmslistchan {
		if *scanStatus {
			wg.Done()
			return
		}
		for k, v := range cmslist {
			for _, feature := range v {
				url := host + feature.Path
				fmt.Print(".")
				resp, err := HeadReq(url)
				if err != nil {
					continue
				}
				if resp.StatusCode != 200 {
					continue
				}
				content, err := GetReq(url)
				if err != nil {
					continue
				}
				if feature.Option == "keyword" {
					if strings.Contains(strings.ToLower(string(content)), strings.ToLower(feature.Content)) {
						resultChan <- fmt.Sprintf("The host: %s matches CMS: %s,it's fingerprint is:path: %s,option: keyword,content: %s", host, k, feature.Path, feature.Content)
						*scanStatus = true
						wg.Done()
						return
					}

				} else if feature.Option == "md5" {
					md5str := fmt.Sprintf("%x", md5.Sum(content))
					if string(md5str) == feature.Content {
						resultChan <- fmt.Sprintf("The host: %s matches CMS: %s,it's fingerprint is:path: %s,option: md5,content: %s", host, k, feature.Path, feature.Content)
						*scanStatus = true
						wg.Done()
						return
					}
				}

			}
		}
		wg.Done()
	}
}
