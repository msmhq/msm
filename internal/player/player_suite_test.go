package player_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlayer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Player Suite")
}
