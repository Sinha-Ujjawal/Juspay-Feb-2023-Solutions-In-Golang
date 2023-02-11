package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"treeOfSpace/naryTree"
)

type config struct {
	inputFilePath          string
	expectedOutputFilePath string
}

func getConfig() config {
	inputFilePath := flag.String("inputFilePath", "", "filepath of the input")
	expectedOutputFilePath := flag.String("expectedOutputFilePath", "", "filepath of the expected output")
	flag.Parse()
	return config{
		inputFilePath:          *inputFilePath,
		expectedOutputFilePath: *expectedOutputFilePath,
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	handleError(err)
	return strings.TrimSpace(line)
}

func readLines(reader *bufio.Reader, n uint) []string {
	var lines []string
	for i := uint(0); i < n; i++ {
		line := readLine(reader)
		lines = append(lines, line)
	}
	return lines
}

func readInt(reader *bufio.Reader) int64 {
	line := readLine(reader)
	value, err := strconv.ParseInt(line, 10, 64)
	handleError(err)
	return value
}

func processFile(filePath *string) <-chan bool {
	file, err := os.Open(*filePath)
	handleError(err)
	stdInReader := bufio.NewReader(file)

	numNodes := uint(readInt(stdInReader))
	branchingFactor := uint(readInt(stdInReader))
	numQueries := uint(readInt(stdInReader))
	nodes := readLines(stdInReader, numNodes)
	tree := naryTree.New(nodes, branchingFactor)

	in := make(chan string, numQueries)
	out := make(chan bool, numQueries)

	go func() {
		for q := uint(0); q < numQueries; q++ {
			line := readLine(stdInReader)
			in <- line
		}
		close(in)
	}()

	go func() {
		for s := range in {
			line := strings.Split(s, " ")
			opCode := line[0]
			nodeId := line[1]
			userId, err := strconv.ParseInt(line[2], 10, 64)
			handleError(err)
			var b bool
			switch opCode {
			case "1":
				b = tree.Lock(nodeId, userId)
			case "2":
				b = tree.Unlock(nodeId, userId)
			case "3":
				b = tree.Upgrade(nodeId, userId)
			default:
				panic(fmt.Sprintf("Did not recognize the opCode: %s", opCode))
			}
			out <- b
		}
		close(out)
	}()

	return out
}

func readExpectedOutput(filePath *string) <-chan bool {
	file, err := os.Open(*filePath)
	handleError(err)
	stdInReader := bufio.NewReader(file)
	out := make(chan bool)
	go func() {
		for {
			line, err := stdInReader.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			switch line {
			case "true":
				out <- true
			case "false":
				out <- false
			default:
				panic(fmt.Sprintf("Did not recognize the boolean string: %s", line))
			}
		}
		close(out)
	}()
	return out
}

func main() {
	config := getConfig()
	fmt.Printf("Processing file: %s\n", config.inputFilePath)
	fmt.Printf("Expected output file: %s\n", config.expectedOutputFilePath)
	outputCh := processFile(&config.inputFilePath)
	expectedOutCh := readExpectedOutput(&config.expectedOutputFilePath)
	for result := range outputCh {
		expectedResult := <-expectedOutCh
		if result != expectedResult {
			println("Test Failed!")
			os.Exit(1)
		}
	}
	println("Test Success!")
}
