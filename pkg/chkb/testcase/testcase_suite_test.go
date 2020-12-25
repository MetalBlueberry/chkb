package testcase_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTestcase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testcase Suite")
}
