package gimple_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGimple(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gimple Suite")
}
