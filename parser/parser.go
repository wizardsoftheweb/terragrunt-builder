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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/hashicorp/hcl/v2"
)

const (
	// DiagIgnoreUnsupportedBlock collects the error message given when an unknown block type is scanned
	DiagIgnoreUnsupportedBlock = "unsupported block"
	// DiagIgnoreUnsupportedAttribute collects the error message given when an unknown attribute is scanned
	DiagIgnoreUnsupportedAttribute = "unsupported attribute"
	// DiagIgnoreUnsupportedArgument collects the error message given when an unknown argument is scanned
	DiagIgnoreUnsupportedArgument = "unsupported argument"
)

var (
	// importantBlocksSchema sets up the blocks we're interested in as we parse TF
	importantBlocksSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}
	// variableBlockSchema grabs only the attributes we're interested in from the variable block
	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "type",
			},
			{
				Name: "default",
			},
		},
	}
	outputBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "type",
			},
			{
				Name: "value",
			},
		},
	}
)

// Variable holds values that may be used for Terragrunt inputs
type Variable struct {
	Name    string
	Default string
}

// Output holds values that may be used for Terragrunt dependencies
type Output struct {
	Name  string
	Value string
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
			ignored := false
			for _, allowedError := range allowedErrors {
				if strings.Contains(strings.ToLower(diag.Error()), strings.ToLower(allowedError)) {
					ignored = true
					break
				}
			}
			if !ignored {
				diagErrors = append(diagErrors, diag)
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

// processSchema is a helper function to process the raw HCL format into something that can be walked and parsed
func processSchema(rawHcl *hcl.File, schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Diagnostics) {
	blocks, diags := rawHcl.Body.Content(schema)
	diagErrs := checkDiagnostics(diags, []string{DiagIgnoreUnsupportedBlock})
	if nil != diagErrs {
		return nil, diagErrs
	}
	return blocks, nil
}

// processVariable turns a variable block into a variable struct
func processVariable(block *hcl.Block) (variable *Variable, diagErr hcl.Diagnostics) {
	if "variable" != block.Type {
		return nil, nil
	}
	blockContent, diags := block.Body.Content(variableBlockSchema)
	diagErr = checkDiagnostics(diags, []string{DiagIgnoreUnsupportedAttribute, DiagIgnoreUnsupportedArgument})
	if nil != diagErr {
		return nil, diagErr
	}
	variable = &Variable{
		Name: block.Labels[0],
	}
	if defaultAttr, ok := blockContent.Attributes["default"]; ok {
		attributeDiags := gohcl.DecodeExpression(defaultAttr.Expr, nil, &variable.Default)
		diagErr = checkDiagnostics(attributeDiags, nil)
		if nil != attributeDiags {
			return nil, diagErr
		}
	}
	return variable, nil
}

// processOutput turns an output block into an output struct
func processOutput(block *hcl.Block) (output *Output, diagErr hcl.Diagnostics) {
	if "output" != block.Type {
		return nil, nil
	}
	blockContent, diags := block.Body.Content(outputBlockSchema)
	diagErr = checkDiagnostics(diags, []string{DiagIgnoreUnsupportedAttribute, DiagIgnoreUnsupportedArgument})
	if nil != diagErr {
		return nil, diagErr
	}
	output = &Output{
		Name: block.Labels[0],
	}
	if valueAttr, ok := blockContent.Attributes["value"]; ok {
		attributeDiags := gohcl.DecodeExpression(valueAttr.Expr, nil, &output.Value)
		diagErr = checkDiagnostics(attributeDiags, nil)
		if nil != attributeDiags {
			return nil, diagErr
		}
	}
	return output, nil
}

func processTerraform(body *hcl.BodyContent) (terraform Terraform, diagErrs hcl.Diagnostics) {
	for _, block := range body.Blocks {
		switch block.Type {
		case "variable":
			variable, diagErr := processVariable(block)
			if nil != diagErr {
				diagErrs = append(diagErrs, diagErr...)
				continue
			}
			terraform.Variables = append(terraform.Variables, variable)
		case "output":
			output, diagErr := processOutput(block)
			if nil != diagErr {
				diagErrs = append(diagErrs, diagErr...)
				continue
			}
			terraform.Outputs = append(terraform.Outputs, output)
		}
	}
	return terraform, diagErrs
}

func processFile(filePath string) (Terraform, error) {
	terraform := Terraform{}
	rawHcl, err := loadFile(filePath)
	if nil != err {
		return terraform, err
	}
	body, diagErrs := processSchema(rawHcl, importantBlocksSchema)
	if nil != diagErrs {
		return terraform, diagErrs
	}
	terraform, diagErrs = processTerraform(body)
	if diagErrs.HasErrors() {
		return Terraform{}, diagErrs
	}
	return terraform, nil
}

func Parse(filePath string) (Terraform, error) {
	fileInfo, statErr := os.Stat(filePath)
	if nil != statErr {
		return Terraform{}, statErr
	}
	if fileInfo.IsDir() {
		// We know we're dealing with a directory, so we'll just iterate over the files in it
		files, _ := ioutil.ReadDir(filePath)
		terraform := Terraform{}
		noTerraform := true
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tf") {
				noTerraform = false
				childPath := path.Join(filePath, file.Name())
				childTerraform, childProcessErr := processFile(childPath)
				if nil != childProcessErr {
					return Terraform{}, childProcessErr
				}
				terraform.Variables = append(terraform.Variables, childTerraform.Variables...)
				terraform.Outputs = append(terraform.Outputs, childTerraform.Outputs...)
			}
		}
		if noTerraform {
			return Terraform{}, fmt.Errorf("no Terraform files found in directory %s", filePath)
		}
		return terraform, nil
	}
	return processFile(filePath)
}
