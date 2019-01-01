## nanoserv
This is a trivial http server that serves json from the EndPoints defined in the
yaml config file. I was interested in looking at the net/http and yaml packages,
and mocking up some apis. The size/lines of the code for this nano server is
surpisingly small. 

### The configuration schema

An example yaml configuration: 

```YAML
---
server:
  name: nanoserver
  version: 0.0.2
  maxUriRequest: 2048
  port: 8080
  root: /Users/foobar/golang/src/nanoserv/

  endPoints:
    - name: hello1
      uri: /hello1/
      relpath: /root/hello1/
      data: index.json

    - name: hello2
      uri: /hello2
      relpath: root//hello2/
      data: foo.json
```

Note that relpath is a path relative to the root path, that is with a 

__relpath: /root/hello1/__

AND

__root: /Users/foobar/golang/src/nanoserv/__

the filesystem path would be 

__/Users/foobar/golang/src/nanoserv/root/hello1/index.json__

Note, the paths are cleaned up in the config code. For this example, the endpoint

__http://localhost:8080/hello1/__

would return

__/Users/foobar/golang/src/nanoserv/root/hello1/index.json__

AND

__http://localhost:8080/hello2/__

would return

__/Users/foobar/golang/src/nanoserv/root/hello2/foo.json__


For basic server info and endpoint discovery, the

__http://localhost:8080/__

would return

```json
{ "EndPoints":["/hello1/","/hello2/"],
  "ServerName":["nanoserver"],
  "Version":["0.0.2"] }
```

### Build and Run

To build and run this simple utility: 

```shell
$ go get github.com/gorilla/handlers github.com/gorilla/mux gopkg.in/yaml.v2
$ go build
$ go install
```

Assuming you have the go bin path in your PATH:

```shell
$ nanoserv config.yml
```


