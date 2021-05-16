package floodNetwork

import (
	"encoding/gob"
	"io"
	"math/rand"
	"net"
	"time"
)

type PacketType byte

const (
	MESSAGE PacketType = iota
	CONN_REQ
	CONN_ACK
)

type Packet struct {
	Type      PacketType
	Payload   []byte // arbitrary data type
	Timestamp string // so we can send the same data twice, and it not get rejected
}

// asynchronous function, a different instance is run for each peer
func (n *Network) handlePeer(peer *Peer) {

	n.waitPeers.Add(1)
	n.addPeerChan <- peer
	n.waitPeers.Wait()
	n.log("handling P2P connection: " + peer.Address)

	// dont return in this loop, have some cleaning up to do afterward
	for {
		dec := gob.NewDecoder(peer.Connection)
		packet := &Packet{}
		err := dec.Decode(packet) // blocking till we finish reading message

		if err == io.EOF { // client disconnected
			break
		} else if err != nil { // error decoding message
			n.log(err.Error())
			continue
		}

		// no errors, handle packet
		n.recievePacket(*packet)
	}

	n.log("stopped handling P2P connection: " + peer.Address)

	n.waitPeers.Add(1)
	n.removePeerChan <- peer // update the peer list
	n.waitPeers.Wait()
}

func (n *Network) recievePacket(packet Packet) {
	// check we havent seen this packet before (may not always be a good idea, probably have to change later)
	for _, oldPacket := range n.recentPackets {
		if oldPacket.Timestamp == packet.Timestamp { // should probably have better way of checking this
			return
		}
	}
	// then add it so we dont handle again
	n.recentPackets = append(n.recentPackets, packet)

	switch packet.Type {
	case MESSAGE:
		n.recieveMessage(packet)
	case CONN_REQ:
		n.recieveConnectionRequest(packet)
	}
}

func (n *Network) recieveMessage(packet Packet) {
	// make message available to handle outside of library
	n.alertMessage(packet.Payload)
	n.announcePacket(packet)
}

func (n *Network) recieveConnectionRequest(packet Packet) {
	// double check
	if packet.Type != CONN_REQ {
		n.log("invalid function call, cannot handle packet not of type CONN_REQ")
		return
	}

	address := string(packet.Payload)

	// get random peer
	var peerToPassTo *Peer = nil
	pickedIndex := rand.Intn(2) // 50% chance we pass
	if pickedIndex != 0 && len(n.peers) > 0 {
		pickedIndex = rand.Intn(len(n.peers))
		i := 0
		for peer := range n.peers {
			if i == pickedIndex {
				peerToPassTo = peer
			}
			i++
		}
	}

	if peerToPassTo == nil {
		if address == n.localAddress {
			n.log("cannot request connection to self, throwing out CONN_REQ")
		} else {
			n.log("got P2P connection request from " + address + ", accepting")

			conn, ok := n.requestConnection(address)
			if ok {
				newPeer := Peer{
					Connection: conn,
					Address:    address,
				}
				n.sendAck(conn) // let them know they are a peer now
				go n.handleConnection(conn, &newPeer)
			}
		}
	} else {
		n.log("got P2P connection request from " + address + ", forwarding to " + peerToPassTo.Address)
		packet.Timestamp = time.Now().String() // this makes sure we dont ignore the packet if we get sent it again
		n.sendPacket(peerToPassTo.Connection, packet)
	}
}

// sends packet to all peers
func (n *Network) announcePacket(packet Packet) {
	for peer := range n.peers {
		n.sendPacket(peer.Connection, packet)
	}
}

// sends packet to a peer
func (n *Network) sendPacket(connection net.Conn, packet Packet) {
	n.recentPackets = append(n.recentPackets, packet)

	encoder := gob.NewEncoder(connection)
	err := encoder.Encode(packet) // writes to tcp connection

	if err != nil {
		n.log(err.Error())
	}
}

func (n *Network) SendMessage(payload []byte) {
	msgPacket := Packet{
		Type:      MESSAGE,
		Payload:   payload,
		Timestamp: time.Now().String(),
	}

	n.announcePacket(msgPacket)
}
