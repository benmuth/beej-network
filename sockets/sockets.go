package sockets

// equivalent to the C struct for socket address structures and lookups
type AddrInfo struct {
	ai_flags     int32     // AI_PASSIVE, AI_CANONNAME, etc.
	ai_family    int32     // AF_INET, AF_INET6, AF_UNSPEC
	ai_socktype  int32     // SOCK_STREAM, SOCK_DGRAM
	ai_protocol  int32     // use 0 for "any"
	ai_addrlen   uint64    // size of ai_addr in bytes
	ai_addr      *SockAddr // struct sockaddr_in or _in6
	ai_canonname *byte     // full canonical hostname

	ai_next *AddrInfo // linked list, next node
}

// SockAddr Holds socket address information for many types of sockets
type SockAddr struct {
	SAFamily uint16   // address family, AF_xxx, AF_INET (IPv4) or AF_INET6 (IPv6)
	SAData   [14]byte // protocol address: destination address and port number
}

// SockAddrIn can be cast to SockAddr in C
// SockAddrIn was created so you don't have to pack SockAddr.SAData by hand
// In stands for internet.
// (IPv4 only--see struct sockaddr_in6 for IPv6)
type SockAddrIn struct {
	SinFamily int16    // Address family, AF_INET, corresponds to SockAddr.SAFamily
	SinPort   uint16   // Port number, should be BigEndian
	SinAddr   InAddr   // Internet address
	SinZero   [8]uint8 // Same size as struct sockaddr, padding, should be set to 0
}

// (IPv4 only--see struct in6_addr for IPv6)
// Internet address (a structure for historical reasons)
type InAddr struct {
	SAddr uint32 // that's a 32-bit int (4 bytes)
}

// SockAddrStorage is large enough to hold both IPV4 and IPV6 structures
// useful if you don't know ahead of time what a function will need
type SockAddrStorage struct {
	ss_family SAFamilyT // address family
	// all this is padding, implementation specific, ignore it:
	__ss_pad1  [_SS_PAD1SIZE]byte
	__ss_align int64
	__ss_pad2  [_SS_PAD2SIZE]byte
}
