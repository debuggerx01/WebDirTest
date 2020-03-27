package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var root = flag.String("P", "/", "root path")
var port = flag.Int("p", 9999, "server port")
var serPath = fmt.Sprintf(":%d", *port)

type FileInf struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

func dir(p string) string {
	infos, _ := ioutil.ReadDir(p)

	fileInfs := make([]*FileInf, 0)

	for _, i := range infos {

		f := &FileInf{
			Name:  i.Name(),
			IsDir: i.IsDir(),
			Size:  i.Size(),
		}

		fileInfs = append(fileInfs, f)
	}
	res, _ := json.Marshal(fileInfs)

	return string(res)
}

func dirHandler(w http.ResponseWriter, req *http.Request) {
	queryPath := req.URL.Query().Get("path")
	var p = path.Join(*root, queryPath)
	stat, err := os.Stat(p)
	if err == nil && stat.IsDir() {
		_, _ = fmt.Fprintln(w, runDir(p))
	} else {
		_, _ = fmt.Fprintln(w, "invalid path")
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", dirHandler)
	_ = http.ListenAndServe(serPath, nil)
}
