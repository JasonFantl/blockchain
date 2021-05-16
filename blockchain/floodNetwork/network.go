package floodNetwork

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	Connection net.Conn
	Address    string
}

type PeerList map[*Peer]bool

type Network struct {
	server        net.Listener
	peers         PeerList
	recentPackets []Packet

	localAddress string

	addPeerChan    chan *Peer
	removePeerChan chan *Peer
	waitPeers      sync.WaitGroup

	// callback functions
	log          func(string)
	alertMessage func([]byte)
}

// New creates a network object for you.
// Whenever a message is recieved, messageCallback be called, passing in that message.
func New(messageCallback func([]byte)) Network {
	return Network{
		peers:         make(map[*Peer]bool),
		recentPackets: make([]Packet, 0),

		addPeerChan:    make(chan *Peer),
		removePeerChan: make(chan *Peer),

		alertMessage: messageCallback,
		log:          func(s string) {}, // throws away debug by default
	}
}

// Whenever debug information needs to be printed, logger will be called, passing the debug message.
func (n *Network) SetLogger(logger func(string)) {
	n.log = logger
}

// Join causes the network object to attempt to connect to the actual network.
// myPort can be left as "" for default port.
// returns if we started a server or not. Check debug to see if we entered the network succsesfully
func (n *Network) Join(bootstrapIP, myPort string) bool {

	// init server and get local addr
	server, err := n.initServer(myPort)
	if err != nil {
		n.log(err.Error())
		return false
	}
	n.server = server

	// have to connect to self to get addr
	c, err := net.Dial("tcp4", server.Addr().String())
	n.localAddress = c.RemoteAddr().String()
	c.Close()
	n.log("Listening on: " + n.localAddress)

	go n.listenForConnections()
	go n.run()

	if bootstrapIP != "" {
		n.enterNetwork(bootstrapIP)
	}

	return true
}

func (n *Network) run() {
	defer n.server.Close()
	// was using for loop, but eats up CPU
	for {
		select {
		case newPeer := <-n.addPeerChan:
			n.peers[newPeer] = true
			n.log("added peer " + newPeer.Address)
			n.waitPeers.Done()
		case oldPeer := <-n.removePeerChan:
			_, ok := n.peers[oldPeer]
			if ok {
				delete(n.peers, oldPeer)
				n.log("removed peer " + oldPeer.Address)

				n.log("disconnected, sending out new CONN_REQ")
				connReq := Packet{
					Type:      CONN_REQ,
					Payload:   []byte(n.localAddress),
					Timestamp: time.Now().String(),
				}
				n.recieveConnectionRequest(connReq)
			}

			n.waitPeers.Done()
		}
	}
}

func (n *Network) initServer(port string) (net.Listener, error) {
	n.log("Initing server...")

	if port == "" {
		port = "1234"
	}
	port = ":" + port

	return net.Listen("tcp4", port)
}
