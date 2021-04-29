// gcov2lcov - convert golang coverage files to the lcov format.
//
// Copyright (c) 2019 Jan Delgado
// Copyright (c) 2019 Richard S Allinson
//
// Credits:
// This tool is based on covfmt (https://github.com/ricallinson/covfmt) and
// uses some parts of goveralls (https://github.com/mattn/goveralls).
//
package gcov2lcov

import (
	"bufio"
	"errors"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type block struct {
	startLine  int
	startChar  int
	endLine    int
	endChar    int
	statements int
	covered    int
}

var vscDirs = []string{".git", ".hg", ".bzr", ".svn"}

type cacheEntry struct {
	file string
	err  error
}

var pkgCache = map[string]cacheEntry{}

func ConvertCoverage(in io.Reader, out io.Writer, pathResolverFunc func(string) string) error {
	blocks, err := parseCoverage(in, pathResolverFunc)
	if err != nil {
		return err
	}
	return writeLcov(blocks, out)
}

func AbsolutePathResolver(name string) string {
	return name
}

func RelativePathResolver(name string) string {
	if dir, ok := findRepositoryRoot(name); ok {
		filename := strings.TrimPrefix(name, dir+string(os.PathSeparator))
		return filename
	}
	return name
}

// given a module+file spec (e.g. github.com/jandelgado/gcov2lcov/main.go),
// strip of the module name and return the file name (e.g. main.go).
func findFile(filePath string) (string, error) {
	dir, file := filepath.Split(filePath)
	var result cacheEntry
	var ok bool
	if result, ok = pkgCache[filePath]; !ok {
		pkg, err := build.Import(dir, ".", build.FindOnly)
		if err == nil {
			result = cacheEntry{filepath.Join(pkg.Dir, file), nil}
		} else {
			result = cacheEntry{"", err}
		}
		pkgCache[filePath] = result
	}
	return result.file, result.err
}

// findRepositoryRoot finds the VCS root dir of a given dir
func findRepositoryRoot(dir string) (string, bool) {
	for _, vcsdir := range vscDirs {
		if d, err := os.Stat(filepath.Join(dir, vcsdir)); err == nil && d.IsDir() {
			return dir, true
		}
	}
	nextdir := filepath.Dir(dir)
	if nextdir == dir {
		return "", false
	}
	return findRepositoryRoot(nextdir)
}

func keysOfMap(m map[int]int) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func writeLcovRecord(filePath string, blocks []*block, w io.StringWriter) error {

	writer := func(err error, s string) error {
		if err != nil {
			return err
		}
		_, err = w.WriteString(s)
		return err
	}
	var err error
	err = writer(err, "TN:\nSF:"+filePath+"\n")

	// Loop over functions
	// FN: line,name

	// FNF: total functions
	// FNH: covered functions

	// Loop over functions
	// FNDA: stats,name ?

	total := 0
	covered := 0

	// maps line number to sum of covered
	coverMap := map[int]int{}

	// Loop over each block and extract the lcov data needed.
	for _, b := range blocks {
		// For each line in a block we add an lcov entry and count the lines.
		for i := b.startLine; i <= b.endLine; i++ {
			coverMap[i] += b.covered
		}
	}

	lines := keysOfMap(coverMap)
	sort.Ints(lines)
	for _, line := range lines {
		err = writer(err, "DA:"+strconv.Itoa(line)+","+strconv.Itoa(coverMap[line])+"\n")
		total++
		if coverMap[line] > 0 {
			covered++
		}
	}

	// LH:<number of lines with a non-zero execution count>
	// LF:<number of instrumented lines>
	err = writer(err, "LF:"+strconv.Itoa(total)+"\n")
	err = writer(err, "LH:"+strconv.Itoa(covered)+"\n")

	// Loop over branches
	// BRDA: ?

	// BRF: total branches
	// BRH: covered branches

	return writer(err, "end_of_record\n")
}

func writeLcov(blocks map[string][]*block, f io.Writer) error {
	w := bufio.NewWriter(f)
	for file, fileBlocks := range blocks {
		if err := writeLcovRecord(file, fileBlocks, w); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

// Format being parsed is:
//   name.go:line.column,line.column numberOfStatements count
// e.g.
//   github.com/jandelgado/golang-ci-template/main.go:6.14,8.2 1 1
func parseCoverageLine(line string) (string, *block, error) {
	path := strings.Split(line, ":")
	if len(path) != 2 {
		return "", nil, errors.New("unexpected format (path sep): " + line)
	}
	parts := strings.Split(path[1], " ")
	if len(parts) != 3 {
		return "", nil, errors.New("unexpected format (parts): " + line)
	}
	sections := strings.Split(parts[0], ",")
	if len(sections) != 2 {
		return "", nil, errors.New("unexpected format (pos): " + line)
	}
	start := strings.Split(sections[0], ".")
	end := strings.Split(sections[1], ".")

	safeAtoi := func(err error, s string) (int, error) {
		if err != nil {
			return 0, err
		}
		return strconv.Atoi(s)
	}
	b := &block{}
	var err error
	b.startLine, err = safeAtoi(nil, start[0])
	b.startChar, err = safeAtoi(err, start[1])
	b.endLine, err = safeAtoi(err, end[0])
	b.endChar, err = safeAtoi(err, end[1])
	b.statements, err = safeAtoi(err, parts[1])
	b.covered, err = safeAtoi(err, parts[2])

	return path[0], b, err
}

func parseCoverage(coverage io.Reader, pathResolverFunc func(string) string) (map[string][]*block, error) {
	scanner := bufio.NewScanner(coverage)
	blocks := map[string][]*block{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		if f, b, err := parseCoverageLine(line); err == nil {
			f, err := findFile(f)
			if err != nil {
				log.Printf("warn: %v", err)
				continue
			}

			f = pathResolverFunc(f)

			// Make sure the filePath is a key in the map.
			if _, found := blocks[f]; !found {
				blocks[f] = []*block{}
			}
			blocks[f] = append(blocks[f], b)
		} else {
			log.Printf("warn: %v", err)
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return blocks, nil
}
