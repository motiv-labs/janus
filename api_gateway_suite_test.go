package janus_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApiGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiGateway Suite")
}
