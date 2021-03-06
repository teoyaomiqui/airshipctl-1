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

package redfish

import (
	"fmt"

	aerror "opendev.org/airship/airshipctl/pkg/errors"
)

// ErrRedfishClient describes an error encountered by the go-redfish client.
type ErrRedfishClient struct {
	aerror.AirshipError
	Message string
}

func (e ErrRedfishClient) Error() string {
	return fmt.Sprintf("redfish client encountered an error: %s", e.Message)
}

// ErrRedfishMissingConfig describes an error encountered due to a missing configuration option.
type ErrRedfishMissingConfig struct {
	What string
}

func (e ErrRedfishMissingConfig) Error() string {
	return "missing configuration: " + e.What
}

// ErrOperationRetriesExceeded raised if number of operation retries exceeded
type ErrOperationRetriesExceeded struct {
	What    string
	Retries int
}

func (e ErrOperationRetriesExceeded) Error() string {
	return fmt.Sprintf("operation %s failed. Maximum retries (%d) exceeded", e.What, e.Retries)
}
