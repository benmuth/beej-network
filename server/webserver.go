package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseAddr(addr string) [4]byte {
	seg := ""
	parsedAddr := [4]byte{}
	i, j := 0, 0
	for j <= len(addr) {
		if j == len(addr) || addr[j] == '.' {
			b, err := strconv.ParseInt(seg, 10, 9)
			check(err)
			parsedAddr[i] = byte(b)
			i++
			seg = ""
		} else {
			seg += string(addr[j])
		}
		j++
	}
	return parsedAddr
}

func main() {
	port := 0
	var err error
	if len(os.Args) == 2 {
		port, err = strconv.Atoi(os.Args[1])
	} else {
		port = 28333
	}
	fmt.Println("running on port ", port)

	socketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	check(err)

	defer unix.Close(socketFD)

	fmt.Println("socket: ", socketFD)

	check(unix.Bind(socketFD, &unix.SockaddrInet4{Port: port, Addr: [4]byte{127, 0, 0, 1}}))
	check(unix.Listen(socketFD, 0))

	const msg = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 6\r\nConnection: close\r\n\r\nHello!\r\n"
	for {
		nfd, _, err := unix.Accept(socketFD)
		check(err)
		response := make([]byte, 0)
		total := 0

		for {
			chunk := make([]byte, 1024)
			n, _, err := unix.Recvfrom(nfd, chunk, 0)
			check(err)
			response = append(response, chunk...)
			total += n
			if bytes.Contains(chunk, []byte("\r\n\r\n")) {
				break
			}
		}
		check(unix.Send(nfd, []byte(msg), 0))

		fmt.Printf("response:\n%s\n", response[:total])
		check(unix.Close(nfd))
	}

}
