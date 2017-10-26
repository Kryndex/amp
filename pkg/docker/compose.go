package docker

import (
	"os"
	"strings"

	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/cli/cli/compose/loader"
	composetypes "github.com/appcelerator/amp/docker/cli/cli/compose/types"
	"golang.org/x/net/context"
)

// ComposeParse parses a compose file
func (d *Docker) ComposeParse(ctx context.Context, composeFile []byte, environment []string) (*composetypes.Config, error) {
	var details composetypes.ConfigDetails
	var err error

	// WorkingDir
	details.WorkingDir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Parsing compose file
	config, err := loader.ParseYAML(composeFile)
	if err != nil {
		return nil, err
	}
	details.ConfigFiles = []composetypes.ConfigFile{{
		Filename: "filename",
		Config:   config,
	}}
	details.Environment, err = stack.BuildEnvironment(environment)
	return loader.Load(details)
}

// ComposeIsAuthorized checks if the given compose file is authorized
func (d *Docker) ComposeIsAuthorized(compose *composetypes.Config) bool {
	for _, reservedSecret := range constants.Secrets {
		if _, exists := compose.Secrets[reservedSecret]; exists {
			return false
		}
	}

	for _, service := range compose.Services {
		for _, reservedSecret := range constants.Secrets {
			for _, secret := range service.Secrets {
				if strings.EqualFold(secret.Source, reservedSecret) {
					return false
				}
			}
		}
		for _, reservedLabel := range constants.Labels {
			if _, exists := service.Labels[reservedLabel]; exists {
				return false
			}
			if _, exists := service.Deploy.Labels[reservedLabel]; exists {
				return false
			}
		}
	}
	return true
}
