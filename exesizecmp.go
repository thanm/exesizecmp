// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

// Overview: given a set of object files, look for definitions and references
// to import symbols.

const DefaultDumper = "llvm-objdump-14"

var inputsflag = flag.String("i", "", "Comma-separated list of input files (omit to read from stdin)")
var objdumpflag = flag.String("objdump", DefaultDumper, "Name of objdump program to invoke")

func usage(msg string) {
	if len(msg) > 0 {
		fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	}
	fmt.Fprintf(os.Stderr, "usage: execsizecmp [flags] -i=X,Y\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func secSizes(infile string) map[string]int64 {
	cmd := exec.Command(*objdumpflag,
		"-h", // section headers
		"--wide",
		infile)
	out, err := cmd.Output()
	if err != nil {
		fatal("running %s on %s: %v", *objdumpflag, infile, err)
	}

	sm := make(map[string]int64)
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	secre := regexp.MustCompile(`^\s+([0-9]+)\s+(\S+)\s+(\S+)\s+.*`)
	for sc.Scan() {
		line := sc.Text()
		if line == "Sections:" {
			// eat next line with titles
			sc.Scan()
			for sc.Scan() {
				line := sc.Text()
				m := secre.FindStringSubmatch(line)
				if len(m) == 0 {
					fatal("bad match on sections line %q", line)
				}
				sname := m[2]
				ssize := m[3]
				var ssiz int64
				if n, err := fmt.Sscanf(ssize, "%x", &ssiz); n != 1 || err != nil {
					fatal("can't parse sec size in line %s in sections table", line)
				}
				if ssiz == 0 {
					continue
				}
				sm[sname] = ssiz
			}
		}
	}
	return sm
}

func diffTabs(m1 map[string]int64, m2 map[string]int64) {
	ks1 := []string{}
	for k := range m1 {
		ks1 = append(ks1, k)
	}
	tabber := tabwriter.NewWriter(os.Stderr, 1, 8, 1, '\t', 0)
	defer tabber.Flush()
	sort.Strings(ks1)
	for _, k := range ks1 {
		s1 := m1[k]
		s2 := m2[k]
		diff := s2 - s1
		if diff == 0 {
			continue
		}
		perc := float64(diff) / float64(s1)
		fmt.Fprintf(tabber, "%s\t%d\t%d\t%d\tp=%2.1f%%\n", k, s1, s2, diff, perc)
	}
}

func main() {
	flag.Parse()
	if *inputsflag == "" {
		usage("supply input files with -i option")
	}
	infiles := strings.Split(*inputsflag, ",")
	if len(infiles) != 2 {
		usage("supply exactly two input files with -i option")
	}
	tab1, tab2 := secSizes(infiles[0]), secSizes(infiles[1])
	diffTabs(tab1, tab2)
}
