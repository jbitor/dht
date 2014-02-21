// Package dht provides an implementation of a "mainline" BitTorrent
// Distributed Hash Table (DHT) node, as specified in BEP 5 (http://www.bittorrent.org/beps/bep_0005.html),
// and a higher-level client interface for querying the DHT.
package dht

import "net"

var defaultNodes = []RemoteNode{
	// Possible node: hope for another node running locally, on the default port
	{
		Address: net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 6881,
		},
	},

	// Bootstrap node: router.bittorrent.com
	{
		Address: net.UDPAddr{
			IP:   net.IPv4(67, 215, 242, 139),
			Port: 6881,
		},
		BootstrapOnly: true,
	},

	// Bootstrap node: dht.transmissionbt.com
	{
		Address: net.UDPAddr{
			IP:   net.IPv4(91, 121, 60, 42),
			Port: 6881,
		},
		BootstrapOnly: true,
	},
}
