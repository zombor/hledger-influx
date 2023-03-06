package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
)

func main() {
	var (
		header []string
		enc    lineprotocol.Encoder
	)

	r := csv.NewReader(bufio.NewReader(os.Stdin))

	enc.SetPrecision(lineprotocol.Microsecond)

	for {
		var (
			t time.Time
			v float64
		)

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if header == nil {
			header = record
			continue
		}

		t, err = time.Parse("2006-01-02", record[0])
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(header)-1; i++ {
			if i == 0 {
				continue // date
			}

			if string(record[i][0]) == "$" {
				v, err = strconv.ParseFloat(strings.Replace(record[i][1:], ",", "", -1), 64)
				if err != nil {
					panic(err)
				}
			} else {
				v, err = strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
				if err != nil {
					panic(err)
				}
			}

			enc.StartLine("account")
			enc.AddTag("name", header[i])
			enc.AddField("balance", lineprotocol.MustNewValue(v))
			enc.EndLine(t)
			if err := enc.Err(); err != nil {
				panic(fmt.Errorf("encoding error: %v", err))
			}
		}
	}

	fmt.Printf("%s", enc.Bytes())
}
