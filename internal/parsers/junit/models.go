package junit

import (
	"encoding/xml"
)

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Properties struct {
	Property []Property `xml:"property"`
}

type TestCase struct {
	Name       string     `xml:"name,attr"`
	ClassName  string     `xml:"classname,attr"`
	Assertions int        `xml:"assertions,attr"`
	Time       float64    `xml:"time,attr"`
	File       string     `xml:"file,attr"`
	Line       int        `xml:"line,attr"`
	Skipped    *Skipped   `xml:"skipped"`
	Failure    *Failure   `xml:"failure"`
	Error      *Error     `xml:"error"`
	SystemOut  string     `xml:"system-out"`
	SystemErr  string     `xml:"system-err"`
	Properties Properties `xml:"properties"`
}

type Skipped struct {
	Message string `xml:"message,attr"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type Error struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr"`
	Tests      int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Errors     int        `xml:"errors,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Assertions int        `xml:"assertions,attr"`
	Time       float64    `xml:"time,attr"`
	File       string     `xml:"file,attr"`
	Properties Properties `xml:"properties"`
	SystemOut  string     `xml:"system-out"`
	SystemErr  string     `xml:"system-err"`
	TestCases  []TestCase `xml:"testcase"`
}

type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Name       string      `xml:"name,attr"`
	Tests      int         `xml:"tests,attr"`
	Failures   int         `xml:"failures,attr"`
	Errors     int         `xml:"errors,attr"`
	Skipped    int         `xml:"skipped,attr"`
	Assertions int         `xml:"assertions,attr"`
	Time       float64     `xml:"time,attr"`
	TestSuites []TestSuite `xml:"testsuite"`
}
