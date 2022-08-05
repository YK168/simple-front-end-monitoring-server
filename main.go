package main

import (
	"log"
	"simple_front_end_monitoring_server/conf"
	"simple_front_end_monitoring_server/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	if err := r.Run(conf.HttpPort); err != nil {
		log.Fatalln(err)
	}
}
