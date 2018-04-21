package command

import (
	"bufio"
	"fmt"
	"os"

	"github.com/arbll/barkstatsd/pkg/bark"
	"github.com/spf13/cobra"
)

var barkCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	RunE: doBark,
}

var (
	host string
	port int
	pps  int64
)

func init() {
	barkCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Host to bark at")
	barkCmd.PersistentFlags().IntVarP(&port, "port", "p", 8125, "Port to bark at")
	barkCmd.PersistentFlags().Int64VarP(&pps, "pps", "r", 1000, "Target PPS")
}

func Bark() {
	if err := barkCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doBark(cmd *cobra.Command, args []string) error {

	generator := bark.DogStatsD{}
	client := bark.NewClient(host, port, pps, &generator)
	client.Bark()
	buf := bufio.NewReader(os.Stdin)
	buf.ReadBytes('\n')
	return nil
}
