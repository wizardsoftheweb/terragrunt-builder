# Spiking the Parser

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC_BY_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)

## Overview

I wanted to quickly understand how to parse Terraform using [the HCL lib](https://github.com/hashicorp/hcl). There isn't an official example of this, as far as I know. After some quick searching I found a great SO answer to get started. The code in this directory is that setup with some modifications to achieve my goals.

* Parse anything that might show up in TF without dying on things we haven't defined
* Keep every `variable`, ignoring parts of the `variable` we don't care about

## References

* <https://github.com/hashicorp/hcl/issues/363>
* <https://stackoverflow.com/a/66620345/2877698>
