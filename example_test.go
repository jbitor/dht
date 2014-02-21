package dht_test

import (
	"fmt"
	"github.com/jbitor/jbitor/dht"
	"github.com/jbitor/jbitor/torrent"
	"testing"
)

// Attempts to find peers for an Ubuntu Torrent.
func ExampleClient() {
	infoHash, _ := torrent.BTIDFromHex("5497a53543938b77ef660939d3b32e02be7bc213")

	c, err := dht.OpenClient(".dht", true)
	if err != nil {
		panic(err)
		return
	}

	defer c.Close()

	peers, err := c.GetPeers(infoHash)
	if err != nil {
		fmt.Printf("Unable to find peers for %v.\n", infoHash)
		return
	}

	fmt.Printf("Found peers: %v.\n", peers)
}

func TestExamples(t *testing.T) {
	t.Skip("the examples don't actually work yet")
	ExampleClient()
}
