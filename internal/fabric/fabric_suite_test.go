package fabric_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFabric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fabric Suite")
}
