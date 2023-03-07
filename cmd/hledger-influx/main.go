package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/zombor/hledger-influx/pkg/hledger-influx"
)

//go:embed VERSION.txt
var version string

func main() {
	fs := flag.NewFlagSet("hledger-influx", flag.ContinueOnError)

	version := &ffcli.Command{
		Name:       "version",
		ShortUsage: "hledger-influx version [<arg> ...]",
		ShortHelp:  "Prints the program version.",
		Exec:       func(context.Context, []string) error { fmt.Println(version); return nil },
	}

	root := &ffcli.Command{
		ShortUsage: "hledger-influx",
		ShortHelp:  "Converts an hledger csv file of account balances to influxdb line format",
		LongHelp: `Converts an hledger csv file of account balances to influxdb line format.
Use the 'bal', '-O csv', and '-DH' options of 'hledger bal' to print a csv report in the right format for this program:
hledger bal --infer-market-prices -V -X=$ not:tag:clopen -O csv -DH --transpose | hledger-influx | influx write ...
		`,
		Subcommands: []*ffcli.Command{version},
		FlagSet:     fs,
		Exec:        func(context.Context, []string) error { return hledger.Convert(bufio.NewReader(os.Stdin), os.Stdout) },
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
