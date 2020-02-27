package sourcelog

import (
	"fmt"
	"net"
	"os"
)

func Run_server() {
	packet, err := net.ListenPacket("udp", ":33344")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for {
		buf := make([]byte, 1024)
		_, _, err := packet.ReadFrom(buf) /* Should probably track addrs in the future */
		if err != nil {
			continue
		}
		fmt.Println(string(buf[7:]))
	}
}
