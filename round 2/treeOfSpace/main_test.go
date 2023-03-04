package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	"treeOfSpace/naryTree"
	"treeOfSpace/sliceUtils"
)

const testCasesDir string = "./test_cases"

func readOutputFile(filePath string) ([]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	ret := []bool{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch strings.TrimSpace(line) {
		case "true":
			ret = append(ret, true)
		case "false":
			ret = append(ret, false)
		default:
			return nil, errors.New(fmt.Sprintf("Something went wrong, the file contained unrecognised line: %s\n", line))
		}
	}
	return ret, nil
}

func TestProcessingInSeq(t *testing.T) {
	testCases, err := os.ReadDir(testCasesDir)
	if err != nil {
		t.Error(err)
		return
	}

	sort.SliceStable(testCases, func(i, j int) bool {
		return testCases[i].Name() < testCases[j].Name()
	})

	for _, testCase := range testCases {
		if !strings.HasPrefix(testCase.Name(), "test_case") {
			continue
		}
		inputFilePath := path.Join(testCasesDir, testCase.Name(), "input.txt")
		outputFilePath := path.Join(testCasesDir, testCase.Name(), "output.txt")
		outputResponses, err := processFile(&inputFilePath, true)
		if err != nil {
			t.Error(err)
			return
		}
		output := sliceUtils.Map(
			outputResponses,
			func(res naryTree.Response) bool { return res.Result },
		)
		expectedOutput, err := readOutputFile(outputFilePath)
		if err != nil {
			handleError(err)
			t.Fail()
		}
		if !sliceUtils.Equal(output, expectedOutput) {
			fmt.Printf("Failed on test case: inputFilePath: %s, outputFilePath: %s\n", inputFilePath, outputFilePath)
			fmt.Printf("Output: %#v\n", output)
			fmt.Printf("Expected Output: %#v\n", expectedOutput)
			t.Fail()
			break
		}
	}
}

func TestProcessingInPar(t *testing.T) {
	testCases, err := os.ReadDir(testCasesDir)
	if err != nil {
		t.Error(err)
		return
	}

	sort.SliceStable(testCases, func(i, j int) bool {
		return testCases[i].Name() < testCases[j].Name()
	})

	for _, testCase := range testCases {
		if !strings.HasPrefix(testCase.Name(), "test_case") {
			continue
		}
		inputFilePath := path.Join(testCasesDir, testCase.Name(), "input.txt")
		tInput, err := parseFile(&inputFilePath)
		if err != nil {
			t.Error(err)
			return
		}

		outputResponses := tInput.processInPar(nil)

		requestsSeq := sliceUtils.Map(
			outputResponses,
			func(res naryTree.Response) naryTree.Request { return res.Request },
		)
		output := sliceUtils.Map(
			outputResponses,
			func(res naryTree.Response) bool { return res.Result },
		)

		expectedOutput := sliceUtils.Map(
			tInput.processInSeq(requestsSeq),
			func(res naryTree.Response) bool { return res.Result },
		)

		if !sliceUtils.Equal(output, expectedOutput) {
			fmt.Printf("Failed on test case: inputFilePath: %s\n", inputFilePath)
			fmt.Printf("Output: %#v\n", output)
			fmt.Printf("Expected Output: %#v\n", expectedOutput)
			t.Fail()
			break
		}
	}
}
