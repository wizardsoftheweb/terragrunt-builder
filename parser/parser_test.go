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
	"path"
	"testing"

	"github.com/hashicorp/hcl/v2"

	"github.com/stretchr/testify/suite"
)

const (
	// fixtureDirectory is the directory containing the fixtures
	fixtureDirectory = "test_fixtures"
	// fixtureDirectoryTerraform is the directory with only Terraform files
	fixtureDirectoryTerraform = "terraform"
	// fixtureFileHclWontParse is a file that will not parse because of a syntax error
	fixtureFileHclWontParse = "hcl_wont_parse.hcl"
	// fixtureFileDoesntExist is a file that does not exist (do not create it!)
	fixtureFileDoesntExist = "doesnt_exist.hcl"
	// fixtureFileParseableHcl is a file that will parse because it is valid HCL
	fixtureFileParseableHcl = "parseable_hcl.hcl"
	// fixtureFileBadTypeDefault is a file containing a variable whose default does not match its type
	fixtureFileBadTypes = "bad_types.tf"
	// fixtureFileTerraformOnlyVariables is a file containing only variable declarations
	fixtureFileTerraformOnlyVariables = "only_variables.tf"
	// fixtureFileTerraformOnlyOutputs is a file containing only variable declarations
	fixtureFileTerraformOnlyOutputs = "only_outputs.tf"
	// fixtureFileTerraformCombined has both variables and outputs
	fixtureFileTerraformCombined = "combined.tf"
)

type ParserTestSuite struct {
	suite.Suite
	fixtureDirectory          string
	terraformFixtureDirectory string
}

func (suite *ParserTestSuite) SetupSuite() {
	suite.fixtureDirectory = path.Join(".", fixtureDirectory)
	suite.terraformFixtureDirectory = path.Join(".", fixtureDirectory, fixtureDirectoryTerraform)
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) Test_checkDiagnostics_NoAllowedErrors() {
	diags := hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Not allowed",
			Detail:   "This is not allowed",
		},
	}
	parsedDiags := checkDiagnostics(diags, nil)
	suite.Equalf(diags, parsedDiags, "Diagnostics should be %v", diags)
}

func (suite *ParserTestSuite) Test_checkDiagnostics_AllAllowedErrors() {
	diags := hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Allowed diagnostic",
			Detail:   "This is allowed",
		},
	}
	allowedErrors := []string{
		"Allowed diagnostic",
	}
	parsedDiags := checkDiagnostics(diags, allowedErrors)
	suite.Nilf(parsedDiags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_checkDiagnostics_MixOfErrors() {
	diags := hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Not allowed",
			Detail:   "This is not allowed",
		},
		{
			Severity: hcl.DiagError,
			Summary:  "Allowed diagnostic",
			Detail:   "This is allowed",
		},
	}
	allowedErrors := []string{
		"Allowed diagnostic",
	}
	expectedDiags := hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Not allowed",
			Detail:   "This is not allowed",
		},
	}
	parsedDiags := checkDiagnostics(diags, allowedErrors)
	suite.Equalf(expectedDiags, parsedDiags, "Diagnostics should be %v", expectedDiags)
}

