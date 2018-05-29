package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/arbll/barkstatsd/pkg/bark"
	"github.com/spf13/cobra"
)

var barkCmd = &cobra.Command{
	Use:   "barkstatsd",
	Short: "A statsd & dogstatsd load testing tool for benchmark purposes.",
	RunE:  doBark,
}

var (
	host     string
	port     int
	pps      int64
	step     int64
	interval time.Duration
	duration time.Duration
)

func init() {
	barkCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Host to bark at")
	barkCmd.PersistentFlags().IntVarP(&port, "port", "p", 8125, "Port to bark at")
	barkCmd.PersistentFlags().Int64VarP(&pps, "pps", "r", 1000, "Initial PPS")
	barkCmd.PersistentFlags().Int64VarP(&step, "step", "s", 0, "PPS step")
	barkCmd.PersistentFlags().DurationVarP(&interval, "step-interval", "i", 1*time.Minute, "Step interval")
	barkCmd.PersistentFlags().DurationVarP(&duration, "duration", "d", 0, "Duration (0 for infinite)")
}

func Bark() {
	if err := barkCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doBark(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires path to a packet file")
	}

	packets, err := readPackets(args[0])
	if err != nil {
		return err
	}

	generator := bark.DogStatsD{
		Packets: packets,
	}
	client := bark.NewClient(host, port, pps, step, interval, duration, &generator)
	client.Bark()
	client.Wait()
	return nil
}

func readPackets(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n\n"), nil
}
