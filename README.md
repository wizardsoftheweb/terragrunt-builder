# `terragrunt-builder`

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY 4.0](https://img.shields.io/badge/License-CC_BY_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)

## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Overview](#overview)
- [References](#references)
- [TODOs](#todos)
  - [Parser](#parser)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview

I want to see if I can parse TF and build some quick Terragrunt scaffolding.

## References

- [Great Stack Overflow answer](https://stackoverflow.com/a/66620345/2877698)

## TODOs

### Parser

- Right now the parser only handles variables whose type is string. It will coerce anything else into a string unless it's something complicated like a list or an object.
  - <https://github.com/hashicorp/terraform/blob/v1.2.6/internal/configs/named_values.go>
  - Probably need to understand [`go-cty`](https://github.com/zclconf/go-cty)
