// This program logs can frames to the console similar to candump from can-utils[1].
//
// [1]: https://github.com/linux-can/can-utils
package main

import (
	"flag"
	"fmt"
	"github.com/ast-dd/can"
	"log"
	"net"
	"os"
	"os/signal"
)

var i = flag.String("if", "", "network interface name")

func main() {
	flag.Parse()
	if len(*i) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	iface, err := net.InterfaceByName(*i)

	if err != nil {
		log.Fatalf("Could not find network interface %s (%v)", *i, err)
	}

	conn, err := can.NewReadWriteCloserForInterface(iface)

	if err != nil {
		log.Fatal(err)
	}

	bus := can.NewBus(conn)
	bus.SubscribeFunc(logCANFrame)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		select {
		case <-c:
			bus.Disconnect()
			os.Exit(1)
		}
	}()

	bus.ConnectAndPublish()
}

// logCANFrame logs a frame with the same format as candump from can-utils.
func logCANFrame(frm can.Frame) {
	data := trimSuffix(frm.Data[:], 0x00)
	length := fmt.Sprintf("[%x]", frm.Length)
	log.Printf("%-3s %-4x %-3s % -24X '%s'\n", *i, frm.ID, length, data, printableString(data[:]))
}

// trim returns a subslice of s by slicing off all trailing b bytes.
func trimSuffix(s []byte, b byte) []byte {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != b {
			return s[:i+1]
		}
	}

	return []byte{}
}

// printableString creates a string from s and replaces non-printable bytes (i.e. 0-32, 127)
// with '.' – similar how candump from can-utils does it.
func printableString(s []byte) string {
	var ascii []byte
	for _, b := range s {
		if b < 32 || b > 126 {
			b = byte('.')

		}
		ascii = append(ascii, b)
	}

	return string(ascii)
}
