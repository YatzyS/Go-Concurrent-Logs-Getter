package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

func create() {
	iso8601 := "2006-01-02T15:04:05.00Z"
	generalFormat := "2006-01-02"
	date := "2019-09-18"
	startDate, _ := time.Parse(generalFormat, date)
	fileNumberStart := 301
	fileNumberEnd := 350
	var wg sync.WaitGroup
	wg.Add(fileNumberEnd-fileNumberStart+1)
	for i := fileNumberStart; i <= fileNumberEnd; i++ {
		go func(startDate time.Time, i int){
			data := ""
			fmt.Println(fmt.Sprintf("startIndex %d",i))
			for j := 1; j <= 60; j++ {
				for k := 1; k <= 1000; k++ {
					logs := ""
					if j == 60 && k == 1000 {
						logs += ",logData G,logData H, logData I"
					} else {
						logs += ",logData G,logData H, logData I\n"
					}
					data += startDate.Format(iso8601) + logs
				}
				startDate = startDate.AddDate(0, 0, 1)
			}
			ioutil.WriteFile(fmt.Sprintf("LogFile-%d.log", i), []byte(data), 0644)
			fmt.Println(fmt.Sprintf("endIndex %d",i))
			wg.Done()
		}(startDate, i)
		incrementDays := 60 * (i % ((fileNumberEnd - fileNumberStart)+1))
		startDate = startDate.AddDate(0, 0, incrementDays)
	}
	wg.Wait()
}
