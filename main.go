package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	p2pd "github.com/libp2p/go-libp2p-daemon"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	secio "github.com/libp2p/go-libp2p-secio"
	"github.com/multiformats/go-multiaddr"
	"os"
	"os/signal"
	"syscall"
)

var controlAddrStr = flag.String("control", "", "control multiaddress")
var listenAddrStr = flag.String("listen", "", "listen multiaddress")
var peerIDStr = flag.String("peerID", "", "peer ID")

func check(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
}

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
		_, _ = os.Stderr.WriteString("\nRemaining arguments are parsed as multiaddresses for specified peer\n")
	}
	flag.Parse()


	controlAddr, err := multiaddr.NewMultiaddr(*controlAddrStr)
	check(err, "invalid control addr")
	listenAddr, err := multiaddr.NewMultiaddr(*listenAddrStr)
	check(err, "invalid listen addr")

	security := libp2p.Security(secio.ID, secio.New)

	options := []libp2p.Option{
		security,
	}

	ctx := context.Background()
	d, err := p2pd.NewDaemon(ctx, controlAddr, "", options...)
	check(err, "failed to create new daemon")

	// concurrently connect to our new daemon
	go connect(controlAddr, listenAddr)

	// now wait for stop signal to close daemon
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		d.Close()
		os.Exit(0)
	}
}

func connect(controlAddr, listenAddr multiaddr.Multiaddr) {
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
