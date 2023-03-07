package hledger_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/zombor/hledger-influx/pkg/hledger-influx"
)

var _ = Describe("Convert", func() {
	var (
		in  io.Reader
		out bytes.Buffer

		err error
	)

	BeforeEach(func() {
		in = strings.NewReader(`"account","test","total"
"2023-01-01","$1,000.00","$1,000.00"`)
	})

	JustBeforeEach(func() {
		err = Convert(in, &out)
	})

	It("Writes the csv to out", func() {
		data, _ := ioutil.ReadAll(&out)
		Expect(string(data), err).To(Equal("account,name=test balance=1000 1672531200000000\n"))
	})

	Describe("When the date is malformed", func() {
		BeforeEach(func() {
			in = strings.NewReader(`"account","test","total"
"not-a-date","$1,000.00","$1,000.00"`)
		})

		It("Returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("When the amount is malformed", func() {
		BeforeEach(func() {
			in = strings.NewReader(`"account","test","total"
"2023-01-01","not-an-amount","$1,000.00"`)
		})

		It("Returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("When the amount is malformed", func() {
		BeforeEach(func() {
			in = strings.NewReader(`"account","test","total"
"2023-01-01","$not-an-amount","$1,000.00"`)
		})

		It("Returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("With two accounts", func() {
		BeforeEach(func() {
			in = strings.NewReader(`"account","test","test2","total"
"2023-01-01","$1,000.00","$1,100.00","$1,000.00"`)
		})

		It("Writes the csv data to out", func() {
			data, _ := ioutil.ReadAll(&out)
			Expect(string(data), err).To(Equal("account,name=test balance=1000 1672531200000000\naccount,name=test2 balance=1100 1672531200000000\n"))
		})

		Describe("When the columns don't match up", func() {
			BeforeEach(func() {
				in = strings.NewReader(`"account","test","test2","total"
"2023-01-01","$1,000.00","$1,000.00"`)
			})

			It("Returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
