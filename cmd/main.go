// gcov2lcov - convert golang coverage files to the lcov format.
//
// Copyright (c) 2019 Jan Delgado
// Copyright (c) 2019 Richard S Allinson
//
// Credits:
// This tool is based on covfmt (https://github.com/ricallinson/covfmt) and
// uses some parts of goveralls (https://github.com/mattn/goveralls).
//
package main

import (
	"flag"
	"github.com/jandelgado/gcov2lcov"
	"log"
	"os"
)

func main() {
	os.Exit(gcovmain())
}

func gcovmain() int {
	infileName := flag.String("infile", "", "go coverage file to read, default: <stdin>")
	outfileName := flag.String("outfile", "", "lcov file to write, default: <stdout>")
	useAbsoluteSourcePath := flag.Bool("use-absolute-source-path", false,
		"use absolute paths for source file in lcov output, default: false")
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		return 1
	}

	infile := os.Stdin
	outfile := os.Stdout
	var err error
	if *infileName != "" {
		infile, err = os.Open(*infileName)
		if err != nil {
			log.Printf("error opening input file: %v\n", err)
			return 2
		}
		defer infile.Close()
	}
	if *outfileName != "" {
		outfile, err = os.Create(*outfileName)
		if err != nil {
			log.Printf("error opening output file: %v\n", err)
			return 3
		}
		defer outfile.Close()
	}

	var pathResolverFunc gcov2lcov.PathResolver
	if *useAbsoluteSourcePath {
		pathResolverFunc = gcov2lcov.AbsolutePathResolver
	} else {
		pathResolverFunc = gcov2lcov.RelativePathResolver
	}

	err = gcov2lcov.ConvertCoverage(infile, outfile, pathResolverFunc)
	if err != nil {
		log.Printf("error: convert: %v", err)
		return 4
	}
	return 0
}

