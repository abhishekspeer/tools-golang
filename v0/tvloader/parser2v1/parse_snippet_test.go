// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v1

import (
	"testing"

	"github.com/swinslow/spdx-go/v0/spdx"
)

// ===== Parser snippet section state change tests =====
func TestParser2_1SnippetStartsNewSnippetAfterParsingSnippetSPDXIDTag(t *testing.T) {
	// create the first snippet
	sid1 := "SPDXRef-s1"

	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "test"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: sid1},
	}
	s1 := parser.snippet
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	// the File's Snippets should have this one only
	if len(parser.file.Snippets) != 1 {
		t.Errorf("Expected len(Snippets) to be 1, got %d", len(parser.file.Snippets))
	}
	if parser.file.Snippets[0] != s1 {
		t.Errorf("Expected snippet %v in Snippets[0], got %v", s1, parser.file.Snippets[0])
	}
	if parser.file.Snippets[0].SnippetSPDXIdentifier != sid1 {
		t.Errorf("expected snippet ID %s in Snippets[0], got %s", sid1, parser.file.Snippets[0].SnippetSPDXIdentifier)
	}

	// now add a new snippet
	sid2 := "SPDXRef-s2"
	err := parser.parsePair2_1("SnippetSPDXID", sid2)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psSnippet2_1 {
		t.Errorf("expected state to be %v, got %v", psSnippet2_1, parser.st)
	}
	// and a snippet should be created
	if parser.snippet == nil {
		t.Fatalf("parser didn't create new snippet")
	}
	// and the snippet ID should be as expected
	if parser.snippet.SnippetSPDXIdentifier != sid2 {
		t.Errorf("expected snippet ID %s, got %s", sid2, parser.snippet.SnippetSPDXIdentifier)
	}
	// and the File's Snippets should be of size 2 and have these two
	if len(parser.file.Snippets) != 2 {
		t.Errorf("Expected len(Snippets) to be 2, got %d", len(parser.file.Snippets))
	}
	if parser.file.Snippets[0] != s1 {
		t.Errorf("Expected snippet %v in Snippets[0], got %v", s1, parser.file.Snippets[0])
	}
	if parser.file.Snippets[0].SnippetSPDXIdentifier != sid1 {
		t.Errorf("expected snippet ID %s in Snippets[0], got %s", sid1, parser.file.Snippets[0].SnippetSPDXIdentifier)
	}
	if parser.file.Snippets[1] != parser.snippet {
		t.Errorf("Expected snippet %v in Snippets[1], got %v", parser.snippet, parser.file.Snippets[1])
	}
	if parser.file.Snippets[1].SnippetSPDXIdentifier != sid2 {
		t.Errorf("expected snippet ID %s in Snippets[1], got %s", sid2, parser.file.Snippets[1].SnippetSPDXIdentifier)
	}
}

