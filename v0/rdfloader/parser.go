// Package rdfparser2_1 contains functions to read, load and parse
// SPDX RDF files.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdfparser2_1

import (
	"fmt"
	"io"
	"strings"

	"github.com/deltamobile/goraptor"
	"github.com/spdx/tools-golang/v0/spdx"
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

// Useful helper pair struct. Might use later.

// type pair struct {
// 	key,val string
// }

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
)

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


// ----------------------------------
// TODO:
// RDF element types in URI format. (RDF classes). 
// ADD RDF element types in URI format
// -----------------------------------
// RDF element types in URI format. (RDF classes). 
// Includes deprecated as well for now. 
var (
	uri_nstype = uri("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")

	typeDocument           = prefix("SpdxDocument")
	typeCreationInfo       = prefix("CreationInfo")
	typePackage            = prefix("Package")
	typeFile               = prefix("File")
	typeVerificationCode   = prefix("PackageVerificationCode")
	typeChecksum           = prefix("Checksum")
	typeArtifactOf         = prefix("doap:Project")
	typeReview             = prefix("Review")
	typeExtractedLicence   = prefix("ExtractedLicensingInfo")
	typeAnyLicence         = prefix("AnyLicenseInfo")
	typeConjunctiveSet     = prefix("ConjunctiveLicenseSet")
	typeDisjunctiveSet     = prefix("DisjunctiveLicenseSet")
	typeLicence            = prefix("License")
	typeAbstractLicenceSet = blank("abstractLicenceSet")
)



// Calls interface to parse a document.
func Parse(input io.Reader, format string) (*spdx.Document2_1, error) {
	parser := NewParser(input, format)
	defer parser.Free()
	return parser.Parse()
}
// type of element t represents the spdx element that ptr builds
type builder struct {
	t        goraptor.Term 
	ptr      interface{}   
	updaters map[string]updater
}


func (b *builder) apply(pred, obj goraptor.Term, meta *spdx.Meta) error {
	property := shortPrefix(pred)
	f, ok := b.updaters[property]
	if !ok {
		return spdx.NewParseError(fmt.Sprintf(msgPropertyNotSupported, property, b.t), meta)
	}
	return f(obj, meta)
}

func (b *builder) has(pred string) bool {
	_, ok := b.updaters[pred]
	return ok
}

type updater func(goraptor.Term, *spdx.Meta) error

type bufferEntry struct {
	*goraptor.Statement
	*spdx.Meta
}


// RDF Parser.
// Use a RDF Parser to parse SPDX RDF files to SPDX documents.
// Use `NewParser()` method to create a new parser.
type Parser struct {
	rdfparser *goraptor.Parser
	input     io.Reader
	index     map[string]*builder
	buffer    map[string][]bufferEntry
	doc       *spdx.Document2_1
}

// This creates a goraptor.Parser object that needs to be freed after use.
// Call Parser.Free() after using the Parser.
func NewParser(input io.Reader, format string) *Parser {
	if format == "rdf" {
		format = "guess"
	}

	return &Parser{
		rdfparser: goraptor.NewP(format),
		input:     input,
		index:     make(map[string]*builder),
		buffer:    make(map[string][]bufferEntry),
	}
}

// Parse the whole input stream and return the resulting spdx.Document2_1 or the first error that occurred.
func (p *Parser) Parse() (*spdx.Document2_1, error) {
	ch := p.rdfparser.Parse(p.input, baseUri)
	locCh := p.rdfparser.LocatorChan()
	var err error
	for statement := range ch {
		locator := <-locCh
		meta := spdx.NewMetaL(locator.Line)
		if err = p.processTruple(statement, meta); err != nil {
			break
		}
	}
	// Consume input channel in case of error. Otherwise goraptor will keep the goroutine busy.
	for _ = range ch {
		<-locCh
	}
	return p.doc, err
}
// Free the goraptor parser after use.
func (p *Parser) Free() {
	p.rdfparser.Free()
	p.doc = nil
}

