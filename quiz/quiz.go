package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func exit(message string) {
	fmt.Fprintf(os.Stderr, "%v", message)
	os.Exit(1)
}

func ReadAndParseCSV(filename string, shuffle bool) []problem {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("failed to open csv file: %s", filename))
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		exit("failed to parse provided csv file")
	}

	problems := make([]problem, len(records))
	for i, line := range records {
		problems[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	return problems
}

func main() {
	var csvFile = flag.String("csv", "problems.csv", "a csv file in the format of 'questions,answers'")
	var timeLimit = flag.Int("limit", 30, "the time limit for the entire quiz in seconds")
	var shuffle = flag.Bool("shuffle", false, "shuffles the problemset before starting the quiz")

	flag.Parse()

	correct := 0
	problems := ReadAndParseCSV(*csvFile, *shuffle)
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

ProblemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s -> ", i+1, p.q)

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()

		select {
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		case <-timer.C:
			fmt.Println()
			break ProblemLoop
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}
