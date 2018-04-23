package bark

import (
	"expvar"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/juju/ratelimit"
)

var (
	pps = expvar.NewInt("pps")
)

type Client struct {
	Host      string
	Port      int
	TargetPPS int64
	Generator Generator
	stop      chan struct{}
	conn      net.Conn
}

// NewClient returns a new client
func NewClient(host string, port int, targetPPS int64, generator Generator) *Client {
	return &Client{
		Host:      host,
		Port:      port,
		TargetPPS: targetPPS,
		Generator: generator,
		stop:      make(chan struct{}),
	}
}

// Bark connects to the targeted host and start sending metrics
func (c *Client) Bark() error {
	conn, err := net.Dial("udp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
	if err != nil {
		return err
	}

	c.conn = conn
	go c.barkLoop()

	return nil
}

func (c *Client) Stop() {
	close(c.stop)
}

func (c *Client) Wait() {
	<-c.stop
}

func (c *Client) barkLoop() {
	limiter := ratelimit.NewBucketWithRate(float64(c.TargetPPS), c.TargetPPS)
	logTicker := time.NewTicker(5 * time.Second)
	count := 0
	retry := 0

	defer logTicker.Stop()
	defer c.conn.Close()

	for {
		select {
		case <-c.stop:
			return
		case <-logTicker.C:
			pps.Set(int64(count / 5))
			count = 0
		default:
			limiter.Wait(1)
			_, err := c.conn.Write(c.Generator.NextDatagram())
			count++
			if err != nil {
				fmt.Println("Bark worker error: ", err)
				retry++
				time.Sleep(5 * time.Second)
				if retry == 5 {
					fmt.Println("Could not connect after five retry, exiting.")
					c.Stop()
				}
			} else {
				retry = 0
			}
		}
	}
}