func TestParser2_1SnippetStartsNewPackageAfterParsingPackageNameTag(t *testing.T) {
	p1Name := "package1"
	f1Name := "f1.txt"
	s1Name := "SPDXRef-s1"
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: p1Name},
		file:    &spdx.File2_1{FileName: f1Name},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: s1Name},
	}
	p1 := parser.pkg
	f1 := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	// now add a new package
	p2Name := "package2"
	err := parser.parsePair2_1("PackageName", p2Name)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should go back to Package
	if parser.st != psPackage2_1 {
		t.Errorf("expected state to be %v, got %v", psPackage2_1, parser.st)
	}
	// and a package should be created
	if parser.pkg == nil {
		t.Fatalf("parser didn't create new pkg")
	}
	// and the package name should be as expected
	if parser.pkg.PackageName != p2Name {
		t.Errorf("expected package name %s, got %s", p2Name, parser.pkg.PackageName)
	}
	// and the package should default to true for FilesAnalyzed
	if parser.pkg.FilesAnalyzed != true {
		t.Errorf("expected FilesAnalyzed to default to true, got false")
	}
	if parser.pkg.IsFilesAnalyzedTagPresent != false {
		t.Errorf("expected IsFilesAnalyzedTagPresent to default to false, got true")
	}
	// and the package should _not_ be an "unpackaged" placeholder
	if parser.pkg.IsUnpackaged == true {
		t.Errorf("package incorrectly has IsUnpackaged flag set")
	}
	// and the Document's Packages should be of size 2 and have these two
	if len(parser.doc.Packages) != 2 {
		t.Errorf("Expected len(Packages) to be 2, got %d", len(parser.doc.Packages))
	}
	if parser.doc.Packages[0] != p1 {
		t.Errorf("Expected package %v in Packages[0], got %v", p1, parser.doc.Packages[0])
	}
	if parser.doc.Packages[0].PackageName != p1Name {
		t.Errorf("expected package name %s in Packages[0], got %s", p1Name, parser.doc.Packages[0].PackageName)
	}
	if parser.doc.Packages[1] != parser.pkg {
		t.Errorf("Expected package %v in Packages[1], got %v", parser.pkg, parser.doc.Packages[1])
	}
	if parser.doc.Packages[1].PackageName != p2Name {
		t.Errorf("expected package name %s in Packages[1], got %s", p2Name, parser.doc.Packages[1].PackageName)
	}
	// and the first Package's Files should be of size 1 and have f1 only
	if len(parser.doc.Packages[0].Files) != 1 {
		t.Errorf("Expected 1 file in Packages[0].Files, got %d", len(parser.doc.Packages[0].Files))
	}
	if parser.doc.Packages[0].Files[0] != f1 {
		t.Errorf("Expected file %v in Files[0], got %v", f1, parser.doc.Packages[0].Files[0])
	}
	if parser.doc.Packages[0].Files[0].FileName != f1Name {
		t.Errorf("expected file name %s in Files[0], got %s", f1Name, parser.doc.Packages[0].Files[0].FileName)
	}
	// and the second Package should have no files
	if len(parser.doc.Packages[1].Files) != 0 {
		t.Errorf("Expected no files in Packages[1].Files, got %d", len(parser.doc.Packages[1].Files))
	}
	// and the current file should be nil
	if parser.file != nil {
		t.Errorf("Expected nil for parser.file, got %v", parser.file)
	}
	// and the current snippet should be nil
	if parser.snippet != nil {
		t.Errorf("Expected nil for parser.snippet, got %v", parser.snippet)
	}
}

func TestParser2_1SnippetMovesToFileAfterParsingFileNameTag(t *testing.T) {
	p1Name := "package1"
	f1Name := "f1.txt"
	s1Name := "SPDXRef-s1"
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: p1Name},
		file:    &spdx.File2_1{FileName: f1Name},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: s1Name},
	}
	p1 := parser.pkg
	f1 := parser.file
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	f2Name := "f2.txt"
	err := parser.parsePair2_1("FileName", f2Name)
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should be correct
	if parser.st != psFile2_1 {
		t.Errorf("expected state to be %v, got %v", psSnippet2_1, parser.st)
	}
	// and current package should remain what it was
	if parser.pkg != p1 {
		t.Fatalf("expected package to remain %v, got %v", p1, parser.pkg)
	}
	// and a file should be created
	if parser.file == nil {
		t.Fatalf("parser didn't create new file")
	}
	// and the file name should be as expected
	if parser.file.FileName != f2Name {
		t.Errorf("expected file name %s, got %s", f2Name, parser.file.FileName)
	}
	// and the Package's Files should be of size 2 and have these two
	if parser.pkg.Files[0] != f1 {
		t.Errorf("Expected file %v in Files[0], got %v", f1, parser.pkg.Files[0])
	}
	if parser.pkg.Files[0].FileName != f1Name {
		t.Errorf("expected file name %s in Files[0], got %s", f1Name, parser.pkg.Files[0].FileName)
	}
	if parser.pkg.Files[1] != parser.file {
		t.Errorf("Expected file %v in Files[1], got %v", parser.file, parser.pkg.Files[1])
	}
	if parser.pkg.Files[1].FileName != f2Name {
		t.Errorf("expected file name %s in Files[1], got %s", f2Name, parser.pkg.Files[1].FileName)
	}
	// and the current snippet should be nil
	if parser.snippet != nil {
		t.Errorf("Expected nil for parser.snippet, got %v", parser.snippet)
	}
}

