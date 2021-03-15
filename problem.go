package kata

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

const problemToken = "PROBLEM:"
const answerToken = "ANSWER:"
const casesToken = "CASES:"

// Case represents a parameter-answer pair
type Case struct {
	Params []string
	Answer string
}

// Problem is a problem statement, with its associated answers
type Problem struct {
	Statement string
	Answer    string
	File      string
	Cases     []Case
}

// NewProblemFilePath creates a Problem from a file path
func NewProblemFilePath(problemFilePath string) (Problem, error) {
	file, err := os.Open(problemFilePath)
	if err != nil {
		return Problem{}, err
	}
	defer file.Close()
	return NewProblem(file, problemFilePath), err
}

// NewProblem creates a Problem from an io.Reader and the file name
func NewProblem(file io.Reader, fileName string) Problem {
	scanner := bufio.NewScanner(file)

	startProblem := -1
	endProblem := false
	startAnswer := -1
	startCases := -1
	statementLines := []string{}
	casePair := []string{}
	casePairs := [][]string{}
	var answer string
	for i := 0; scanner.Scan(); i++ {
		switch scanner.Text() {
		case problemToken:
			startProblem = i
		case answerToken:
			startAnswer = i
			endProblem = true
		case casesToken:
			startCases = i
			endProblem = true
		}
		if startProblem > -1 && startProblem < i && !endProblem {
			statementLines = append(statementLines, scanner.Text())
		} else if i == startAnswer+1 {
			answer = scanner.Text()
		} else if i > startCases {
			current := scanner.Text()
			if current == "" {
				casePairs = append(casePairs, casePair[0:])
				casePair = []string{}
			} else {
				casePair = append(casePair, scanner.Text())
			}
		}
	}
	if len(casePair) > 1 {
		casePairs = append(casePairs, casePair[0:])
	}

	cases := []Case{}
	if len(casePairs) == 0 && answer != "" {
		c := Case{Params: []string{}, Answer: answer}
		cases = append(cases, c)
	} else {
		cases = parseCasePairs(casePairs)
		answer = ""
	}

	return Problem{
		Statement: cleanParagraph(statementLines),
		Answer:    answer,
		File:      fileName,
		Cases:     cases,
	}
}

// Introduction provides the formatted introduction and can be the problem
// statement if present or a formatted file path
func (p Problem) Introduction() string {
	if len(p.Statement) > 0 {
		return p.Statement
	}
	return p.FileIntro()
}

// FileIntro returns a formatted file path
func (p Problem) FileIntro() string {
	return fileAsIntro(p.File)
}

func fileAsIntro(filename string) string {
	parts := strings.Split(filename, "/")
	lastIndex := len(parts) - 1
	fparts := strings.Split(parts[lastIndex], ".")
	if len(fparts) > 1 {
		parts[lastIndex] = strings.Join(fparts[:len(fparts)-1], ".")
	} else {
		parts[lastIndex] = strings.Join(fparts, "")
	}

	return strings.Join(parts, " ")
}

func cleanParagraph(lines []string) string {
	for i := len(lines) - 1; i > -1; i-- {
		hasNext := i < (len(lines) - 1)
		if hasNext {
			next := lines[i+1]
			indented, _ := regexp.MatchString(`^\s`, next)
			if lines[i] != "" && next != "" && !indented {
				lines[i] = lines[i] + " " + lines[i+1]
				lines = append(lines[:i+1], lines[i+2:]...)
			}
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func parseCasePairs(pairs [][]string) []Case {
	cases := []Case{}
	for _, pair := range pairs {
		if len(pair) < 2 {
			continue
		}
		first := pair[0]
		next := pair[1:]
		cases = append(cases, Case{parseCaseParams(first), strings.Join(next, "\n")})
	}
	return cases
}

func parseCaseParams(str string) []string {
	result := []string{}
	for _, param := range strings.Split(str, ",") {
		result = append(result, strings.TrimSpace(param))
	}
	return result
}
