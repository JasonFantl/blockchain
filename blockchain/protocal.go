package blockchain

import (
	"encoding/json"
)

type PacketType byte

const (
	BLOCK PacketType = iota
	CHAIN
)

type Packet struct {
	Type    PacketType
	Payload []byte // arbitrary data type
}

func (bc *Blockchain) announcePacket(packetType PacketType, payload []byte) error {
	jsonPacket, err := json.Marshal(Packet{packetType, payload})
	if err != nil {
		return err
	}

	bc.network.SendMessage(jsonPacket)

	return nil
}

func (bc *Blockchain) announceChain() error {

	jsonList, err := bc.chain.ToBytes()
	if err != nil {
		return err
	}
	return bc.announcePacket(CHAIN, jsonList)
}

func (bc *Blockchain) recievePacket(data []byte) {
	var packet Packet
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return
	}

	if packet.Type == CHAIN {
		bc.chain.Update(packet.Payload)
	}
	if packet.Type == BLOCK {
	}
}
