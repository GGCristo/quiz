package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	stringFlag   string
	limitFlag    int
	isRandomFlag bool
)

func init() {
	flag.StringVar(&stringFlag, "csv", "problems.csv",
		`a csv file in the format of 'question,answer'`)
	flag.IntVar(&limitFlag, "limit", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&isRandomFlag, "random", true,
		"Randomize the order of the questions")
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	problems, err := readProblems()
	if err != nil {
		log.Fatal(err)
	}

	correctsCounter := run(problems)
	fmt.Printf("\nYou scored %v out of %v.\n", correctsCounter, len(problems))
}

func stopWatch(done chan<- bool) {
	timer := time.NewTimer(time.Second * time.Duration(limitFlag))
	<-timer.C
	done <- true
}

func run(problems [][]string) (correctsCounter int) {
	scanner := bufio.NewScanner(os.Stdin)
	done := make(chan bool)
	fmt.Printf("You have %v seconds. Press enter to start", limitFlag)
	fmt.Scanln()
	go stopWatch(done)
	go func(done chan<- bool) {
		for i, problem := range problems {
			fmt.Printf("Problem #%v: %v = ", i+1, problem[0])
			scanner.Scan()
			answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if answer == problem[1] {
				correctsCounter++
			}
		}
		done <- true
	}(done)
	<-done
	return
}

func readProblems() ([][]string, error) {
	var problems [][]string
	problemsFile, err := os.Open(stringFlag)
	if err != nil {
		return problems, err
	}
	defer problemsFile.Close()
	reader := csv.NewReader(problemsFile)
	problems, err = reader.ReadAll()
	if err != nil {
		return problems, err
	}
	if isRandomFlag {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}
	return problems, nil
}
