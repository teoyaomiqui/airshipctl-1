/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package document

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipctl/pkg/environment"
)

// NewDocumentPullCommand creates a new command for pulling airship document repositories
func NewDocumentPullCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	settings := pull.Settings{AirshipCTLSettings: rootSettings}
	documentPullCmd := &cobra.Command{
		Use:   "pull",
		Short: "pulls documents from remote git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return settings.Pull()
		},
	}

	return documentPullCmd
}
