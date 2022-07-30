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
	// fixtureFileTerraformOnlyVariables is a file containing only variable declarations
	fixtureFileTerraformOnlyVariables = "only_variables.tf"
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
	fmt.Printf("%+v\n", variable)
	fmt.Printf("%+v\n", diags)
	suite.NotNilf(variable, "Variable should not be nil")
	suite.Nilf(diags, "Diagnostics should be nil")
}
