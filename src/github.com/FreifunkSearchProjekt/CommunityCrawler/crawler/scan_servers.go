package crawler

import (
	"github.com/anvie/port-scanner"
	"net"
	"time"
)

// ScanServer scans the Server for open ports to crawl from port 1 to 30000
func ScanServer(ip string) (openPorts []int) {
	// scan ip with a 2 second timeout per port in 5 concurrent threads
	ps := portscanner.NewPortScanner(ip, 2*time.Second, 5)

	openPorts = ps.GetOpenedPort(1, 30000)
	return
}

// Hosts calculates the Ips in a subnet
func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
