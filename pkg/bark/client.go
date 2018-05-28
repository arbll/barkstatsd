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
	pps    = expvar.NewInt("pps")
	pcount = expvar.NewInt("pcount")
)

type Client struct {
	Host         string
	Port         int
	TargetPPS    int64
	StepPPS      int64
	IntervalStep time.Duration
	Duration     time.Duration
	Generator    Generator
	stop         chan struct{}
	conn         net.Conn
}

// NewClient returns a new client
func NewClient(host string, port int, targetPPS, stepPPS int64, stepInterval time.Duration, duration time.Duration, generator Generator) *Client {
	return &Client{
		Host:         host,
		Port:         port,
		TargetPPS:    targetPPS,
		StepPPS:      stepPPS,
		IntervalStep: stepInterval,
		Duration:     duration,
		Generator:    generator,
		stop:         make(chan struct{}),
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
	targetPPS := c.TargetPPS
	limiter := ratelimit.NewBucketWithRate(float64(targetPPS), targetPPS)
	logTicker := time.NewTicker(5 * time.Second)
	stepTicker := time.NewTicker(c.IntervalStep)
	duration := c.Duration
	if duration == 0 {
		duration = time.Hour * 24 * 365
	}
	timerStop := time.NewTimer(duration)
	count := 0
	retry := 0

	defer stepTicker.Stop()
	defer logTicker.Stop()
	defer c.conn.Close()
	defer timerStop.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-timerStop.C:
			c.Stop()
		case <-logTicker.C:
			pps.Set(int64(count / 5))
			count = 0
		case <-stepTicker.C:
			if c.StepPPS > 0 {
				targetPPS += c.StepPPS
				limiter = ratelimit.NewBucketWithRate(float64(targetPPS), targetPPS)
			}
		default:
			limiter.Wait(1)
			_, err := c.conn.Write(c.Generator.NextDatagram())
			if err != nil {
				fmt.Println("Bark worker error: ", err)
				retry++
				time.Sleep(5 * time.Second)
				if retry == 5 {
					fmt.Println("Could not connect after five retry, exiting.")
					c.Stop()
				}
			} else {
				count++
				pcount.Add(1)
				retry = 0
			}
		}
	}
}
