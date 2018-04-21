package bark

import (
	"math/rand"
	"strings"
)

type DogStatsD struct {
	Seed int64
	rand *rand.Rand
}

var (
	metricNames  = [...]string{"page.views", "fuel.level", "song.length", "users.uniques", "users.online"}
	metricValues = [...]string{"1", "0.5", "240", "1234", "1337", "19932242.324"}
	metricType   = [...]string{"c", "g", "ms", "h", "s"}
	sampleRates  = [...]string{"@0", "@0.5", "@1"}
	tags         = [...]string{"#country:china,contry:france", "#tagname:test", "#tag1:value,tag2:value,tag3:value,tag4:value,tag5:value"}
)

// NextDatagram generates the next dogstatsd datagram
func (g *DogStatsD) NextDatagram() []byte {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(g.Seed))
	}
	return g.nextMetric()
}

func (g *DogStatsD) nextMetric() []byte {
	var metricBuilder strings.Builder
	metricName := metricNames[g.rand.Int31n(int32(len(metricNames)))]
	metricValue := metricValues[g.rand.Int31n(int32(len(metricValues)))]
	metricType := metricType[g.rand.Int31n(int32(len(metricType)))]
	tag := tags[g.rand.Int31n(int32(len(tags)))]

	metricBuilder.WriteString(metricName)
	metricBuilder.WriteByte(':')
	metricBuilder.WriteString(metricValue)
	metricBuilder.WriteByte('|')
	metricBuilder.WriteString(metricType)
	metricBuilder.WriteByte('|')
	if metricType == "h" || metricType == "s" {
		sampleRate := sampleRates[g.rand.Int31n(int32(len(sampleRates)))]
		metricBuilder.WriteString(sampleRate)
		metricBuilder.WriteByte('|')
	}
	metricBuilder.WriteString(tag)

	return []byte(metricBuilder.String())
}
