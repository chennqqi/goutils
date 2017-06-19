//net perf function {$host}:{$port}/debug/pprof
package netperf

import (
	"net/http"
	_ "net/http/pprof"
	ppf "runtime/pprof"
)

func ListenAndServe(host string) error {
	if host == "" {
		host = ":6060"
	}
	return http.ListenAndServe(host, nil)
}
