// +build linux

package memreflect

import (
	"log"
	"net"

	tproxy "github.com/LiamHaworth/go-tproxy"
)

// MemReflect reflect a kill switch to
// the affected memcached server. Mitigating
// the effect for DRDoS.
type MemReflect struct {
	ln       *net.UDPConn // the listener
	Port     int          // The port to listen on
	Shutdown bool         // Whether or not to shutdown the server
}

// ListenAndServe on a port to reflect command
func ListenAndServe(port int, shutdown bool) error {
	m := MemReflect{
		Port:     port,
		Shutdown: shutdown,
	}
	return m.ListenAndServe()
}

// ListenAndServe on a port to reflect command
func (m *MemReflect) ListenAndServe() error {
	var err error
	log.Printf("Binding UDP TProxy listener to 0.0.0.0: %d\n", m.Port)
	udpListener, err := tproxy.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: m.Port})
	if err != nil {
		return err
	}
	defer udpListener.Close()

	m.ln = udpListener
	return m.listenUDP()
}

func (m *MemReflect) listenUDP() error {
	for {
		buff := make([]byte, 1500)
		_, srcAddr, dstAddr, err := tproxy.ReadFromUDP(m.ln, buff)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("Temporary error while reading data: %s", netErr)
				continue
			}
			log.Printf("Unrecoverable error while reading data: %s", err)
			return err
		}

		log.Printf("UDP attack from %s to %s", srcAddr.String(), dstAddr.String())
		go m.killMemcached(srcAddr, dstAddr)
	}
}

func (m *MemReflect) killMemcached(srcAddr, dstAddr *net.UDPAddr) {
	// send kill command to affected Memcached server
	conn, err := tproxy.DialUDP("udp", dstAddr, srcAddr)
	if err != nil {
		log.Printf("Failed to connect to original UDP source [%s]: %s", srcAddr.String(), err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(m.killCommand())
	if err != nil {
		log.Printf("Encountered error while writing to remote [%s]: %s", conn.RemoteAddr(), err)
		return
	}
}

func (m *MemReflect) killCommand() []byte {
	if m.Shutdown {
		return []byte("\x00\x00\x00\x00\x00\x01\x00\x00flush_all\r\nshutdown\r\n")
	}
	return []byte("\x00\x00\x00\x00\x00\x01\x00\x00flush_all\r\n")
}
