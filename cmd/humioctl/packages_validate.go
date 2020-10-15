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
	"path/filepath"

	"github.com/humio/cli/api"

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func validatePackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "validate [flags] <repo-or-view-name> <package-dir>",
		Short: "Validate a package's content.",
		Args:  cobra.ExactArgs(2),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			repoOrViewName := args[0]
			dirPath := args[1]

			if !filepath.IsAbs(dirPath) {
				var err error
				dirPath, err = filepath.Abs(dirPath)
				if err != nil {
					return nil, fmt.Errorf("invalid path: %w", err)
				}
				dirPath += "/"
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			validationResult, apiErr := client.Packages().Validate(repoOrViewName, dirPath)
			if apiErr != nil {
				return nil, fmt.Errorf("errors validating package: %w", apiErr)
			}

			if validationResult.IsValid() {
				return fmt.Sprintf("Package in %s is valid.", dirPath), nil
			} else {
				return nil, fmt.Errorf("%s", FormatResult(validationResult, true))
			}
		}),
	}

	return &cmd
}

func printValidation(out *prompt.Prompt, validationResult *api.ValidationResponse) {
	out.Error("Package is not valid")
	out.Error(out.List(validationResult.InstallationErrors))
	out.Error(out.List(validationResult.ParseErrors))
}
