// Use of this source code is governed by the BSD 3-Clause 
// License that can be found in the LICENSE file.

// This supports the trivial http server by handling 
// the simple yaml config files.
package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Some default values, of which the config could override.
const (
	NPORT    = 8080
	NSERVER  = "nanoserver"
	NVERSION = "0.0.1"
	NROOT    = "."
	MAXREQ   = 2048
	NDATA    = "index.json"
)

// A basic structure that maps to the yaml file scheme
type NanoServerConfigInfo struct {
	Config struct {
		Name          string `yaml:"name"`
		Port          int    `yaml:"port"`
		MaxUriRequest int    `yaml:"maxUriRequest"`
		Root          string `yaml:"root"`
	    Version       string `yaml:"version"`
		EndPoints     []struct {
			Name string `yaml:"name"`
			Uri  string `yaml:"uri"`
			Path string `yaml:"relpath"`
			Data string `yaml:"data"`
		} `yaml:"endPoints"`
	} `yaml:"server"`
}

func NanoServerLoadConfig(configpath string) NanoServerConfigInfo {
	var server NanoServerConfigInfo

	yamlFile, err := ioutil.ReadFile(configpath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &server)
	if err != nil {
		panic(err)
	}

	if server.Config.Port == 0 {
		server.Config.Port = NPORT
	}

	if server.Config.MaxUriRequest == 0 {
		server.Config.MaxUriRequest = MAXREQ
	}

	if server.Config.Name == "" {
		server.Config.Name = NSERVER
	}

    if server.Config.Version == "" {
		server.Config.Version = NVERSION
	}

	if server.Config.Root == "" {
		server.Config.Root = NROOT
	}
	server.Config.Root = filepath.Clean(server.Config.Root)

	/* cleanup file paths and endpoints, etc... */
	for i, endpoint := range server.Config.EndPoints {
		server.Config.EndPoints[i].Path = server.Config.Root + "/" + filepath.Clean(endpoint.Path)
		stmp := strings.TrimSpace(endpoint.Uri)
		if !strings.HasSuffix(stmp, "/") {
			server.Config.EndPoints[i].Uri = stmp + "/"
		} else {
			server.Config.EndPoints[i].Uri = stmp
		}
		if endpoint.Data == "" {
			server.Config.EndPoints[i].Data = NDATA
		} else {
			server.Config.EndPoints[i].Data = strings.TrimSpace(endpoint.Data)
		}
	}

	return server
}
