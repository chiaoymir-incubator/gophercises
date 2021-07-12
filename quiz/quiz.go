package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Color string

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	ColorBlack         = "\u001b[30m"
	ColorRed           = "\u001b[31m"
	ColorGreen         = "\u001b[32m"
	ColorYellow        = "\u001b[33m"
	ColorBlue          = "\u001b[34m"
	ColorMagneta       = "\u001b[35m"
	ColorCyan          = "\u001b[36m"
	ColorWhite         = "\u001b[37m"
	ColorGray          = "\u001b[90m"
	ColorBrightRed     = "\u001b[91m"
	ColorBrightGreen   = "\u001b[92m"
	ColorBrightYellow  = "\u001b[93m"
	ColorBrightBlue    = "\u001b[94m"
	ColorBrightMagneta = "\u001b[95m"
	ColorBrightCyan    = "\u001b[96m"
	ColorBrightWhite   = "\u001b[97m"
	ColorReset         = "\u001b[0m"
)

const (
	CORRECT_POINT = 3
	PASS_POINT    = 1
	WRONG_POINT   = -1
)

type QuizRecord struct {
	Q string
	A int
}

func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func PrintResult(totalQ, correct, pass, wrong int) {
	fmt.Println("")
	fmt.Println("Result:")

	totalQMsg := fmt.Sprintf("  Number of Questions: %d", totalQ)
	colorize(ColorBrightWhite, totalQMsg)

	correctMsg := fmt.Sprintf("  Correct: %d", correct)
	wrongMsg := fmt.Sprintf("  Wrong: %d", wrong)
	passMsg := fmt.Sprintf("  Pass: %d", pass)
	colorize(ColorBrightGreen, correctMsg)
	colorize(ColorBrightRed, wrongMsg)
	colorize(ColorBrightYellow, passMsg)

	totalPoints := correct*CORRECT_POINT + pass*PASS_POINT + wrong*WRONG_POINT

	totalPointsMsg := fmt.Sprintf("  Total Points: %d", totalPoints)
	colorize(ColorBrightCyan, totalPointsMsg)
}

func readCSV(fileName *string) *os.File {
	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, *fileName)

	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Sprintf("Failed to open the CSV file: %s!", filePath)
		colorize(ColorRed, msg)
		os.Exit(1)
	}

	return file
}

func parseCSV(file *os.File) []QuizRecord {
	quizs := []QuizRecord{}
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		Q, a := record[0], record[1]

		A, err := strconv.Atoi(a)

		if err != nil {
			msg := fmt.Sprintf("The record answer has wrong format: %d", A)
			colorize(ColorRed, msg)
		}

		quizs = append(quizs, QuizRecord{
			Q: Q,
			A: A,
		})

	}

	return quizs
}

func main() {
	fmt.Println("Quiz Game Start!")

	// fileName := "problems.csv"
	fileName := flag.String("csv", "problems.csv", "the input csv file, the format is question,answer.")
	timeLimit := flag.Int("limit", 30, "the question time limit")
	flag.Parse()

	if *timeLimit < 1 {
		colorize(ColorRed, "The time should be positive")
		return
	}

	file := readCSV(fileName)
	quizs := parseCSV(file)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quizs), func(i, j int) { quizs[i], quizs[j] = quizs[j], quizs[i] })

	scanner := bufio.NewScanner(os.Stdin)
	correct := 0
	wrong := 0
	pass := 0
	for i, q := range quizs {
		fmt.Printf("Problem #%d: %s = ", i+1, q.Q)
		ch := make(chan string, 1)

		go func() {
			if scanner.Scan() {
				s := strings.TrimSpace(scanner.Text())
				ch <- s
			}
		}()

		select {
		case s := <-ch:
			num, err := strconv.Atoi(s)

			if err != nil {
				colorize(ColorYellow, "Please enter a number!")
				pass++
				continue
			}

			if num != q.A {
				colorize(ColorRed, "Wrong!")
				wrong++
				continue
			}

			colorize(ColorGreen, "Correct!")
			correct++
		case <-time.After(time.Duration(*timeLimit) * time.Second):
			fmt.Println("")
			colorize(ColorRed, "Time out!")
			PrintResult(len(quizs), correct, 0, len(quizs)-correct)
			return
		}
	}

	PrintResult(len(quizs), correct, pass, wrong)

}
