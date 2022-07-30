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
	"testing"

	"github.com/hashicorp/hcl/v2"

	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
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
	suite.Equal(diags, parsedDiags)
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
