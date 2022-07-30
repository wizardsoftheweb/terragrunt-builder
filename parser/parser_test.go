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
	fixtureFileHclWontParse = "hcl_wont_parse.hcl"
)

type ParserTestSuite struct {
	suite.Suite
	fixtureDirectory string
}

func (suite *ParserTestSuite) SetUpSuite() {
	suite.fixtureDirectory = path.Join(".", "test_fixtures")
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

func (suite *ParserTestSuite) Test_loadFile_WontParse() {
	filePath := path.Join(suite.fixtureDirectory, fixtureFileHclWontParse)
	rawHcl, parseErr := loadFile(filePath)
	suite.Nilf(rawHcl, "Raw HCL should be nil")
	suite.NotNilf(parseErr, "Parse error should not be nil")
}

func (suite *ParserTestSuite) Test_loadFile_DoesntExist() {
	filePath := path.Join(suite.fixtureDirectory, "doesnt_exist.hcl")
	rawHcl, parseErr := loadFile(filePath)
	suite.Nilf(rawHcl, "Raw HCL should be nil")
	suite.NotNilf(parseErr, "Parse error should not be nil")
}
