package BTClient

import "net"

/*构建peer的struct*/
type Peer struct {
	IP       net.IP     `ip`
	PeerID   [20]byte 	`peer id`
	Port     uint16		`port`
}
