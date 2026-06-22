package pprofserver

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

type PprofServerAdapter struct {
	addr string
}

func NewPprofServerAdapter(addr string) *PprofServerAdapter {
	return &PprofServerAdapter{
		addr: addr,
	}
}

func (ps *PprofServerAdapter) Start() error {
	log.Printf("pprof server starting at http://localhost%v...", ps.addr)
	
	return http.ListenAndServe(ps.addr, nil)
}
