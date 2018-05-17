package main

import "os"

func main() {
	var app App

	app.Initialize()
	hostname := os.Getenv("HOST_NAME")
	address := hostname + ":9000"
	app.Run(address)
}
