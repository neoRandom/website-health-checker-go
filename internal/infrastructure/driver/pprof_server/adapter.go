package pprofserver

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

type PprofServerAdapter struct {}

func NewPprofServerAdapter() *PprofServerAdapter {
	return &PprofServerAdapter{}
}

func (ps *PprofServerAdapter) Start() error {
	log.Println("Initializing pprof at http://localhost:6060...")
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		log.Printf("Error initializing pprof: %v", err)
		return err
	}
	return nil
}
