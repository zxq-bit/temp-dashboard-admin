package main

import (
	"log"

	"github.com/caicloud/dashboard-admin/pkg/admin/server"
)

func main() {
	s, e := server.NewServer()
	if e != nil {
		log.Fatalf("NewServer failed, %v", e)
	}
	s.Run()
}
