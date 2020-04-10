package main

import (
	"sync"
)


func CreateWorkGroup(len int) sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(len)

	return wg
}