func TestParser2_1SnippetMovesToOtherLicenseAfterParsingLicenseIDTag(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: "SPDXRef-s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	err := parser.parsePair2_1("LicenseID", "LicenseRef-TestLic")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psOtherLicense2_1 {
		t.Errorf("expected state to be %v, got %v", psOtherLicense2_1, parser.st)
	}
}

func TestParser2_1SnippetMovesToReviewAfterParsingReviewerTag(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: "SPDXRef-s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	err := parser.parsePair2_1("Reviewer", "Person: John Doe")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psReview2_1 {
		t.Errorf("expected state to be %v, got %v", psReview2_1, parser.st)
	}
}

func TestParser2_1SnippetStaysAfterParsingRelationshipTags(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: "SPDXRef-s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	err := parser.parsePair2_1("Relationship", "blah CONTAINS blah-else")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should remain unchanged
	if parser.st != psSnippet2_1 {
		t.Errorf("expected state to be %v, got %v", psSnippet2_1, parser.st)
	}
	// and the relationship should be in the Document's Relationships
	if len(parser.doc.Relationships) != 1 {
		t.Fatalf("expected doc.Relationships to have len 1, got %d", len(parser.doc.Relationships))
	}
	if parser.doc.Relationships[0].RefA != "blah" {
		t.Errorf("expected RefA to be %s, got %s", "blah", parser.doc.Relationships[0].RefA)
	}

	err = parser.parsePair2_1("RelationshipComment", "blah")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	// state should still remain unchanged
	if parser.st != psSnippet2_1 {
		t.Errorf("expected state to be %v, got %v", psSnippet2_1, parser.st)
	}
}

func TestParser2_1SnippetStaysAfterParsingAnnotationTags(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: "SPDXRef-s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	err := parser.parsePair2_1("Annotator", "Person: John Doe ()")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psSnippet2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet2_1)
	}

	err = parser.parsePair2_1("AnnotationDate", "2018-09-15T00:36:00Z")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psSnippet2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet2_1)
	}

	err = parser.parsePair2_1("AnnotationType", "REVIEW")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psSnippet2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet2_1)
	}

	err = parser.parsePair2_1("SPDXREF", "SPDXRef-45")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psSnippet2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet2_1)
	}

	err = parser.parsePair2_1("AnnotationComment", "i guess i had something to say about this particular file")
	if err != nil {
		t.Errorf("got error when calling parsePair2_1: %v", err)
	}
	if parser.st != psSnippet2_1 {
		t.Errorf("parser is in state %v, expected %v", parser.st, psSnippet2_1)
	}

	// and the annotation should be in the Document's Annotations
	if len(parser.doc.Annotations) != 1 {
		t.Fatalf("expected doc.Annotations to have len 1, got %d", len(parser.doc.Annotations))
	}
	if parser.doc.Annotations[0].Annotator != "John Doe ()" {
		t.Errorf("expected Annotator to be %s, got %s", "John Doe ()", parser.doc.Annotations[0].Annotator)
	}
}

