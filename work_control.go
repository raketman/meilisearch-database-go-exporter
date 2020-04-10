package main

import "log"

func WorkControl(works chan int, done chan bool, name string, counter int) {
	go func() {
		log.Println(name,"COUNTER", counter)
		for {
			_, more := <-works

			counter = counter - 1
			log.Println(name,"CHANGE COUNTER", counter)
			if counter == 0 || !more {
				log.Println(name,"FINISH", counter)
				done <- true
				return
			}
		}
	}()
}