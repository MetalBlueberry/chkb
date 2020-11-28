package chkb_test

import (
	"syscall"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestChkb(t *testing.T) {
	RegisterFailHandler(Fail)
	logrus.SetOutput(GinkgoWriter)
	RunSpecs(t, "Chkb Suite")
}

func Time(t time.Time) syscall.Timeval {
	return syscall.Timeval{
		Sec:  t.Unix(),
		Usec: int64(t.UnixNano() / 1000 % 1000000),
	}
}

func Elapsed(ms int64) time.Time {
	return time.
		Date(2020, 11, 20, 12, 0, 0, 0, time.UTC).
		Add(time.Duration(ms) * time.Millisecond)
}
