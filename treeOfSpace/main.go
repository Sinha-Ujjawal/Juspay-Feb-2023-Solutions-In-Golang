package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"treeOfSpace/naryTree"
	"treeOfSpace/sliceUtils"
)

type config struct {
	inputFilePath string
	runInParallel bool
}

func getConfig() config {
	inputFilePath := flag.String("inputFilePath", "", "filepath of the input")
	runInParallel := flag.Bool("runInParallel", false, "if true, the requests will be submitted in parallel, hence can be out of order")
	flag.Parse()
	return config{inputFilePath: *inputFilePath, runInParallel: *runInParallel}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err == nil {
		return strings.TrimSpace(line), nil
	}
	return line, err
}

func readLines(reader *bufio.Reader, n uint) ([]string, error) {
	lines := make([]string, n)
	for i := uint(0); i < n; i++ {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		lines[i] = line
	}
	return lines, nil
}

func readInt(reader *bufio.Reader) (int64, error) {
	line, err := readLine(reader)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(line, 10, 64)
}

type treeInput struct {
	nodeIds         []string
	branchingFactor uint
	requests        []naryTree.Request
}

func (tInput *treeInput) processInSeq(requests []naryTree.Request) []naryTree.Response {
	if requests == nil {
		requests = tInput.requests
	}
	requestCh := sliceUtils.AsChannel(requests)
	responseCh := naryTree.ProcessSeq(
		tInput.nodeIds,
		tInput.branchingFactor,
		requestCh,
	)
	responses := make([]naryTree.Response, len(tInput.requests))
	i := 0
	for response := range responseCh {
		responses[i] = response
		i += 1
	}
	return responses
}

func (tInput *treeInput) processInPar(requests []naryTree.Request) []naryTree.Response {
	if requests == nil {
		requests = tInput.requests
	}
	requestCh := sliceUtils.AsChannel(requests)
	responseCh := naryTree.ProcessPar(
		tInput.nodeIds,
		tInput.branchingFactor,
		requestCh,
	)
	responses := make([]naryTree.Response, len(tInput.requests))
	i := 0
	for response := range responseCh {
		responses[i] = response
		i += 1
	}
	return responses
}

func parseFile(filePath *string) (*treeInput, error) {
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	numNodes, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	branchingFactor, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	numQueries, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	nodeIds, err := readLines(reader, uint(numNodes))
	if err != nil {
		return nil, err
	}

	requests := make([]naryTree.Request, numQueries)
	for q := int64(0); q < numQueries; q++ {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}

		parts := strings.Split(line, " ")

		opCode, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		nodeId := parts[1]

		userId, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, err
		}

		requests[q] = naryTree.Request{
			Operation: uint(opCode),
			UserId:    userId,
			NodeId:    nodeId,
		}
	}
	return &treeInput{
		nodeIds:         nodeIds,
		branchingFactor: uint(branchingFactor),
		requests:        requests,
	}, nil
}

func processFile(filePath *string, processInSeq bool) ([]naryTree.Response, error) {
	tInput, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}
	if processInSeq {
		return tInput.processInSeq(nil), nil
	}
	return tInput.processInPar(nil), nil
}

func main() {
	config := getConfig()
	fmt.Printf("Processing file: %s\n", config.inputFilePath)
	if config.runInParallel {
		fmt.Println("Processing concurrently")
	}
	responses, err := processFile(&config.inputFilePath, !config.runInParallel)
	handleError(err)
	for _, response := range responses {
		fmt.Printf("Request: %+v, Result: %t\n", response.Request, response.Result)
	}
}