// Set the type of node to t.
// If the node does not exist, a builder of the required type is created and the buffered
// statements will be applied in fifo order.
// If the node exists and the types are not compatible, a ParseError is returned.
func (p *Parser) setType(node, t goraptor.Term, meta *spdx.Meta) (interface{}, error) {
	nodeStr := termStr(node)
	bldr, ok := p.index[nodeStr]
	if ok {
		if !equalTypes(bldr.t, t) && bldr.has("ns:type") {
			//apply the type change
			if err := bldr.apply(uri("ns:type"), t, meta); err != nil {
				return nil, err
			}
			return bldr.ptr, nil
		}
		if !compatibleTypes(bldr.t, t) {
			return nil, spdx.NewParseError(fmt.Sprintf(msgIncompatibleTypes, node, bldr.t, t), meta)
		}
		return bldr.ptr, nil
	}

// Set the type of node to t.
// If the node does not exist, a builder of the required type is created and the buffered
// statements will be applied in fifo order.
// If the node exists and the types are not compatible, a ParseError is returned.
func (p *Parser) setType(node, t *goraptor.Term, meta *spdx.Meta) (interface{}, error) {
	nodeStr := termStr(node)
	bldr, ok := p.index[nodeStr]
	if ok {
		if !equalTypes(bldr.t, t) && bldr.has("ns:type") {
			//apply the type change
			if err := bldr.apply(uri("ns:type"), t, meta); err != nil {
				return nil, err
			}
			return bldr.ptr, nil
		}
		if !compatibleTypes(bldr.t, t) {
			return nil, spdx.NewParseError(fmt.Sprintf(msgIncompatibleTypes, node, bldr.t, t), meta)
		}
		return bldr.ptr, nil
	}

	// new builder by type
	switch {
	case t.Equals(typeDocument):
		p.doc = &spdx.Document2_1{Meta: meta}
		bldr = p.documentMap(p.doc)
	case t.Equals(typeCreationInfo):
		bldr = p.creationInfoMap(&spdx.CreationInfo{Meta: meta})
	case t.Equals(typePackage):
		bldr = p.packageMap(&spdx.Package{Meta: meta})
	case t.Equals(typeChecksum):
		bldr = p.checksumMap(&spdx.Checksum{Meta: meta})
	case t.Equals(typeVerificationCode):
		bldr = p.verificationCodeMap(&spdx.VerificationCode{Meta: meta})
	case t.Equals(typeFile):
		bldr = p.fileMap(&spdx.File{Meta: meta})
	case t.Equals(typeReview):
		bldr = p.reviewMap(&spdx.Review{Meta: meta})
	case t.Equals(typeArtifactOf):
		artif := &spdx.ArtifactOf{Meta: meta}
		if artifUri, ok := node.(*goraptor.Uri); ok {
			artif.ProjectUri.Val = termStr(artifUri)
			artif.ProjectUri.Meta = meta
		}
		bldr = p.artifactOfMap(artif)
	case t.Equals(typeExtractedLicence):
		bldr = p.extractedLicensingInfoMap(&spdx.ExtractedLicence{Meta: meta})
	case t.Equals(typeAnyLicence):
		switch t := node.(type) {
		case *goraptor.Uri: // licence in spdx licence list
			bldr = p.licenceReferenceBuilder(node, meta)
		case *goraptor.Blank: // licence reference or abstract set
			if strings.HasPrefix(strings.ToLower(termStr(t)), "licenseref") {
				bldr = p.extractedLicensingInfoMap(&spdx.ExtractedLicence{Meta: meta})
			} else {
				bldr = p.licenceSetMap(&spdx.LicenceSet{
					Members: make([]spdx.AnyLicence, 0),
					Meta:    meta,
				})
			}
		}
	case t.Equals(typeLicence):
		bldr = p.licenceReferenceBuilder(node, meta)
	case t.Equals(typeAbstractLicenceSet):
		bldr = p.licenceSetMap(&spdx.LicenceSet{
			Members: make([]spdx.AnyLicence, 0),
			Meta:    meta,
		})
	case t.Equals(typeConjunctiveSet):
		bldr = p.conjunctiveSetBuilder(meta)
	case t.Equals(typeDisjunctiveSet):
		bldr = p.disjuntiveSetBuilder(meta)
	default:
		return nil, spdx.NewParseError(fmt.Sprintf(msgUnknownType, t), meta)
	}

	p.index[nodeStr] = bldr

	// run buffer
	buf := p.buffer[nodeStr]
	for _, stm := range buf {
		if err := bldr.apply(stm.Predicate, stm.Object, stm.Meta); err != nil {
			return nil, err
		}
	}
	delete(p.buffer, nodeStr)

	return bldr.ptr, nil
}



// Process a SPDX Truple.
func (p *Parser) processTruple(stm *goraptor.Statement, meta *spdx.Meta) error {
	node := termStr(stm.Subject)
	if stm.Predicate.Equals(uri_nstype) {
		_, err := p.setType(stm.Subject, stm.Object, meta)
		return err
	}

	// apply function if it's a builder
	bldr, ok := p.index[node]
	if ok {
		return bldr.apply(stm.Predicate, stm.Object, meta)
	}

	// buffer statement
	if _, ok := p.buffer[node]; !ok {
		p.buffer[node] = make([]bufferEntry, 0)
	}
	p.buffer[node] = append(p.buffer[node], bufferEntry{stm, meta})

	return nil
}



// Checks if found is any of the need types. Note: a type term of type
// goraptor.Uri is not the same type as one of type goraptor.Blank; same
// applies for other combinations.
func equalTypes(found goraptor.Term, need ...goraptor.Term) bool {
	for _, b := range need {
		if found == b || found.Equals(b) {
			return true
		}
	}
	return false
}

// Checks if found is the same as need.
//
// If need is any of typeLicence, typeDisjunctiveSet, typeConjunctiveSet
// and typeExtractedLicence and found is AnyLicence, it  is permitted and
// the function returns true.
func compatibleTypes(found, need goraptor.Term) bool {
	if equalTypes(found, need) {
		return true
	}
	if equalTypes(need, typeAnyLicence) {
		return equalTypes(found, typeExtractedLicence, typeConjunctiveSet, typeDisjunctiveSet, typeLicence)
	}
	return false
}

// WIP

// Request that a SPDX Element has a specific type. If it does not have, a
// parse error is returned. If `node` is not parsed yet, it is created and set
// to be of type `t`.
//
// If the node is found and the types match, this method returns a pointer to
// that element, but of type interface{} and a nil error. To get a more specific
// element, use one of the other req* functions (reqDocument, reqFile, etc.).
//
// Parser.req* functions are supposed to get the node from either the index check,
// if it's the required type and return a pointer to the relevant spdx.* object.
func (p *Parser) reqType(node, t goraptor.Term) (interface{}, error) {
	bldr, ok := p.index[termStr(node)]
	if ok {
		if !compatibleTypes(bldr.t, t) {
			return nil, fmt.Errorf(msgIncompatibleTypes, node, bldr.t, t)
		}
		return bldr.ptr, nil
	}
	return p.setType(node, t, nil)
}

func (p *Parser) reqDocument(node goraptor.Term) (*spdx.Document2_1, error) {
	obj, err := p.reqType(node, typeDocument)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.Document2_1), err
}
func (p *Parser) reqCreationInfo(node goraptor.Term) (*spdx.CreationInfo, error) {
	obj, err := p.reqType(node, typeCreationInfo)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.CreationInfo), err
}
func (p *Parser) reqPackage(node goraptor.Term) (*spdx.Package, error) {
	obj, err := p.reqType(node, typePackage)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.Package), err
}
func (p *Parser) reqFile(node goraptor.Term) (*spdx.File, error) {
	obj, err := p.reqType(node, typeFile)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.File), err
}
func (p *Parser) reqVerificationCode(node goraptor.Term) (*spdx.VerificationCode, error) {
	obj, err := p.reqType(node, typeVerificationCode)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.VerificationCode), err
}
func (p *Parser) reqChecksum(node goraptor.Term) (*spdx.Checksum, error) {
	obj, err := p.reqType(node, typeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.Checksum), err
}
func (p *Parser) reqReview(node goraptor.Term) (*spdx.Review, error) {
	obj, err := p.reqType(node, typeReview)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.Review), err
}
func (p *Parser) reqExtractedLicence(node goraptor.Term) (*spdx.ExtractedLicence, error) {
	obj, err := p.reqType(node, typeExtractedLicence)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.ExtractedLicence), err
}
func (p *Parser) reqAnyLicence(node goraptor.Term) (spdx.AnyLicence, error) {
	obj, err := p.reqType(node, typeAnyLicence)
	if err != nil {
		return nil, err
	}
	switch lic := obj.(type) {
	case *spdx.AnyLicence:
		return *lic, nil
	case *spdx.ConjunctiveLicenceSet:
		return *lic, nil
	case *spdx.DisjunctiveLicenceSet:
		return *lic, nil
	case *[]spdx.AnyLicence:
		return nil, nil
	case *spdx.Licence:
		return *lic, nil
	case *spdx.ExtractedLicence:
		return lic, nil
	default:
		return nil, fmt.Errorf("Unexpected error, an element of type AnyLicence cannot be casted to any licence type. %s || %#v", node, obj)
	}
}
func (p *Parser) reqArtifactOf(node goraptor.Term) (*spdx.ArtifactOf, error) {
	obj, err := p.reqType(node, typeArtifactOf)
	if err != nil {
		return nil, err
	}
	return obj.(*spdx.ArtifactOf), err
}
