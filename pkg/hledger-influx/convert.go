package hledger

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
)

func Convert(in io.Reader, out io.Writer) error {
	var (
		header []string
		enc    lineprotocol.Encoder
	)

	r := csv.NewReader(in)

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
			return err
		}

		if header == nil {
			header = record
			continue
		}

		t, err = time.Parse("2006-01-02", record[0])
		if err != nil {
			return err
		}

		for i := 0; i < len(header)-1; i++ {
			if i == 0 {
				continue // date
			}

			if string(record[i][0]) == "$" {
				v, err = strconv.ParseFloat(strings.Replace(record[i][1:], ",", "", -1), 64)
				if err != nil {
					return err
				}
			} else {
				v, err = strconv.ParseFloat(strings.Replace(record[i], ",", "", -1), 64)
				if err != nil {
					return err
				}
			}

			enc.StartLine("account")
			enc.AddTag("name", header[i])
			enc.AddField("balance", lineprotocol.MustNewValue(v))
			enc.EndLine(t)
			if err := enc.Err(); err != nil {
				return fmt.Errorf("encoding error: %v", err)
			}
		}
	}

	fmt.Fprintf(out, "%s", enc.Bytes())

	return nil
}
