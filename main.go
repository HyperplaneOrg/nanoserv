// Use of this source code is governed by the BSD 3-Clause 
// License that can be found in the LICENSE file.

//  This is a trivial http server that serves json from the EndPoints defined in the
//  yaml config file. I was interested in looking at the net/http and yaml packages,
//  and mocking up some apis. The size/lines of the code for this nano server is
//  small enough to fit in your back pocket :~)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type NanoServer struct {
	Port       string
	Root       string
	MaxUri     int
	router     *mux.Router
	EndPoints  map[string]string
	ServerInfo string
}

var nanSrv = NanoServer{}

func NanoSeverInfo(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(rw, nanSrv.ServerInfo)
}

func NanoSeverIntError(rw http.ResponseWriter) {
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(rw, "500 internal server error")
}

func NanoSeverJson(rw http.ResponseWriter, r *http.Request) {
	requri := r.URL.RequestURI()
	if len(requri) > nanSrv.MaxUri {
		NanoSeverIntError(rw)
		return
	}
	path := nanSrv.EndPoints[requri]
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		Openfile, err := os.Open(path)
		defer Openfile.Close()
		if err != nil {
			NanoSeverIntError(rw)
			return
		} else {
			FileStat, _ := Openfile.Stat()
			FileSize := strconv.FormatInt(FileStat.Size(), 10)
			rw.Header().Set("Content-Type", "application/json")
			rw.Header().Set("Content-Length", FileSize)
			io.Copy(rw, Openfile)
			return
		}

	} else {
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(rw, "404 page not found")
		return
	}
}

func (n *NanoServer) InitNanoServer(sconf *NanoServerConfigInfo) {
	n.Port = strconv.Itoa(sconf.Config.Port)
	n.Root = sconf.Config.Root
	n.MaxUri = sconf.Config.MaxUriRequest
	n.router = mux.NewRouter()
	n.EndPoints = make(map[string]string)
	endnames := []string{}
	for _, endpoint := range sconf.Config.EndPoints {
		os.MkdirAll(endpoint.Path, 0755)
		n.EndPoints[endpoint.Uri] = endpoint.Path + "/" + endpoint.Data
		n.router.HandleFunc(endpoint.Uri, NanoSeverJson)
		endnames = append(endnames, endpoint.Uri)
	}
	n.router.HandleFunc("/", NanoSeverInfo)
	sinfo := make(map[string][]string)
	sinfo["EndPoints"] = endnames
	sinfo["ServerName"] = []string{sconf.Config.Name}
	sinfo["Version"] = []string{sconf.Config.Version}
	tret, _ := json.Marshal(sinfo)
	n.ServerInfo = string(tret)
}

func NanoServUsage() {
    helpstring := "\nUsage: nanoserv <config.yml>\nSee the manual for the basic config.yml schema...\n\n"
    fmt.Printf(helpstring)
}

func main() {
	flag.Usage = NanoServUsage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	yamlFname := flag.Args()[0]
	srvrconf := NanoServerLoadConfig(yamlFname)
    fmt.Println(srvrconf)
	nanSrv.InitNanoServer(&srvrconf)
	loggedRouter := handlers.LoggingHandler(os.Stdout, nanSrv.router)
	err := http.ListenAndServe(":"+nanSrv.Port, loggedRouter)
	if err != nil {
		log.Fatal("nanoserv http.ListenAndServe Error : ", err)
	}
}
