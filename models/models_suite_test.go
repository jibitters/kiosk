package models_test

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pgHost string

func init() {
	flag.StringVar(&pgHost, "pg.host", "localhost", "")
}

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}
