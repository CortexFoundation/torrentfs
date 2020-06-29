// Copyright 2020 The CortexTheseus Authors
// This file is part of the CortexTheseus library.
//
// The CortexTheseus library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The CortexTheseus library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the CortexTheseus library. If not, see <http://www.gnu.org/licenses/>.
package torrentfs

import (
	"fmt"
	"github.com/CortexFoundation/CortexTheseus/log"
	"github.com/CortexFoundation/CortexTheseus/p2p"
	"github.com/CortexFoundation/CortexTheseus/rlp"
	mapset "github.com/ucwong/golang-set"
	"sync"
)

type Peer struct {
	host *TorrentFS
	peer *p2p.Peer
	ws   p2p.MsgReadWriter

	trusted bool

	known mapset.Set // Messages already known by the peer to avoid wasting bandwidth
	quit  chan struct{}

	wg sync.WaitGroup
}

func newPeer(host *TorrentFS, remote *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	return &Peer{
		host:    host,
		peer:    remote,
		ws:      rw,
		trusted: false,
		known:   mapset.NewSet(),
		quit:    make(chan struct{}),
	}
}

func (p *Peer) Start() error {
	return nil
}

func (peer *Peer) handshake() error {
	log.Info("Nas handshake", "peer", peer.ID())
	errc := make(chan error, 1)
	peer.wg.Add(1)
	go func() {
		defer peer.wg.Done()

		errc <- p2p.SendItems(peer.ws, statusCode, ProtocolVersion, 0, nil, false)
	}()
	if err := <-errc; err != nil {
		return fmt.Errorf("peer [%x] failed to send status packet: %v", peer.ID(), err)
	}
	// Fetch the remote status packet and verify protocol match
	packet, err := peer.ws.ReadMsg()
	if err != nil {
		return err
	}
	if packet.Code != statusCode {
		return fmt.Errorf("peer [%x] sent packet %x before status packet", peer.ID(), packet.Code)
	}
	s := rlp.NewStream(packet.Payload, uint64(packet.Size))
	_, err = s.List()
	if err != nil {
		return fmt.Errorf("peer [%x] sent bad status message: %v", peer.ID(), err)
	}
	peerVersion, err := s.Uint()
	if err != nil {
		return fmt.Errorf("peer [%x] sent bad status message (unable to decode version): %v", peer.ID(), err)
	}
	if peerVersion != ProtocolVersion {
		return fmt.Errorf("peer [%x]: protocol version mismatch %d != %d", peer.ID(), peerVersion, ProtocolVersion)
	}

	if err := <-errc; err != nil {
		return fmt.Errorf("peer [%x] failed to send status packet: %v", peer.ID(), err)
	}
	log.Info("Nas p2p hanshake success", "id", peer.ID(), "status", packet.Code, "version", peerVersion)
	return nil
}

func (p *Peer) Stop() error {
	close(p.quit)
	p.wg.Wait()
	return nil
}
func (peer *Peer) ID() []byte {
	id := peer.peer.ID()
	return id[:]
}
