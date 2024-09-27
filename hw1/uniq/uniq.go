package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func printLine(output io.Writer, line string) {
	if output != nil {
		fmt.Fprintln(output, line)
	} else {
		fmt.Println(line)
	}
}

func processFile(input io.Reader, output io.Writer, countFlag, duplicatesFlag, uniqueFlag, ignoreCase bool, fieldCount, charCount int) {
	scanner := bufio.NewScanner(input)
	if (countFlag && duplicatesFlag) || (countFlag && uniqueFlag) || (duplicatesFlag && uniqueFlag) {
		printLine(output, "Error: Options -c, -d, and -u are mutually exclusive.")
		return
	}

	lastLine := ""
	cur_in := 0
	origLast := ""
	origCur := ""

	for scanner.Scan() {
		line := scanner.Text()
		origCur = line
		if ignoreCase {
			line = strings.ToLower(line)
		}

		fields := strings.Fields(line)

		if fieldCount > 0 && len(fields) > fieldCount {
			line = strings.Join(fields[fieldCount:], " ")
		} else if fieldCount > 0 {
			line = ""
		}

		if charCount > 0 && len(line) > charCount {
			line = line[charCount:]
		}

		if cur_in == 0 {
			origLast = origCur
		}
		if line != lastLine && cur_in != 0 {
			if countFlag {
				printLine(output, fmt.Sprintf("%d %s", cur_in, origLast))
			} else if duplicatesFlag {
				if cur_in > 1 {
					printLine(output, origLast)
				}
			} else if uniqueFlag {
				if cur_in == 1 {
					printLine(output, origLast)
				}
			} else {
				printLine(output, origLast)
			}
			cur_in = 0
			origLast = origCur
		}

		cur_in++
		lastLine = line
	}

	if cur_in > 0 { // Обработка последней строки
		if countFlag {
			printLine(output, fmt.Sprintf("%d %s", cur_in, origLast))
		} else if duplicatesFlag {
			if cur_in > 1 {
				printLine(output, origLast)
			}
		} else if uniqueFlag {
			if cur_in == 1 {
				printLine(output, origLast)
			}
		} else {
			printLine(output, origLast)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(output, "Error reading input:", err)
		return
	}
}

func main() {
	countFlag := flag.Bool("c", false, "count consecutive occurrences of each line")
	duplicatesFlag := flag.Bool("d", false, "only print duplicate lines")
	uniqueFlag := flag.Bool("u", false, "only print unique lines")
	ignoreCase := flag.Bool("i", false, "ignore case differences")
	fieldCount := flag.Int("f", 0, "ignore the first num_fields fields")
	charCount := flag.Int("s", 0, "ignore the first num_chars characters")

	flag.Parse()

	var output io.Writer

	if len(flag.Args()) > 1 {
		outputFile, err := os.Create(flag.Args()[1])
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer outputFile.Close()
		output = outputFile
	} else {
		output = os.Stdout // Вывод в консоль
	}

	var input io.Reader
	if len(flag.Args()) > 0 {
		inputFile, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer inputFile.Close()
		input = inputFile
	} else {
		input = os.Stdin
	}

	processFile(input, output, *countFlag, *duplicatesFlag, *uniqueFlag, *ignoreCase, *fieldCount, *charCount)
}
