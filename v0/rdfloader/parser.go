// Package rdfparser2_1 contains functions to read, load and parse
// SPDX RDF files.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfparser2_1

import (
	"strings"

	"github.com/deltamobile/goraptor"
)

// Constants representing RDF formats supported by raptor.
// Currently working on RDF support only.
const (
	// Fmt_ntriples     = "ntriples"      // for N-Triples
	// Fmt_turtle       = "turtle"        // for Turtle Terse RDF Triple Language
	Fmt_rdfxmlXmp    = "rdfxml-xmp"    // for RDF/XML (XMP Profile)
	Fmt_rdfxmlAbbrev = "rdfxml-abbrev" // for RDF/XML (Abbreviated)
	Fmt_rdfxml       = "rdfxml"        // for RDF/XML
	// Fmt_rss          = "rss-1.0"       // for RSS 1.0
	// Fmt_atom         = "atom"          // for Atom 1.0
	// Fmt_dot          = "dot"           // for GraphViz DOT format
	// Fmt_jsonTriples  = "json-triples"  // for RDF/JSON Triples
	// Fmt_json         = "json"          // for RDF/JSON Resource-Centric
	// Fmt_html         = "html"          // for HTML Table
	// Fmt_nquads       = "nquads"        // for N-Quads
)

// Useful RDF URIs
const (
	baseUri    = "http://spdx.org/rdf/terms#"
	licenseUri = "http://spdx.org/licenses/"
)

// Common RDF prefixes used in SPDX RDF Representations.
var rdfPrefixes = map[string]string{
	"":      baseUri,
	"rdfs:": "http://www.w3.org/2000/01/rdf-schema#",
	"ns:":   "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"doap:": "http://usefulinc.com/ns/doap#",
}

// Common RDF parser error messages.
const (
	msgIncompatibleTypes    = "%s is already set to be type %s and cannot be changed to type %s."
	msgPropertyNotSupported = "Property %s is not supported for %s."
	msgAlreadyDefined       = "Property already defined."
	msgUnknownType          = "Found type %s which is unknown."
)

// Expands the prefixes to their full URIs.
// If there is no ":" or there is another prefix, it expands to baseUri.
func prefix(k string) *goraptor.Uri {
	var pref string
	rest := k
	if i := strings.Index(k, ":"); i >= 0 {
		pref = k[:i+1]
		rest = k[i+1:]
	}
	if long, ok := rdfPrefixes[pref]; ok {
		pref = long
	}
	uri := goraptor.Uri(pref + rest)
	return &uri
}

// Change the RDF prefixes to their short forms.
func shortPrefix(t goraptor.Term) string {
	str := termStr(t)
	for short, long := range rdfPrefixes {
		if strings.HasPrefix(str, long) {
			return strings.Replace(str, long, short, 1)
		}
	}
	return str
}

// goraptor.Term to string. Returns empty string if the term given is not one of
// the following types: *goraptor.Uri, *goraptor.Blank or *goraptor.Literal.
func termStr(term goraptor.Term) string {
	switch t := term.(type) {
	case *goraptor.Uri:
		return string(*t)
	case *goraptor.Blank:
		return string(*t)
	case *goraptor.Literal:
		return t.Value
	default:
		return ""
	}
}

// Create *goraptor.Blank from string
func blank(b string) *goraptor.Blank {
	return (*goraptor.Blank)(&b)
}

// Create *goraptor.Uri from string
func uri(uri string) *goraptor.Uri {
	return (*goraptor.Uri)(&uri)
}

// Create *goraptor.Literal from string
func literal(lit string) *goraptor.Literal {
	return &goraptor.Literal{Value: lit}
}

// Checks if `fmt` is one of the raptor supported formats (Fmt* constantsa above).
// The special "rdf" value is considered invalid by this function.
// Parser error messages.

func FormatOk(fmt string) bool {
	fmts := []string{
		Fmt_rdfxmlXmp,
		Fmt_rdfxmlAbbrev,
		Fmt_rdfxml,
	}

	for _, f := range fmts {
		if fmt == f {
			return true
		}

	}
	return false
}
