///bin/true; exec /usr/bin/env go run "$0" "$@"
// -*- Mode: Go; indent-tabs-mode: t; tab-width: 8 -*-
/*
 * Copyright 2015 Michael Terry
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if _, err := exec.LookPath("pbnsolve"); err != nil {
		fmt.Println("You forgot to set the PATH for pbnsolve")
		os.Exit(1)
	}

	n := flag.Uint("n", 1, "Number of puzzles to generate")
	d := flag.String("d", "1+", "Difficulty of puzzles to generate")
	flag.Parse()

	if len(flag.Args()) != 3 {
		fmt.Println(flag.Args())
		fmt.Println("Usage: findpuzzle.go HEIGHT WIDTH SEED")
		os.Exit(1)
	}
	height, err := strconv.ParseUint(flag.Arg(0), 10, 0)
	if err != nil {
		fmt.Println("Could not understand height")
		os.Exit(1)
	}
	width, err := strconv.ParseUint(flag.Arg(1), 10, 0)
	if err != nil {
		fmt.Println("Could not understand width")
		os.Exit(1)
	}
	seed, err := strconv.ParseInt(flag.Arg(2), 10, 64)
	if err != nil {
		fmt.Println("Could not understand the seed")
		os.Exit(1)
	}

	minDifficulty := 1
	maxDifficulty := 0

	if (*d)[len(*d)-1] == '+' {
		minDifficulty, _ = strconv.Atoi((*d)[:len(*d)-1])
	} else {
		minDifficulty, _ = strconv.Atoi(*d)
		maxDifficulty = minDifficulty
	}

	maxSeedLength := len(fmt.Sprintf("%v", uint64(math.MaxUint64)))

	for count := uint(0); count < *n; {
		cmd := exec.Command("go", "run", "./generate.go", fmt.Sprintf("%v", height), fmt.Sprintf("%v", width), fmt.Sprintf("%v", seed))
		c := make(chan string)
		go generate(height, width, seed, c)
		output := []byte(<-c)

		file, _ := ioutil.TempFile(".", "")
		file.Write(output)
		file.Close()

		cmd = exec.Command("pbnsolve", "-u", "-aLHEC", "-t", "-f", "non", "-d", "5", file.Name())
		output, err = cmd.Output()
		if err != nil {
			panic(err)
		}

		if ok, _ := regexp.Match(".*UNIQUE.*", output); ok {
			diffexp, _ := regexp.Compile("\nLines Processed:.* .((?m:.*))00%")
			matches := diffexp.FindSubmatch(output)
			difficulty, _ := strconv.Atoi(string(matches[1]))

			if difficulty >= minDifficulty && (maxDifficulty == 0 || difficulty <= maxDifficulty) {
				newfile := fmt.Sprintf("%02vx%02v.%03v.%0*v.non", height, width, difficulty, maxSeedLength, uint64(seed))
				os.Rename(file.Name(), newfile)
				fmt.Println("Wrote", newfile)
				count++
			}
		}
		seed++

		os.Remove(file.Name())
	}
}

func generate(height, width uint64, seed int64, c chan<- string) {
	goal := ""
	for i := uint64(0); i < width*height; i++ {
		goal += strconv.Itoa(rand.Intn(2))
	}

	rowHints := make([]string, height)
	for y := uint64(0); y < height; y++ {
		row := goal[y*width : y*width+width]
		rowChunks := strings.Split(row, "0")
		var currentRowHints []string
		for _, chunk := range rowChunks {
			if l := len(chunk); l > 0 {
				currentRowHints = append(currentRowHints, strconv.Itoa(l))
			}
		}
		rowHints[y] = strings.Join(currentRowHints, ",")
	}

	colHints := make([]string, width)
	for x := uint64(0); x < width; x++ {
		col := ""
		for y := uint64(0); y < height; y++ {
			col += string(goal[y*width+x])
		}
		colChunks := strings.Split(col, "0")
		var currentColHints []string
		for _, chunk := range colChunks {
			if l := len(chunk); l > 0 {
				currentColHints = append(currentColHints, strconv.Itoa(l))
			}
		}
		colHints[x] = strings.Join(currentColHints, ",")
	}

	content := fmt.Sprintf(`catalogue "generate.go with seed %v"
title "Random #%v"
by "Michael Terry"
copyright "© 2015 Michael Terry"
license CC-BY-SA-4.0
height %v
width %v

rows
%s

columns
%s

goal "%v"`, seed, seed, height, width, strings.Join(rowHints, "\n"), strings.Join(colHints, "\n"), goal)
	c <- content
}
