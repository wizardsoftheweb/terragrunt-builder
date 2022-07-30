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

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var (
	configFileSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}

	variableBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: "description",
			},
			{
				Name: "type",
			},
			{
				Name: "sensitive",
			},
		},
	}
)

type Config struct {
	Variables []*Variable
}

type Variable struct {
	Name        string
	Description string
	Type        string
	Sensitive   bool
}

func main() {
	config := configFromFile("test.tf")
	for _, v := range config.Variables {
		fmt.Printf("%+v\n", v)
	}
}

func configFromFile(filePath string) *Config {
	content, err := os.ReadFile(filePath) // go 1.16
	if err != nil {
		log.Fatal(err)
	}

	file, diags := hclsyntax.ParseConfig(content, filePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		log.Fatal("ParseConfig", diags)
	}

	bodyCont, diags := file.Body.Content(configFileSchema)
	if diags.HasErrors() {
		for _, diagErr := range diags.Errs() {
			if !strings.Contains(strings.ToLower(diagErr.Error()), "unsupported block type") {
				log.Printf("%+v", diagErr)
			}
		}
	}

	res := &Config{}

	for _, block := range bodyCont.Blocks {
		v := &Variable{
			Name: block.Labels[0],
		}

		blockCont, diags := block.Body.Content(variableBlockSchema)
		if diags.HasErrors() {
			log.Fatal("block content", diags)
		}

		if attr, exists := blockCont.Attributes["description"]; exists {
			diags := gohcl.DecodeExpression(attr.Expr, nil, &v.Description)
			if diags.HasErrors() {
				log.Fatal("description attr", diags)
			}
		}

		if attr, exists := blockCont.Attributes["sensitive"]; exists {
			diags := gohcl.DecodeExpression(attr.Expr, nil, &v.Sensitive)
			if diags.HasErrors() {
				log.Fatal("sensitive attr", diags)
			}
		}

		if attr, exists := blockCont.Attributes["type"]; exists {
			v.Type = hcl.ExprAsKeyword(attr.Expr)
			if v.Type == "" {
				log.Fatal("type attr", "invalid value")
			}
		}

		res.Variables = append(res.Variables, v)
	}
	return res
}