// ===== Snippet data section tests =====
func TestParser2_1CanParseSnippetTags(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	// Snippet SPDX Identifier
	err := parser.parsePairFromSnippet2_1("SnippetSPDXID", "SPDXRef-s1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetSPDXIdentifier != "SPDXRef-s1" {
		t.Errorf("got %v for SnippetSPDXIdentifier", parser.snippet.SnippetSPDXIdentifier)
	}

	// Snippet from File SPDX Identifier
	err = parser.parsePairFromSnippet2_1("SnippetFromFileSPDXID", "SPDXRef-f1")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetFromFileSPDXIdentifier != "SPDXRef-f1" {
		t.Errorf("got %v for SnippetFromFileSPDXIdentifier", parser.snippet.SnippetFromFileSPDXIdentifier)
	}

	// Snippet Byte Range
	err = parser.parsePairFromSnippet2_1("SnippetByteRange", "20:320")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetByteRangeStart != 20 {
		t.Errorf("got %v for SnippetByteRangeStart", parser.snippet.SnippetByteRangeStart)
	}
	if parser.snippet.SnippetByteRangeEnd != 320 {
		t.Errorf("got %v for SnippetByteRangeEnd", parser.snippet.SnippetByteRangeEnd)
	}

	// Snippet Line Range
	err = parser.parsePairFromSnippet2_1("SnippetLineRange", "5:12")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetLineRangeStart != 5 {
		t.Errorf("got %v for SnippetLineRangeStart", parser.snippet.SnippetLineRangeStart)
	}
	if parser.snippet.SnippetLineRangeEnd != 12 {
		t.Errorf("got %v for SnippetLineRangeEnd", parser.snippet.SnippetLineRangeEnd)
	}

	// Snippet Concluded License
	err = parser.parsePairFromSnippet2_1("SnippetLicenseConcluded", "BSD-3-Clause")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetLicenseConcluded != "BSD-3-Clause" {
		t.Errorf("got %v for SnippetLicenseConcluded", parser.snippet.SnippetLicenseConcluded)
	}

	// License Information in Snippet
	lics := []string{
		"Apache-2.0",
		"GPL-2.0-or-later",
		"CC0-1.0",
	}
	for _, lic := range lics {
		err = parser.parsePairFromSnippet2_1("LicenseInfoInSnippet", lic)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
	for _, licWant := range lics {
		flagFound := false
		for _, licCheck := range parser.snippet.LicenseInfoInSnippet {
			if licWant == licCheck {
				flagFound = true
			}
		}
		if flagFound == false {
			t.Errorf("didn't find %s in LicenseInfoInSnippet", licWant)
		}
	}
	if len(lics) != len(parser.snippet.LicenseInfoInSnippet) {
		t.Errorf("expected %d licenses in LicenseInfoInSnippet, got %d", len(lics),
			len(parser.snippet.LicenseInfoInSnippet))
	}

	// Snippet Comments on License
	err = parser.parsePairFromSnippet2_1("SnippetLicenseComments", "this is a comment about the licenses")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetLicenseComments != "this is a comment about the licenses" {
		t.Errorf("got %v for SnippetLicenseComments", parser.snippet.SnippetLicenseComments)
	}

	// Snippet Copyright Text
	err = parser.parsePairFromSnippet2_1("SnippetCopyrightText", "copyright (c) John Doe and friends")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetCopyrightText != "copyright (c) John Doe and friends" {
		t.Errorf("got %v for SnippetCopyrightText", parser.snippet.SnippetCopyrightText)
	}

	// Snippet Comment
	err = parser.parsePairFromSnippet2_1("SnippetComment", "this is a comment about the snippet")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetComment != "this is a comment about the snippet" {
		t.Errorf("got %v for SnippetComment", parser.snippet.SnippetComment)
	}

	// Snippet Name
	err = parser.parsePairFromSnippet2_1("SnippetName", "from some other package called abc")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if parser.snippet.SnippetName != "from some other package called abc" {
		t.Errorf("got %v for SnippetName", parser.snippet.SnippetName)
	}
}

func TestParser2_1SnippetUnknownTagFails(t *testing.T) {
	parser := tvParser2_1{
		doc:     &spdx.Document2_1{},
		st:      psSnippet2_1,
		pkg:     &spdx.Package2_1{PackageName: "package1"},
		file:    &spdx.File2_1{FileName: "f1.txt"},
		snippet: &spdx.Snippet2_1{SnippetSPDXIdentifier: "SPDXRef-s1"},
	}
	parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
	parser.pkg.Files = append(parser.pkg.Files, parser.file)
	parser.file.Snippets = append(parser.file.Snippets, parser.snippet)

	err := parser.parsePairFromSnippet2_1("blah", "something")
	if err == nil {
		t.Errorf("expected error from parsing unknown tag")
	}
}
