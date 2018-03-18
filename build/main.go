// +build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/Max-Sum/memreflect"
	"github.com/coreos/go-iptables/iptables"
)

func main() {
	log.Println("Starting memreflect server")
	// Listen to interrupt signal
	interruptListener := make(chan os.Signal)
	signal.Notify(interruptListener, os.Interrupt)

	s := flag.Bool("s", false, "Whether or not to shutdown the memcached server")
	p := flag.Int("p", 11211, "The port to listen on")
	// Setup iptables and ip route
	defer shutdown()
	if err := setup(*p); err != nil {
		log.Println(err)
		return
	}

	go func() {
		err := memreflect.ListenAndServe(*p, *s)
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-interruptListener
	log.Println("memreflect server closing")
}

// setup iptables rules to listen
func setup(port int) error {
	CHAIN := "MEMREFLECT"
	whitelist := []string{"0.0.0.0/8", "127.0.0.0/8", "224.0.0.0/8", "240.0.0.0/8"}
	PORT := fmt.Sprintf("%d", port)
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	// Check if any left chains
	chains, err := ipt.ListChains("mangle")
	if err != nil {
		return err
	}
	for _, c := range chains {
		if c == CHAIN {
			if err = shutdown(); err != nil {
				return err
			}
		}
	}
	if err = ipt.NewChain("mangle", CHAIN); err != nil {
		return err
	}
	// whitelist ips will returns
	for _, rule := range whitelist {
		if err = ipt.AppendUnique("mangle", CHAIN, "-d", rule, "-j", "RETURN"); err != nil {
			return err
		}
	}
	if err = ipt.AppendUnique("mangle", CHAIN, "-p", "udp", "--sport", "11211", "-j",
		"TPROXY", "--on-port", PORT, "--tproxy-mark", "0x1f/0x1f"); err != nil {
		return err
	}
	// append chain into prerouting and output
	if err = ipt.AppendUnique("mangle", "PREROUTING", "-p", "udp", "-j", CHAIN); err != nil {
		return err
	}
	// Set up route
	if err = exec.Command("ip", "rule", "add", "fwmark", "0x1f/0x1f", "table", "11211").Run(); err != nil {
		return err
	}
	if err = exec.Command("ip", "route", "add", "local", "0.0.0.0/0", "dev", "lo", "table", "11211").Run(); err != nil {
		return err
	}
	return nil
}

func shutdown() error {
	CHAIN := "MEMREFLECT"
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	if err = ipt.ClearChain("mangle", CHAIN); err != nil {
		log.Println(err)
	}
	if err = ipt.Delete("mangle", "PREROUTING", "-p", "udp", "-j", CHAIN); err != nil {
		log.Println(err)
	}
	if err = ipt.DeleteChain("mangle", CHAIN); err != nil {
		log.Println(err)
	}
	// clear route
	if err = exec.Command("ip", "rule", "del", "table", "11211").Run(); err != nil {
		log.Println(err)
	}
	if err = exec.Command("ip", "route", "flush", "table", "11211").Run(); err != nil {
		log.Println(err)
	}
	return err
}
