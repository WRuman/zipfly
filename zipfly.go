package main

/**
	Zipfly - A utility for archiving files of any size with a small memory footprint
	Copyright (C) 2016 Will Ruman and Scott Parsons 

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
)

const BLOCKSIZE int64 = 50 * 1024 * 1024 // 50 Mb

// Returns a handle to an existing file
func getInfile(filename string) *os.File {
	f, err := os.Open(filename)
	check(err)
	return f
}

// Returns a handle to a new file
func getOutfile(filename string) *os.File {
	f, err := os.Create(filename)
	check(err)
	return f
}

// Compresses the contents of infile BLOCKSIZE-bytes at 
// a time and stores the compressed data in outfile
func compress(infile, outfile *os.File) {
	reader := bufio.NewReader(infile)
	zipper := gzip.NewWriter(outfile)
	defer zipper.Close()

	var buf = make([]byte, BLOCKSIZE)
	var done = false
	for !done {
		v, _ := reader.Read(buf)
		// No bytes read, signal loop termination
		if v == 0 {
			done = true
			continue
		}
		// Compress the exact number of bytes read into buf,
		// which can be <= len(buf)
		zipper.Write(buf[:v - 1])
	}
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	flag.Parse()
	fname := flag.Arg(0)
	if len(fname) < 1 {
		fmt.Println("Usage : zipfly <filename>")
		os.Exit(1)
	}
	// Get file handles
	var dest = fname + ".gz"
	infile := getInfile(fname) // 
	outfile := getOutfile(dest)
	compress(infile, outfile)
	// Close all files before doing any deletions
	infile.Close()
	outfile.Close()

	fmt.Printf("%s compressed to %s\n", fname, dest)
	// Remove original input file
	rmerr := os.RemoveAll(fname)
	check(rmerr)
}


