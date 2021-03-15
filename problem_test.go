package kata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProblem(t *testing.T) {
	assert := assert.New(t)
	problem, _ := NewProblemFilePath("testdata/default.txt")
	assert.Equal("testdata/default.txt", problem.File)
	assert.Equal("testdata default", problem.Introduction())
	numberOfCases := len(problem.Cases)
	assert.Equal(1, numberOfCases)
	if numberOfCases > 0 {
		answer := problem.Cases[0]
		assert.Equal([]string{}, answer.Params)
		assert.Equal("Wooo", answer.Answer)
	}
}

func TestSimpleProblem(t *testing.T) {
	assert := assert.New(t)
	problem, _ := NewProblemFilePath("testdata/simple.txt")
	// assert.Equal("This is my answer", problem.Answer)
	assert.Equal("testdata/simple.txt", problem.File)
	expectedIntro := `This is a statement of the problem. It can span multiple lines for a paragraph.

It can have a paragraph by having 2 newlines between.

Here is some indented text which preserves newlines:

  The quick brown fox
  Jumped over the lazy dog
    Poor sod.`

	assert.Equal(expectedIntro, problem.Introduction())

	numberOfCases := len(problem.Cases)
	assert.Equal(1, numberOfCases)
	if numberOfCases > 0 {
		answer := problem.Cases[0]
		assert.Equal([]string{}, answer.Params)
		assert.Equal("This is my answer", answer.Answer)
	}
}

func TestMultiCaseProblem(t *testing.T) {
	assert := assert.New(t)
	problem, _ := NewProblemFilePath("testdata/cases.txt")
	assert.Equal("testdata/cases.txt", problem.File)
	expectedIntro := `This demonstrates string case and answers.`
	assert.Equal(expectedIntro, problem.Introduction())
	assert.Equal("", problem.Answer)

	numberOfCases := len(problem.Cases)
	assert.Equal(3, numberOfCases)
	if numberOfCases == 3 {
		answer := problem.Cases[0]
		assert.Equal([]string{`1`, `2`}, answer.Params)
		assert.Equal("obviously julie", answer.Answer)

		answer = problem.Cases[1]
		assert.Equal([]string{"the", "quick", "brown", "fox"}, answer.Params)
		assert.Equal("foo\nbar", answer.Answer)

		answer = problem.Cases[2]
		assert.Equal([]string{"the quick", "brown", "fox"}, answer.Params)
		assert.Equal("bar\nbaz", answer.Answer)
	}
}
