// gcov2lcov - convert golang coverage files to the lcov format.
// (c) 2019 Jan Delgado
package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeysOfMapReturnsAllKeysOfMap(t *testing.T) {
	m := map[int]int{1: 10, 10: 100}

	keys := keysOfMap(m)
	assert.Contains(t, keys, 1)
	assert.Contains(t, keys, 10)
	assert.Equal(t, 2, len(keys))
}

func TestParseCoverageLineFailsOnInvalidLines(t *testing.T) {
	_, _, err := parseCoverageLine("main.go")
	assert.NotNil(t, err)

	_, _, err = parseCoverageLine("main.go:A B")
	assert.NotNil(t, err)

	_, _, err = parseCoverageLine("main.go:A B C")
	assert.NotNil(t, err)

	_, _, err = parseCoverageLine("main.go:6.14,8.3 X 1")
	assert.NotNil(t, err)
}

func TestParseCoverageLineOfParsesValidLineCorrectly(t *testing.T) {
	line := "github.com/jandelgado/gcov2lcov/main.go:6.14,8.3 2 1"
	file, b, err := parseCoverageLine(line)

	assert.Nil(t, err)
	assert.Equal(t, "github.com/jandelgado/gcov2lcov/main.go", file)
	assert.Equal(t, 6, b.startLine)
	assert.Equal(t, 14, b.startChar)
	assert.Equal(t, 8, b.endLine)
	assert.Equal(t, 3, b.endChar)
	assert.Equal(t, 2, b.statements)
	assert.Equal(t, 1, b.covered)
}

func TestParseCoverage(t *testing.T) {

	// note: in this integrative test, the package path must match the actual
	// repository name of this project.
	cov := `mode: set
github.com/jandelgado/gcov2lcov/main.go:6.14,8.3 2 1`

	reader := strings.NewReader(cov)
	res, err := parseCoverage(reader)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	for k, blks := range res {
		assert.Equal(t, 1, len(blks))
		b := blks[0]
		assert.Equal(t, "main.go", k)
		assert.Equal(t, 6, b.startLine)
		assert.Equal(t, 14, b.startChar)
		assert.Equal(t, 8, b.endLine)
		assert.Equal(t, 3, b.endChar)
		assert.Equal(t, 2, b.statements)
		assert.Equal(t, 1, b.covered)
	}
}

func TestConvertCoverage(t *testing.T) {
	// note: in this integrative test, the package path must match the actual
	// repository name of this project. Format:
	//   name.go:line.column,line.column numberOfStatements count
	cov := `mode: set
github.com/jandelgado/gcov2lcov/main.go:6.14,8.3 2 1
github.com/jandelgado/gcov2lcov/main.go:7.14,9.3 2 0
github.com/jandelgado/gcov2lcov/main.go:10.1,11.10 2 2`

	in := strings.NewReader(cov)
	out := bytes.NewBufferString("")
	err := convertCoverage(in, out)

	expected := `TN:
SF:main.go
DA:6,1
DA:7,1
DA:8,1
DA:9,0
DA:10,2
DA:11,2
LF:6
LH:5
end_of_record
`
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}
