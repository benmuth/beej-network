package main

import (
	"fmt"
	"net"
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
	if len(os.Args) < 2 {
		fmt.Println("provide a url like example.com and optionally a port (default 80)")
		return
	}
	host := os.Args[1]
	port := 0
	var err error
	if len(os.Args) == 3 {
		port, err = strconv.Atoi(os.Args[2])
	} else {
		port = 80
	}

	fmt.Println("running...")
	addrs, err := net.LookupHost(host)
	check(err)

	socketFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	check(err)
	fmt.Println("socket: ", socketFD)

	err = unix.Connect(socketFD, &unix.SockaddrInet4{Port: port, Addr: parseAddr(addrs[0])})
	check(err)

	msg := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", host)
	check(unix.Send(socketFD, []byte(msg), 0))

	response := make([]byte, 0)
	total := 0

	for {
		chunk := make([]byte, 1024)
		n, _, err := unix.Recvfrom(socketFD, chunk, 0)
		check(err)
		response = append(response, chunk...)
		total += n
		if n == 0 {
			break
		}
	}

	fmt.Printf("response:\n%s\n", response[:total])
	defer check(unix.Close(socketFD))
}

// ip1, _ = strconv.ParseUint(ip[0], 10, 8)

// import (
// 	"fmt"
// 	"log"
// "strconv"
// 	"strings"

// 	// "golang.org/x/sys/unix"
// 	"unix"
// )

// var (

// 	MAXMSGSIZE = 8000
// )

// func main() {
// 	args := os.Args[1:]
// 	if len(args) != 2 {
// 		fmt.Println("./client [IPv4] [Port]")
// 	}

// 	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)

// 	if err != nil {
// 		log.Fatal("Socket: ", err)
// 	}

// 	port, err := strconv.Atoi(args[1])

// 	if err != nil || (port < 0 || port > 100000) {
// 		os.Stderr.WriteString("Invalid port format\n")

// 		return
// 	}

// 	serverAddr := &unix.SockaddrInet4{
// 		Port: port,
// 		Addr: inetAddr(args[0]),
// 	}

// 	err = unix.Connect(serverFD, serverAddr)
// 	if err != nil {

// 		if err == unix.ECONNREFUSED {
// 			fmt.Println("* Connection failed")
// 			unix.Close(serverFD)
// 			return
// 		}
// 	}

// 	var msg string
// 	var response []byte

// 	response = make([]byte, MAXMSGSIZE)

// 	print("> ")
// 	fmt.Scanln(&msg)
// 	err = unix.Sendmsg(
// 		serverFD,
// 		[]byte(msg),
// 		nil, serverAddr, unix.MSG_DONTWAIT)
// 	if err != nil {
// 		fmt.Println("Sendmsg: ", err)
// 	}
// 	_, _, err = unix.Recvfrom(serverFD, response, 0)
// 	if err != nil {
// 		fmt.Println("Recvfrom: ", err)
// 		unix.Close(serverFD)
// 		return
// 	}
// 	fmt.Printf("< %s\n", string(response))
// 	unix.Close(serverFD)
// 	return
// }

// func inetAddr(ipaddr string) [4]byte {
// 	var (
// 		ip                 = strings.Split(ipaddr, ".")
// 		ip1, ip2, ip3, ip4 uint64
// 	)
// 	ip1, _ = strconv.ParseUint(ip[0], 10, 8)
// 	ip2, _ = strconv.ParseUint(ip[1], 10, 8)
// 	ip3, _ = strconv.ParseUint(ip[2], 10, 8)
// 	ip4, _ = strconv.ParseUint(ip[3], 10, 8)
// 	return [4]byte{byte(ip1), byte(ip2), byte(ip3), byte(ip4)}
// }
