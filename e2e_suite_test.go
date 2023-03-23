package inchori_e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	manager *Manager
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("set up e2e test suite...")

	var err error
	s.manager, err = NewManager("local-panacea")
	s.Require().NoError(err)
}
