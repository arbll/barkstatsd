package bark

import (
	"math/rand"
)

type DogStatsD struct {
	Seed int64
	rand *rand.Rand
}

// NextDatagram generates the next dogstatsd datagram
func (g *DogStatsD) NextDatagram() []byte {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(g.Seed))
	}
	return g.nextMetric()
}

func (g *DogStatsD) nextMetric() []byte {
	return []byte("barkstatsd.metric:1|c")
}
