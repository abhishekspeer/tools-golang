// Package rdfParser2_1 contains functions to read, load and parse
// SPDX RDF files.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfParser2_1

import (
	"fmt"
	
	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/deltamobile/goraptor"
	"io"
	"strings"

)

// Constants representing RDF formats supported by raptor.
const (
	Fmt_ntriples     = "ntriples"      // for N-Triples
	Fmt_turtle       = "turtle"        // for Turtle Terse RDF Triple Language
	Fmt_rdfxmlXmp    = "rdfxml-xmp"    // for RDF/XML (XMP Profile)
	Fmt_rdfxmlAbbrev = "rdfxml-abbrev" // for RDF/XML (Abbreviated)
	Fmt_rdfxml       = "rdfxml"        // for RDF/XML
	Fmt_rss          = "rss-1.0"       // for RSS 1.0
	Fmt_atom         = "atom"          // for Atom 1.0
	Fmt_dot          = "dot"           // for GraphViz DOT format
	Fmt_jsonTriples  = "json-triples"  // for RDF/JSON Triples
	Fmt_json         = "json"          // for RDF/JSON Resource-Centric
	Fmt_html         = "html"          // for HTML Table
	Fmt_nquads       = "nquads"        // for N-Quads
)

// Parser error messages.
const {
	msgUnknownType          = "Found type %s which is unknown."
	msgIncompatibleTypes    = "%s is already set to be type %s and cannot be changed to type %s."
	msgAlreadyDefined       = "Property already defined."
	msgPropertyNotSupported = "Property %s is not supported for %s."

}

// Checks if `fmt` is one of the raptor supported formats (Fmt* constantsa above). 
// The special "rdf" value is considered invalid by this function.
func FormatOk(fmt string) bool {
	fmts := []string{
		Fmt_ntriples,
		Fmt_turtle,
		Fmt_rdfxmlXmp,
		Fmt_rdfxmlAbbrev,
		Fmt_rdfxml,
		Fmt_rss,
		Fmt_atom,
		Fmt_dot,
		Fmt_jsonTriples,
		Fmt_json,
		Fmt_html,
		Fmt_nquads,
	}
	for _, f := range fmts {
		if fmt == f {
			return true
		}
	}
	return false
}

// Calls interface to parse a document.
func Parse(input io.Reader, format string) (*spdx.Document2_1, error) {
	parser := rdfParser2_1(input, format)
	defer parser.Free()
	return parser.Parse()
}

// RDF Parser. 
// Use a RDF Parser to parse SPDX RDF files to SPDX documents.
// Use `rdfParser2_1()` method to create a new parser.
type Parser struct {
	rdfparser *goraptor.Parser
	input     io.Reader
	index     map[string]*builder
	buffer    map[string][]bufferEntry
	doc       *spdx.Document2_1
}

// This creates a goraptor.Parser object that needs to be freed after use.
// Call Parser.Free() after using the Parser.
func rdfParser2_1(input io.Reader, format string) *Parser {
	if format == "rdf" {
		format = "guess"
	}

	return &Parser{
		rdfparser: goraptor.rdfParser2_1(format),
		input:     input,
		index:     make(map[string]*builder),
		buffer:    make(map[string][]bufferEntry),
	}
}

// Free the goraptor parser.
func (p *Parser) Free() {
	p.rdfparser.Free()
	p.doc = nil
}

