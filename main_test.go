// gcov2lcov - convert golang coverage files to the lcov format.
// (c) 2019 Jan Delgado
package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// note: in this integrative test the package path must match the actual
	// repository name of this project.
	cov := `mode: set
github.com/jandelgado/gcov2lcov/main.go:6.14,8.3 2 1`

	reader := strings.NewReader(cov)
	res := parseCoverage(reader)

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
