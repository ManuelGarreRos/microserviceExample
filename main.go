package main

import (
	"log"
)

func main() {
	svc := NewMessageService("https://catfact.ninja/fact")
	svc = NewLoggingService(svc)

	ApiServer := NewApiService(svc)
	log.Fatal(ApiServer.start(":8001"))
}
