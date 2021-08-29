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
	"sync"
	"time"
)

var wg sync.WaitGroup

var (
	stringFlag   string
	limitFlag    int
	isRandomFlag bool
)

func init() {
	flag.StringVar(&stringFlag, "csv", "problems.csv",
		`a csv file in the format of 'question,answer'`)
	flag.IntVar(&limitFlag, "limit", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&isRandomFlag, "random", true, "Randomize the order of the questions")
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	problems, err := readProblems()
	if err != nil {
		log.Fatal(err)
	}

	correctsCounter := 0
	fmt.Printf("You have %v seconds. Press enter to start", limitFlag)
	fmt.Scanln()
	wg.Add(1)
	go run(problems, &correctsCounter)
	go stopWatch(float64(limitFlag))
	wg.Wait()
	fmt.Printf("\nYou scored %v out of %v.\n", correctsCounter, len(problems))
}

func stopWatch(limitFlag float64) {
	start := time.Now()
	for time.Since(start).Seconds() < limitFlag {
		time.Sleep(time.Second)
	}
	wg.Done()
}

func run(problems [][]string, correctsCounter *int) {
	scanner := bufio.NewScanner(os.Stdin)
	for i, problem := range problems {
		fmt.Printf("Problem #%v: %v = ", i, problem[0])
		scanner.Scan()
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer == problem[1] {
			(*correctsCounter)++
		}
	}
	wg.Done()
}

func readProblems() ([][]string, error) {
	problemsFile, err := os.Open(stringFlag)
	defer problemsFile.Close()
	var problems [][]string
	if err != nil {
		return problems, err
	}
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
