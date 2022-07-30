// Copyright 2022 CJ Harries
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

package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/hashicorp/hcl/v2"
)

const (
	// DiagIgnoreUnsupportedBlock collects the error message given when an unknown block type is scanned
	DiagIgnoreUnsupportedBlock = "unsupported block"
	// DiagIgnoreUnsupportedAttribute collects the error message given when an unknown attribute is scanned
	DiagIgnoreUnsupportedAttribute = "unsupported attribute"
)

//var (
//	// importantBlocksSchema sets up the blocks we're interested in as we parse TF
//	importantBlocksSchema = &hcl.BodySchema{
//		Blocks: []hcl.BlockHeaderSchema{
//			{
//				Type:       "variable",
//				LabelNames: []string{"name"},
//			},
//			{
//				Type:       "output",
//				LabelNames: []string{"name"},
//			},
//		},
//	}
//	// variableBlockSchema grabs only the attributes we're interested in from the variable block
//	variableBlockSchema = &hcl.BodySchema{
//		Attributes: []hcl.AttributeSchema{
//			{
//				Name: "default",
//			},
//		},
//	}
//)

// Variable holds values that may be used for Terragrunt inputs
type Variable struct {
	Name    string
	Default interface{}
}

// Output holds values that may be used for Terragrunt dependencies
type Output struct {
	Name  string
	Value interface{}
}

// Terraform holds the blocks from TF files we're interested in working with
type Terraform struct {
	Variables []*Variable
	Outputs   []*Output
}

// checkDiagnostics is a simple helper function to ignore diagnostic errors we may not care about. For example, if we're
// parsing for variables, we may only pass in a schema that contains variables and their structure. Things like
// resources and outputs would trigger a diagnostic error.
func checkDiagnostics(diags hcl.Diagnostics, allowedErrors []string) (diagErrors hcl.Diagnostics) {
	if 0 == len(allowedErrors) {
		return diags
	}
	if diags.HasErrors() {
		for _, diag := range diags {
			for _, allowedError := range allowedErrors {
				if !strings.Contains(strings.ToLower(diag.Error()), strings.ToLower(allowedError)) {
					fmt.Println(diag.Error())
					diagErrors = append(diagErrors, diag)
				}
			}
		}
	}
	return diagErrors
}

// loadFile reads the file and parses it into a raw HCL format, ready for unmarshalling
func loadFile(filePath string) (rawHcl *hcl.File, err error) {
	fileContents, fileReadErr := os.ReadFile(filePath)
	if fileReadErr != nil {
		return nil, fileReadErr
	}
	rawHcl, hclParseDiags := hclsyntax.ParseConfig(fileContents, filePath, hcl.Pos{Line: 1, Column: 1})
	if hclParseDiags.HasErrors() {
		return nil, hclParseDiags
	}
	return rawHcl, nil
}
