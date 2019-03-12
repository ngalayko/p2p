package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ngalayko/p2p/client"
	"github.com/ngalayko/p2p/instance"
	"github.com/ngalayko/p2p/logger"
)

var (
	logLevel          = flag.String("log_level", "info", "logging level [debug|info|warning|error|panic]")
	udp6Multicast     = flag.String("udp6_multicast", "[ff02::114]", "multicast addr for udp6 discrvery")
	udp4Multicast     = flag.String("udp4_multicast", "239.255.255.250", "multicast addr for udp4 discrvery")
	consulAddr        = flag.String("consul", "consul:8500", "consul address")
	port              = flag.Int("port", 30000, "port to listen for messages")
	insecurePort      = flag.Int("insecure_port", 30001, "port to listen for greetings")
	discoveryPort     = flag.String("discovery_port", "30002", "port to discover other peers")
	uiPort            = flag.Int("ui_port", 30003, "port to serve ui interface")
	discoveryInterval = flag.Duration("discovery_interval", 1*time.Second, "interval to send discovery broadcast")
	statisPath        = flag.String("static_path", "./client/public", "path to static files for ui")
	keySize           = flag.Int("key_size", 1024, "private key size")
)

func main() {
	flag.Parse()

	log := logger.New(logger.ParseLevel(*logLevel))
	log.Prefix("main").Info("starting...")

	inst := instance.New(
		log,
		*udp4Multicast,
		*udp6Multicast,
		*consulAddr,
		*discoveryPort,
		*uiPort,
		*port,
		*insecurePort,
		*discoveryInterval,
		*keySize,
	)

	client := client.New(
		log,
		fmt.Sprintf("0.0.0.0:%d", *uiPort),
		inst,
		*statisPath,
	)

	ctx := context.Background()

	go func() {
		if err := client.Start(ctx); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := inst.Start(ctx); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
}
