package main

import (
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	"github.com/multiformats/go-multiaddr"
)

var controlAddrStr = flag.String("control", "", "control multiaddress")
var listenAddrStr = flag.String("listen", "", "listen multiaddress")
var peerIDStr = flag.String("peerID", "", "peer ID")

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println("\nRemaining arguments are parsed as multiaddresses for specified peer")
	}
	flag.Parse()

	check := func (err error, msg string) {
		if err != nil {
			panic(fmt.Errorf("%s: %v", msg, err))
		}
	}
	controlAddr, err := multiaddr.NewMultiaddr(*controlAddrStr)
	check(err, "invalid control addr")
	listenAddr, err := multiaddr.NewMultiaddr(*listenAddrStr)
	check(err, "invalid listen addr")

	var peerAddrs []multiaddr.Multiaddr
	for i, addrStr := range flag.Args() {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		check(err, fmt.Sprintf("peer multi addr %d is invalid", i))
		peerAddrs = append(peerAddrs, addr)
	}

	cl, err := p2pclient.NewClient(controlAddr, listenAddr)
	check(err, "cannot create client")

	peerID, err := peer.IDFromString(*peerIDStr)
	check(err, "invalid peer ID")

	check(cl.Connect(peerID, peerAddrs), "cannot connect to peer")
}
