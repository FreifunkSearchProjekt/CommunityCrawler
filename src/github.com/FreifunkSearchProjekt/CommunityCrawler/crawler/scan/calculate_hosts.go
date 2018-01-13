package scan

import "net"

// Hosts calculates the Ips in a subnet TODO: Use map instead slice
func Hosts(cidr string) (map[int]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ips := make(map[int]net.IP)
	var ipslice []net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ipslice = append(ipslice, ip)
	}
	// remove network address and broadcast address
	ipslice = ipslice[1 : len(ips)-1]
	for i, e := range ipslice  {
		ips[i] = e
	}
	return ips, nil
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