func (suite *ParserTestSuite) Test_checkDiagnostics_MultipleAllowedErrors() {
	diags := hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Allowed diagnostic one",
			Detail:   "This is allowed",
		},
		{
			Severity: hcl.DiagError,
			Summary:  "Allowed diagnostic two",
			Detail:   "This is allowed",
		},
	}
	allowedErrors := []string{
		"Allowed diagnostic one",
		"Allowed diagnostic two",
	}
	parsedDiags := checkDiagnostics(diags, allowedErrors)
	suite.Nilf(parsedDiags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_loadFile_WontParse() {
	filePath := path.Join(suite.fixtureDirectory, fixtureFileHclWontParse)
	rawHcl, parseErr := loadFile(filePath)
	suite.Nilf(rawHcl, "Raw HCL should be nil")
	suite.NotNilf(parseErr, "Parse error should not be nil")
}

func (suite *ParserTestSuite) Test_loadFile_DoesntExist() {
	filePath := path.Join(suite.fixtureDirectory, fixtureFileDoesntExist)
	rawHcl, parseErr := loadFile(filePath)
	suite.Nilf(rawHcl, "Raw HCL should be nil")
	suite.NotNilf(parseErr, "Parse error should not be nil")
}

func (suite *ParserTestSuite) Test_loadFile_WillParse() {
	filePath := path.Join(suite.fixtureDirectory, fixtureFileParseableHcl)
	rawHcl, parseErr := loadFile(filePath)
	suite.NotNilf(rawHcl, "Raw HCL should not be nil")
	suite.Nilf(parseErr, "Parse error should be nil")
}

func (suite *ParserTestSuite) Test_processSchema_SchemaWithErrors() {
	rawHcl, _ := loadFile(path.Join(suite.fixtureDirectory, fixtureFileParseableHcl))
	body, diags := processSchema(rawHcl, &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	})
	suite.Nilf(body, "Body should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processSchema_SchemaWithoutErrors() {
	rawHcl, _ := loadFile(path.Join(suite.fixtureDirectory, fixtureFileParseableHcl))
	body, diags := processSchema(rawHcl, &hcl.BodySchema{})
	suite.NotNilf(body, "Body should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processVariables_OnlyVariables() {
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyVariables))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	variable, diags := processVariable(body.Blocks[0])
	suite.NotNilf(variable, "Variable should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processVariables_VariableSchemaFails() {
	oldVariableBlockSchema := variableBlockSchema
	defer (func() { variableBlockSchema = oldVariableBlockSchema })()
	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	}
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyVariables))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	variable, diags := processVariable(body.Blocks[0])
	suite.Nilf(variable, "Variable should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processVariables_NotAVariable() {
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyOutputs))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	variable, diags := processVariable(body.Blocks[0])
	suite.Nilf(variable, "Variable should be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processVariables_BadType() {
	rawHcl, _ := loadFile(path.Join(suite.fixtureDirectory, fixtureFileBadTypes))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	variable, diags := processVariable(body.Blocks[0])
	suite.Nilf(variable, "Variable should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processOutputs_OnlyOutputs() {
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyOutputs))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	output, diags := processOutput(body.Blocks[0])
	suite.NotNilf(output, "Output should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processOutputs_OutputSchemaFails() {
	oldOutputBlockSchema := outputBlockSchema
	defer (func() { outputBlockSchema = oldOutputBlockSchema })()
	outputBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	}
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyOutputs))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	output, diags := processOutput(body.Blocks[0])
	suite.Nilf(output, "Output should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_procesOutputs_NotAnOutput() {
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformOnlyVariables))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	output, diags := processOutput(body.Blocks[0])
	suite.Nilf(output, "Output should be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processOutputs_BadType() {
	rawHcl, _ := loadFile(path.Join(suite.fixtureDirectory, fixtureFileBadTypes))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	output, diags := processOutput(body.Blocks[1])
	suite.Nilf(output, "Output should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processTerraform_BadTypes() {
	rawHcl, _ := loadFile(path.Join(suite.fixtureDirectory, fixtureFileBadTypes))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	terraform, diags := processTerraform(body)
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processTerraform_Success() {
	rawHcl, _ := loadFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformCombined))
	body, _ := processSchema(rawHcl, importantBlocksSchema)
	terraform, diags := processTerraform(body)
	suite.NotNilf(terraform.Variables, "Terraform variables should not be nil")
	suite.NotNilf(terraform.Outputs, "Terraform outputs should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processFile_Success() {
	terraform, diags := processFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformCombined))
	suite.NotNilf(terraform.Variables, "Terraform variables should not be nil")
	suite.NotNilf(terraform.Outputs, "Terraform outputs should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_processFile_IsNotHcl() {
	terraform, diags := processFile("parser_test.go")
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processFile_FileDoesNotExist() {
	terraform, diags := processFile(path.Join(suite.fixtureDirectory, fixtureFileDoesntExist))
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processFile_SchemaFails() {
	oldSchema := importantBlocksSchema
	defer (func() { importantBlocksSchema = oldSchema })()
	importantBlocksSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	}
	terraform, diags := processFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformCombined))
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_processFile_TerraformFails() {
	oldSchema := variableBlockSchema
	defer (func() { variableBlockSchema = oldSchema })()
	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	}
	terraform, diags := processFile(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformCombined))
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_Parse_DoesNotExist() {
	_, diags := Parse(path.Join(suite.fixtureDirectory, fixtureFileDoesntExist))
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_Parse_SingleFile() {
	terraform, diags := Parse(path.Join(suite.terraformFixtureDirectory, fixtureFileTerraformCombined))
	suite.NotNilf(terraform.Variables, "Terraform variables should not be nil")
	suite.NotNilf(terraform.Outputs, "Terraform outputs should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_Parse_DirectoryWithoutTerraform() {
	terraform, diags := Parse(".")
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}

func (suite *ParserTestSuite) Test_Parse_DirectorySuccess() {
	terraform, diags := Parse(suite.terraformFixtureDirectory)
	suite.NotNilf(terraform.Variables, "Terraform variables should not be nil")
	suite.Greaterf(len(terraform.Variables), 1, "There should be several variables")
	suite.NotNilf(terraform.Outputs, "Terraform outputs should not be nil")
	suite.Greaterf(len(terraform.Outputs), 1, "There should be several outputs")
	suite.Nilf(diags, "Diagnostics should be nil")
}

func (suite *ParserTestSuite) Test_Parse_DirectoryWithDiagErrors() {
	oldSchema := variableBlockSchema
	defer (func() { variableBlockSchema = oldSchema })()
	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "missing",
				Required: true,
			},
		},
	}
	terraform, diags := Parse(suite.terraformFixtureDirectory)
	suite.Nilf(terraform.Variables, "Terraform variables should be nil")
	suite.Nilf(terraform.Outputs, "Terraform outputs should be nil")
	suite.NotNilf(diags, "Diagnostics should not be nil")
}
