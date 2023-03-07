package hledger_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHledgerInflux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HledgerInflux Suite")
}
