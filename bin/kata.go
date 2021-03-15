package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/andreyvit/diff"
	"github.com/asartalo/kata"
	"github.com/fatih/color"
	do "github.com/gbevan/godo"
)

func timeTrack(start time.Time, run func()) time.Duration {
	run()
	return time.Since(start)
}

func showAnswer(haveProblemFormal bool, expected string, answer string) {
	suffix := ""
	correct := false
	var output string
	if haveProblemFormal {
		if answer == expected {
			suffix = "✔"
			correct = true
			output = fmt.Sprintf("Answer: %s %s", answer, suffix)
		} else {
			suffix = "✘"
			output = fmt.Sprintf(
				"Answer: %s %s\n%s",
				answer,
				suffix,
				diff.LineDiff(expected, answer),
			)
		}
	}
	colorAnswer(haveProblemFormal, correct, output)
}

func showDuration(duration time.Duration) {
	fmt.Println("Duration:", duration)
	fmt.Println("")
}

func colorAnswer(haveProblemFormal bool, correct bool, output string) {
	if haveProblemFormal {
		if correct {
			color.Green(output)
		} else {
			color.Red(output)
		}
	} else {
		color.Yellow(output)
	}
}

func checkFormalProblem(problemName string) (kata.Problem, bool) {
	// Check if answer exists
	answerFileName := fmt.Sprintf(`problems/%s.txt`, problemName)
	problem, err := kata.NewProblemFilePath(answerFileName)
	if err != nil {
		return problem, false
	}
	return problem, true
}

func showIntro(fileIntro, introText string) {
	hr := strings.Repeat("-", len(fileIntro))
	fmt.Print("\n\n")
	color.Yellow(fileIntro)
	fmt.Println(hr)
	if fileIntro != introText {
		color.Yellow(introText)
	}
	fmt.Print("\n")
}

type runCallback func(string, time.Duration)

func runCommand(args []string, onSuccess runCallback) {
	var (
		output      []byte
		errorOutput []byte
	)
	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		if output, err = ioutil.ReadAll(stdout); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error executing the script ", err)
			os.Exit(1)
		}
		errorOutput, _ = ioutil.ReadAll(stderr)
		done <- cmd.Wait()
	}()

	timeLimit := 1 * time.Minute
	select {
	case <-time.After(timeLimit):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill: ", err)
		}
		log.Println(fmt.Sprintf(
			"Process killed. Program ran for more than %d seconds",
			timeLimit/time.Second,
		))
	case err := <-done:
		duration := time.Since(start)
		if err != nil {
			log.Printf("process done with error = %v", err)
			if len(output) > 0 {
				fmt.Println(
					strings.TrimSpace(string(output[:])),
				)
			}

			if len(errorOutput) > 0 {
				fmt.Println(string(errorOutput[:]))
			}
			fmt.Println()
		} else {
			onSuccess(strings.TrimSpace(string(output[:])), duration)
		}
	}
}

func answerSuccess(expected string, haveProblemFormal bool) runCallback {
	return func(answer string, duration time.Duration) {
		showAnswer(haveProblemFormal, expected, answer)
		showDuration(duration)
	}
}

func checkCommand(cmd []string, scriptFile string) {
	call := fmt.Sprintf("%s %s", strings.Join(cmd, " "), scriptFile)
	r, _ := regexp.Compile(`solutions/([\w-_/]+/)\w+/([^/_]+)(_\w+)*\.\w+$`)
	if r.MatchString(scriptFile) {
		match := r.FindStringSubmatch(scriptFile)
		problemName := match[1] + match[2]
		problem, haveProblemFormal := checkFormalProblem(problemName)
		showIntro(problem.FileIntro(), problem.Introduction())
		for _, c := range problem.Cases {
			all := append(cmd, scriptFile)
			all = append(all, c.Params...)
			cmdTxt := strings.Join(all, " ")
			if len(cmdTxt) > 300 {
				cmdTxt = fmt.Sprintf(`%.300s...`, cmdTxt)
			}
			color.Magenta(cmdTxt)
			fmt.Print("\n")
			if len(c.Params) > 0 {
				txtArgs := strings.Join(c.Params, ", ")
				if len(txtArgs) > 300 {
					color.Yellow(fmt.Sprintf(`Given: %.300s...`, txtArgs))
				} else {
					color.Yellow(fmt.Sprintf(`Given: %s`, txtArgs))
				}
			}
			runCommand(all, answerSuccess(c.Answer, haveProblemFormal))
		}
	} else {
		color.Red(fmt.Sprintf("Don't know how to run %s", call))
	}
}

var runners = map[string][]string{
	"go": {"go", "run"},
	"js": {"node"},
	"py": {"python3"},
	"rb": {"ruby"},
}

var extensionMatcher, _ = regexp.Compile(`\.(\w+)$`)

func detectRunner(file string) []string {
	extension := extensionMatcher.FindStringSubmatch(file)[1]
	return runners[extension]
}

func tasks(p *do.Project) {
	do.Env = `GOPATH=.vendor::$GOPATH`
	p.Task("default", do.S{"problems"}, nil)

	p.Task("problems", nil, func(c *do.Context) {
		if c.FileEvent == nil {
			return
		}
		scriptFile := c.FileEvent.Path
		ifScriptFile(scriptFile, func() {
			detectRunner(scriptFile)
			checkCommand(detectRunner(scriptFile), scriptFile)
		})
	}).Src("solutions/**/*")
}

func ifScriptFile(scriptFile string, fn func()) {
	file, err := os.Stat(scriptFile)
	if err != nil {
		return
	}
	if file.IsDir() {
		return
	}
	fn()
}

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Setenv(`KATA_PATH`, dir)
	do.Godo(tasks)
}
