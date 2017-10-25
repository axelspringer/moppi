package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMoppi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Moppi Suite")
}
