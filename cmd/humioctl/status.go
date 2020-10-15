// Copyright © 2018 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Shows general status information",
		Args:  cobra.ExactArgs(0),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)
			serverStatus, serverErr := client.Status()
			if serverErr != nil {
				return nil, fmt.Errorf("error getting server status: %w", serverErr)
			}

			username, usernameErr := client.Viewer().Username()
			if usernameErr != nil {
				return nil, fmt.Errorf("error getting the current user: %w", usernameErr)
			}

			data := struct {
				Status   string
				Address  string
				Version  string
				Username string
			}{
				Status:   formatStatusText(serverStatus.Status),
				Address:  viper.GetString(viperkey.Address),
				Version:  serverStatus.Version,
				Username: username,
			}

			return data, nil
		}),
	}

	cmd.AddCommand(newLicenseInstallCmd())
	cmd.AddCommand(newLicenseShowCmd())

	return cmd
}

func formatStatusText(statusText string) string {
	switch statusText {
	case "OK":
		return prompt.Colorize("[green]OK[reset]")
	case "WARN":
		return prompt.Colorize("[yellow]WARN[reset]")
	default:
		return prompt.Colorize(fmt.Sprintf("[red]%s[reset]", statusText))
	}
}
