[![Build Status](https://travis-ci.org/spdx/tools-golang.svg?branch=master)](https://travis-ci.org/spdx/tools-golang)
[![Coverage Status](https://coveralls.io/repos/github/spdx/tools-golang/badge.svg)](https://coveralls.io/github/spdx/tools-golang)

# tools-golang

tools-golang is a collection of Go packages intended to make it easier for
Go programs to work with [SPDX®](https://spdx.org/) files.

This software is in an early state, and its API may change significantly
(hence the "v0/" directory).

## What it does

tools-golang currently works with files conformant to version 2.1 of the
SPDX specification, available at: https://spdx.org/specifications

tools-golang provides the following packages:

* *v0/spdx* - in-memory data model for the sections of an SPDX document
* *v0/tvloader* - tag-value file loader
* *v0/tvsaver* - tag-value file saver
* *v0/rdfloader* - RDF file loader
* *v0/rdfsaver* - RDF file saver
* *v0/builder* - builds "empty" SPDX document (with hashes) for directory contents
* *v0/idsearcher* - searches for [SPDX short-form IDs](https://spdx.org/ids/) and builds SPDX document
* *v0/licensediff* - compares concluded licenses between files in two packages
* *v0/reporter* - generates basic license count report from SPDX document
* *v0/utils* - various utility functions that support the other tools-golang packages

Examples for how to use these packages can be found in the `examples/`
directory.

## What it doesn't do

tools-golang doesn't currently do any of the following:

* work with files under any version of the SPDX spec *other than* v2.1
* convert between RDF and tag-value files, or between different versions
* enable applications to interact with SPDX files without needing to care
  (too much) about the particular SPDX file version

We are working towards adding functionality for all of these. Code contributions
are welcome!

## Requirements

* [libraptor2](http://librdf.org/raptor/) for parsing and serializing RDF
* [goraptor](https://github.com/deltamobile/goraptor) fork of [goraptor](https://bitbucket.org/ww/goraptor/src) by William Waites for RDF

## Licenses

As indicated in `LICENSE-code.txt`, tools-golang **source code files** are
provided and may be used, at your option, under *either*:
* Apache License, version 2.0 (**Apache-2.0**), **OR**
* GNU General Public License, version 2.0 or later (**GPL-2.0-or-later**).

As indicated in `LICENSE-docs.txt`, tools-golang **documentation files** are
provided and may be used under the Creative Commons Attribution
4.0 International license (**CC-BY-4.0**).

This `README.md` file is documentation:

`SPDX-License-Identifier: CC-BY-4.0`
