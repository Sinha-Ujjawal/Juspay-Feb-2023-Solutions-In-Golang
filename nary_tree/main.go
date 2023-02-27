package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"treeOfSpace/chanUtils"
	"treeOfSpace/naryTree"
	"treeOfSpace/sliceUtils"
)

type config struct {
	inputFilePath string
}

func getConfig() config {
	inputFilePath := flag.String("inputFilePath", "", "filepath of the input")
	flag.Parse()
	return config{inputFilePath: *inputFilePath}
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

type request struct {
	idx       uint
	operation naryTree.Operation
	userId    naryTree.UserId
	nodeId    naryTree.NodeId
}

func (r *request) Operation() naryTree.Operation {
	return r.operation
}

func (r *request) UserId() naryTree.UserId {
	return r.userId
}

func (r *request) NodeId() naryTree.NodeId {
	return r.nodeId
}

type response struct {
	request *request
	result  bool
}

func mkResponse(request *request, result bool) response {
	return response{request, result}
}

type treeInput struct {
	nodeIds         []string
	branchingFactor uint
	requests        []request
}

func (tInput *treeInput) processInSeq(requests []request) []response {
	if requests == nil {
		requests = tInput.requests
	}
	return chanUtils.AsSlice(
		naryTree.ProcessSeq(
			tInput.nodeIds,
			tInput.branchingFactor,
			mkResponse,
			sliceUtils.AsChannel(requests),
		),
	)
}

func (tInput *treeInput) processInPar(requests []request) []response {
	if requests == nil {
		requests = tInput.requests
	}
	return chanUtils.AsSlice(
		naryTree.ProcessPar(
			tInput.nodeIds,
			tInput.branchingFactor,
			mkResponse,
			sliceUtils.AsChannel(requests),
		),
	)
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

	requests := make([]request, numQueries)
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

		requests[q] = request{
			idx:       uint(q),
			operation: uint(opCode),
			userId:    userId,
			nodeId:    nodeId,
		}
	}
	return &treeInput{
		nodeIds:         nodeIds,
		branchingFactor: uint(branchingFactor),
		requests:        requests,
	}, nil
}

func processFile(filePath *string, processInSeq bool) ([]response, error) {
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
	responses, err := processFile(&config.inputFilePath, true)
	handleError(err)
	for _, response := range responses {
		fmt.Printf("Request: %+v, Result: %t\n", response.request, response.result)
	}
}
