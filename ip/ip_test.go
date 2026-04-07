package ip

import (
	"log"
	"net"
	"net/netip"
	"testing"

	"github.com/skerkour/go-benchmarks/ip/art"
)

var Sink bool

var routes = []string{"192.168.1.0/24", "192.168.5.5/32", "192.168.0.0/16",
	"192.168.129.0/24", "192.168.128.0/17", "192.168.0.0/16",
	"192.168.5.0/24", "192.168.0.0/16", "192.0.0.0/8",
	"192.168.0.0/16", "192.168.0.0/16", "192.0.0.0/8"}

func BenchmarkNet(b *testing.B) {
	ip := net.ParseIP("10.0.0.0")

	cidrs := make([]*net.IPNet, 0, len(routes))
	for _, route := range routes {
		_, cidr, err := net.ParseCIDR(route)
		if err != nil {
			log.Fatal(err)
		}
		cidrs = append(cidrs, cidr)
	}

	for i := 0; i < b.N; i++ {
		for _, cidr := range cidrs {
			Sink = cidr.Contains(ip)
		}
	}
}

func BenchmarkNetIP(b *testing.B) {
	ip := netip.MustParseAddr("10.0.0.0")

	cidrs := make([]netip.Prefix, 0, len(routes))
	for _, route := range routes {
		cidr := netip.MustParsePrefix(route)
		cidrs = append(cidrs, cidr)
	}

	for i := 0; i < b.N; i++ {
		for _, cidr := range cidrs {
			Sink = cidr.Contains(ip)
		}
	}
}

func BenchmarkArt(b *testing.B) {
	ip := netip.MustParseAddr("10.0.0.0")
	table := &art.Table[struct{}]{}

	for _, route := range routes {
		cidr := netip.MustParsePrefix(route)
		table.Insert(cidr, struct{}{})
	}

	for i := 0; i < b.N; i++ {
		_, Sink = table.Get(ip)
	}
}
