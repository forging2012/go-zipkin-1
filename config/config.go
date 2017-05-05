package config

import "flag"

const (
	//Server ...
	Server = "server"
	//Client ...
	Client = "client"
)

var (
	//ZipkinURL ...
	ZipkinURL = flag.String("url", "http://localhost:9411/api/v1/spans", "Zipkin server URL")
	//ServerPort ...
	ServerPort = flag.String("port", "8000", "server port")
	//ActorKind ...
	ActorKind = flag.String("actor", "server", "server or client")
)

//Get the config
func Get() {
	flag.Parse()
}
