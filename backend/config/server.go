package config

import (
	"errors"
	"os"
)

var (
	Secret               = os.Getenv("SECRET")
	HasuraApiUrl         = os.Getenv("HASURA_API_URL")
	HasuraApiToken       = os.Getenv("HASURA_API_TOKEN")
	QoveryApiToken       = os.Getenv("API_TOKEN_QOVERY")
	QoveryOrganizationId = os.Getenv("ORGANIZATION_ID_QOVERY")
)

func CheckServerConfig() []error {
	var configErrors []error

	if Secret == "" {
		configErrors = append(configErrors, errors.New("SECRET env required"))
	}
	if HasuraApiUrl == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_URL env required"))
	}
	if HasuraApiToken == "" {
		configErrors = append(configErrors, errors.New("HASURA_API_TOKEN env required"))
	}
	if QoveryApiToken == "" {
		configErrors = append(configErrors, errors.New("API_TOKEN_QOVERY env required"))
	}
	if QoveryOrganizationId == "" {
		configErrors = append(configErrors, errors.New("ORGANIZATION_ID_QOVERY env required"))
	}

	return configErrors
}
