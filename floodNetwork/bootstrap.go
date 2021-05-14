package floodNetwork

import (
	"encoding/gob"
	"io"
	"net"
	"time"
)

// listens for new TCP connections
func (n *Network) listenForConnections() {
	for {
		c, err := n.server.Accept()
		if err != nil {
			n.log(err.Error())
			continue
		}
		go n.handleConnection(c, nil)
	}
}

// handleConnection is either waiting for a packet, or handles a new peers connection.
// if newPeeer is nil, then we await a packet, otherwise we just handle that peer.
// for tmp connections before they become peers, only accepts one packet, then closes.
func (n *Network) handleConnection(conn net.Conn, newPeer *Peer) {
	n.log("handling TCP connection: " + conn.RemoteAddr().String())

	if newPeer != nil {
		n.handlePeer(newPeer)
	} else {
		dec := gob.NewDecoder(conn)
		packet := Packet{}
		err := dec.Decode(&packet) // blocking till we finish reading message

		if err == io.EOF { // client disconnected
			// do nothing, closing connection later
		} else if err != nil { // error decoding message
			n.log(err.Error())
		} else { // no errors, handle packet
			switch packet.Type {
			case CONN_REQ:
				n.recieveConnectionRequest(packet)
			case CONN_ACK:
				n.recieveConnectionAcknowledgment(conn, packet)
			}
		}
	}

	n.log("stopped handling TCP connection: " + conn.RemoteAddr().String())
	conn.Close()
}

func (n *Network) recieveConnectionAcknowledgment(conn net.Conn, packet Packet) {
	// double check
	if packet.Type != CONN_ACK {
		n.log("invalid function call, cannot handle packet not of type CONN_ACK")
		return
	}

	n.log("got P2P connection acknowledge from " + packet.Origin)

	for peer := range n.Peers {
		if packet.Origin == peer.Address {
			n.log("already P2P connected to " + packet.Origin + ", ACK rejected")
			return
		}
	}

	newPeer := Peer{
		Connection: conn,
		Address:    packet.Origin,
	}
	n.handlePeer(&newPeer)
}

// creates the connection to a machine
func (n *Network) requestConnection(destinationAddr string) (net.Conn, bool) {
	n.log("requesting TCP connection to " + destinationAddr)

	// verify you can connect
	if destinationAddr == n.localAddress {
		n.log("Tried to TCP connect to self, request rejected")
		return nil, false
	}
	for peer := range n.Peers {
		if destinationAddr == peer.Address {
			n.log("Already TCP connected to " + destinationAddr + ", request rejected")
			return nil, false
		}
	}

	conn, err := net.Dial("tcp4", destinationAddr)
	if err != nil {
		n.log(err.Error())
		return nil, false
	}
	n.log("TCP connection established with " + destinationAddr)
	return conn, true
}

// creates connection, sends request, then closes. We will get a new connection if someone accepts
// should only be used by a node not connected to any nodes, otherwise send request through peers
func (n *Network) enterNetwork(bootstrapIP string) {
	tmpConn, ok := n.requestConnection(bootstrapIP)
	if !ok {
		return
	}
	n.sendConnReq(tmpConn)
	tmpConn.Close()
	n.log("closed TCP connection to " + bootstrapIP + "\n")
}

func (n *Network) sendAck(c net.Conn) {
	n.log("sending CONN_ACK to " + c.RemoteAddr().String())
	ack := Packet{
		Type:      CONN_ACK,
		Origin:    n.localAddress,
		Timestamp: time.Now().String(),
	}
	n.sendPacket(c, ack)
}

func (n *Network) sendConnReq(c net.Conn) {
	n.log("sending CONN_REQ to " + c.RemoteAddr().String())
	connReq := Packet{
		Type:      CONN_REQ,
		Origin:    n.localAddress,
		Timestamp: time.Now().String(),
	}
	n.sendPacket(c, connReq)
}
