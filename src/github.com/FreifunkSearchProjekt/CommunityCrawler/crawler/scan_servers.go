package crawler

import (
	"github.com/anvie/port-scanner"
	"time"
)

func ScanServer(ip string) (openPorts []int) {
	// scan ip with a 2 second timeout per port in 5 concurrent threads
	ps := portscanner.NewPortScanner(ip, 2*time.Second, 5)

	openPorts = ps.GetOpenedPort(1, 30000)
	return
}
