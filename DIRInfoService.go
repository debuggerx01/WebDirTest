package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

var maxRoutineNum = 5

type DIRInfo struct {
	Path      string `json:"path"`
	DirCount  int    `json:"dirCount"`
	FileCount int    `json:"fileCount"`
	TotalSize int64  `json:"totalSize"`
	parent    *DIRInfo
}

var dirs []*DIRInfo
var lock sync.Mutex

var ch = make(chan int, maxRoutineNum)

func handleDir(p string, parent *DIRInfo) {
	d := &DIRInfo{
		Path:      p,
		DirCount:  0,
		FileCount: 0,
		TotalSize: 0,
		parent:    parent,
	}

	lock.Lock()
	dirs = append(dirs, d)
	lock.Unlock()

	infos, _ := ioutil.ReadDir(p)

	hasDir := false

	for _, i := range infos {
		if i.IsDir() {
			d.DirCount++
			fullPath := path.Join(p, i.Name())
			hasDir = true
			ch <- 1
			go handleDir(fullPath, d)
		} else {
			d.FileCount++
			d.TotalSize += i.Size()
		}
	}

	if !hasDir {
		var parent = d.parent
		for parent != nil {
			parent.TotalSize += d.TotalSize
			parent.FileCount += d.FileCount
			parent.DirCount += d.DirCount + 1
			parent = parent.parent
		}
	}
	<-ch
}

func runDir(p string) string {
	dirs = make([]*DIRInfo, 0)
	go handleDir(p, nil)
	time.Sleep(time.Second)
	res, _ := json.Marshal(dirs)
	r, _ := json.Marshal(dirs[0])
	fmt.Println(string(r))
	return string(res)

}
