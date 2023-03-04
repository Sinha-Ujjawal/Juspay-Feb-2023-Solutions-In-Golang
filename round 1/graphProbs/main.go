package main

import (
	"errors"
	"fmt"
	"graphProbs/graphProbs"
	"strconv"
	"strings"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseEdge(s string) (*graphProbs.Edge, error) {
	words := strings.Split(s, " ")
	if len(words) != 2 {
		return nil, errors.New(fmt.Sprintf("Unable to parse %s to an edge\n", s))
	}
	return &graphProbs.Edge{Frm: words[0], To: words[1]}, nil
}

func parseWeightedEdge(s string) (*graphProbs.Edge, error) {
	words := strings.Split(s, " ")
	err := errors.New(fmt.Sprintf("Unable to parse %s to an edge\n", s))
	if len(words) != 3 {
		return nil, err
	}
	wt, e := strconv.Atoi(words[2])
	if e != nil {
		return nil, err
	}
	return &graphProbs.Edge{Frm: words[0], To: words[1], Wt: uint(wt)}, nil
}

type simpleGraphInput struct {
	g         graphProbs.Graph
	follower  graphProbs.Node
	following graphProbs.Node
}

func parseSimpleGraphInput(input string) (*simpleGraphInput, error) {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil, errors.New("Inputed string is empty")
	}

	numNodes, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	lines = lines[1:]
	if len(lines) < numNodes {
		return nil, errors.New(fmt.Sprintf("Not enough nodes provided, expected: %d, found: %d\n", numNodes, len(lines)))
	}
	nodes := make([]graphProbs.Node, numNodes)
	for i := 0; i < numNodes; i += 1 {
		nodes[i] = lines[i]
	}
	lines = lines[numNodes:]

	numEdges, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	lines = lines[1:]
	if len(lines) < numEdges {
		return nil, errors.New(fmt.Sprintf("Not enough edges provided, expected: %d, found: %d\n", numEdges, len(lines)))
	}
	edges := make([]graphProbs.Edge, numEdges)
	for i := 0; i < numEdges; i += 1 {
		edge, err := parseEdge(lines[i])
		if err != nil {
			return nil, err
		}
		edges[i] = *edge
	}
	lines = lines[numEdges:]

	if len(lines) != 2 {
		return nil, errors.New(fmt.Sprintf("Unable to parse follower and following, expected 2 lines, got %d lines\n", len(lines)))
	}
	followerNode := lines[0]
	followingNode := lines[1]

	g := graphProbs.MkGraph(edges, nodes)

	return &simpleGraphInput{g: g, follower: followerNode, following: followingNode}, nil
}

func solveFindReachability(input string) {
	simpleInput, err := parseSimpleGraphInput(input)
	handleError(err)
	if simpleInput.g.CanReach(simpleInput.follower, simpleInput.following) {
		fmt.Println("1")
	} else {
		fmt.Println("0")
	}
}

type weightedGraphInput struct {
	g         graphProbs.Graph
	follower  graphProbs.Node
	following graphProbs.Node
}

func parseWeightedGraphInput(input string) (*weightedGraphInput, error) {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil, errors.New("Inputed string is empty")
	}

	numNodes, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	lines = lines[1:]
	if len(lines) < numNodes {
		return nil, errors.New(fmt.Sprintf("Not enough nodes provided, expected: %d, found: %d\n", numNodes, len(lines)))
	}
	nodes := make([]graphProbs.Node, numNodes)
	for i := 0; i < numNodes; i += 1 {
		nodes[i] = lines[i]
	}
	lines = lines[numNodes:]

	numEdges, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	lines = lines[1:]
	if len(lines) < numEdges {
		return nil, errors.New(fmt.Sprintf("Not enough edges provided, expected: %d, found: %d\n", numEdges, len(lines)))
	}
	edges := make([]graphProbs.Edge, numEdges)
	for i := 0; i < numEdges; i += 1 {
		edge, err := parseWeightedEdge(lines[i])
		if err != nil {
			return nil, err
		}
		edges[i] = *edge
	}
	lines = lines[numEdges:]

	if len(lines) != 2 {
		return nil, errors.New(fmt.Sprintf("Unable to parse follower and following, expected 2 lines, got %d lines\n", len(lines)))
	}
	followerNode := lines[0]
	followingNode := lines[1]

	g := graphProbs.MkGraph(edges, nodes)

	return &weightedGraphInput{g: g, follower: followerNode, following: followingNode}, nil
}

func solveShortestTime(input string) {
	weightedInput, err := parseWeightedGraphInput(input)
	handleError(err)
	shortestTime := weightedInput.g.ShortestTime(weightedInput.follower, weightedInput.following)
	if shortestTime == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(*shortestTime)
	}
}

func main() {
	solveFindReachability(`5
1
2
3
4
5
3
2 1
1 5
1 3
2
5`)
	solveShortestTime(`5
1
2
3
4
5
5
2 1 1
1 3 1
1 5 2
3 4 1
4 5 1
2
5`)
}
