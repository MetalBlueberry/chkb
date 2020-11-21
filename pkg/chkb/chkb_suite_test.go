package chkb_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChkb(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetOutput(GinkgoWriter)
	RunSpecs(t, "Chkb Suite")
}
