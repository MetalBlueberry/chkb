/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package chkb_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/MetalBlueberry/chkb/pkg/chkb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var (
	TapDelayMs = chkb.DefaultTapDelay.Milliseconds()

	BeforeTap = TapDelayMs - 50
	AfterTap  = TapDelayMs + 50
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

var InitialTime = time.Date(2020, 11, 20, 12, 0, 0, 0, time.UTC)

func Elapsed(ms int64) time.Time {
	return InitialTime.
		Add(time.Duration(ms) * time.Millisecond)
}
