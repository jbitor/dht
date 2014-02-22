package dht

import (
	"errors"
	"github.com/jbitor/bittorrent"
	"time"
)

type GetPeersOptions struct {
	// The search isn't finished until
	TargetPeers       int           // [we've found at least TargetPeers and
	MinNodes          int           // we've contacted at least MinNodes] or
	MaxNodes          int           // we've contacted at least MaxNodes or
	Timeout           time.Duration // this much time has passed.
	ConcurrentQueries int           // We may send up to this many queries at once.
}

var GetPeersDefaultOptions = GetPeersOptions{
	TargetPeers:       8,
	MinNodes:          4,
	MaxNodes:          32,
	Timeout:           10 * time.Minute,
	ConcurrentQueries: 8,
}

type GetPeersSearch struct {
	Infohash           bittorrent.BTID
	Options            GetPeersOptions
	QueriedNodes       map[string]*RemoteNode
	PeersFound         map[string]*bittorrent.RemotePeer
	OutstandingQueries map[string]*RpcQuery
	StartTime          time.Time
	finished           bool
	localNode          *localNode
}

// The main loop of a GetPeers search.
func (s *GetPeersSearch) run() {
	for !s.Finished() {
		nodes := s.localNode.NodesByCloseness(s.Infohash, false)

		if s.AdditionalQueriesDue() {
			var remote *RemoteNode = nil
			for _, candidate := range nodes {
				// XXX(JB): .String() is not a clean way to do this
				if _, ok := s.QueriedNodes[candidate.Address.String()]; ok {
					continue // already queried this node
				}

				remote = candidate
				s.QueriedNodes[remote.Address.String()] = remote
				break
			}

			if remote != nil {
				logger.Printf("Request peers for %v from %v.\n", s.Infohash, remote)

				go func() {
					peersResult, nodesResult, errorResult := s.localNode.GetPeers(remote, s.Infohash)

					select {
					case peers := <-peersResult:
						logger.Printf("Got peers from %v")
						for _, peer := range peers {
							// XXX(JB): .String() is not a clean way to do this
							s.PeersFound[peer.Address.String()] = peer
						}
					case _ = <-nodesResult:
						// nothing to do -- nodes will already have been recorded
					case err := <-errorResult:
						logger.Printf("Error response to GetPeers: %v\n", err)
					}
				}()
			}
		}

		time.Sleep(5 * time.Second)
	}
}

// Immediately causes the search to be finished.
// Any incomplete queries will be discarded.
func (s *GetPeersSearch) Terminate() {
	s.finished = true
}

// Whether this search is finished, with a result that will not change.
func (s *GetPeersSearch) Finished() bool {
	if s.finished {
		return s.finished
	}

	s.finished =
		time.Now().After(s.StartTime.Add(s.Options.Timeout)) ||
			(len(s.OutstandingQueries) == 0 &&
				len(s.QueriedNodes) >= s.Options.MinNodes &&
				(len(s.PeersFound) >= s.Options.TargetPeers ||
					len(s.QueriedNodes) >= s.Options.MaxNodes))

	return s.finished
}

// Whether this search is currently due to send any more queries.
func (s *GetPeersSearch) AdditionalQueriesDue() bool {
	return !s.Finished() &&
		len(s.OutstandingQueries) < s.Options.ConcurrentQueries
}

// Returns a channel which is notified each time new
// peers are recieved, and is closed when the request is finished.
func (s *GetPeersSearch) ReadNewPeers() (peersSource <-chan []bittorrent.RemotePeer) {
	return nil
}

// Blocks until we have any peers, then returns all known peers.
// Returns an error if the request fails to find any peers.
func (s *GetPeersSearch) AnyPeers() (peers []bittorrent.RemotePeer, err error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

// Blocks until the request is finished, then returns all known peers.
// Returns an error if the request fails to find any peers.
func (s *GetPeersSearch) AllPeers() (peers []bittorrent.RemotePeer, err error) {
	return nil, errors.New("NOT IMPLEMENTED")
}
