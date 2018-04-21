package bark

type Generator interface {
	NextDatagram() []byte
}
