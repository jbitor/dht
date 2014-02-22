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
	QueriedNodes       map[bittorrent.BTID]*RemoteNode
	PeersFound         map[bittorrent.BTID]*bittorrent.RemotePeer
	OutstandingQueries map[string]*RpcQuery
	StartTime          time.Time
	finished           bool
}

// XXX(JB): We probably won't actually use this anywhere, instances
// XXX(JB): will just be created in the one location that they're
// XXX(JB): actually needed.
func newGetPeersSearch(infohash bittorrent.BTID) (s *GetPeersSearch) {
	s = &GetPeersSearch{
		Infohash:  infohash,
		Options:   GetPeersDefaultOptions,
		StartTime: time.Now(),
	}
	return
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
