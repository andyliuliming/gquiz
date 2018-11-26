package gquiz_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGquiz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gquiz Suite")
}
