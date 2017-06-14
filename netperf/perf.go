package netperf

import (
	"net/http"
	_ "net/http/pprof"
	ppf "runtime/pprof"
)

func RunNetPerf(host string) error {
	if host == "" {
		host = ":6060"
	}
	http.HandleFunc("/debug/goroutine", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		p := ppf.Lookup("goroutine")
		p.WriteTo(w, 1)
	})

	return http.ListenAndServe(host, nil)
}
