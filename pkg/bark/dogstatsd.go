package bark

import (
	"math/rand"
)

type DogStatsD struct {
	Seed    int64
	Packets []string
	i       int
	rand    *rand.Rand
}

// NextDatagram generates the next dogstatsd datagram
func (g *DogStatsD) NextDatagram() []byte {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(g.Seed))
	}
	return g.nextMetric()
}

func (g *DogStatsD) nextMetric() []byte {
	packet := g.Packets[g.i]
	g.i = (g.i + 1) % len(g.Packets)

	return []byte(packet)
}
