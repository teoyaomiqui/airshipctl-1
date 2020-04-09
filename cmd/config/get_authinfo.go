/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/environment"
)

var (
	getAuthInfoLong = `Display a specific user's information, or all defined users if no name is provided`

	getAuthInfoExample = `# List all the users airshipctl knows about
airshipctl config get-credentials

# Display a specific user information
airshipctl config get-credentials e2e`
)

// NewCmdConfigGetAuthInfo returns a Command instance for 'config -AuthInfo' sub command
// An AuthInfo refers to a particular user for a cluster
func NewCmdConfigGetAuthInfo(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	o := &config.AuthInfoOptions{}
	cmd := &cobra.Command{
		Use:     "get-credentials NAME",
		Short:   "Gets a user entry from the airshipctl config",
		Long:    getAuthInfoLong,
		Example: getAuthInfoExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				o.Name = args[0]
			}
			return config.RunGetAuthInfo(o, cmd.OutOrStdout(), rootSettings.Config())
		},
	}

	return cmd
}
