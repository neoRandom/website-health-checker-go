package main

import (
	"http-server/internal/adapter/driver"
)

func main() {
	s := driver.ServerAdapter{
		Addr: "localhost:8080",
	}

	s.Init()
}
