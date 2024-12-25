package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Node struct {
	ID          int
	FingerTable []Finger
	Successor   *Node
	Predecessor *Node
}

type Finger struct {
	Start int
	Node  *Node
}

const (
	M = 10 
)

func NewNode(id int) *Node {
	node := &Node{
		ID:          id,
		FingerTable: make([]Finger, M),
	}
	for i := range node.FingerTable {
		node.FingerTable[i].Start = (id + int(math.Pow(2, float64(i)))) % int(math.Pow(2, float64(M)))
	}
	return node
}

func (n *Node) UpdateFingers() {
	for i := range n.FingerTable {
		successor := n.FindSuccessor(n.FingerTable[i].Start)
		n.FingerTable[i].Node = successor
	}
}

func (n *Node) FindSuccessor(id int) *Node {
	if id < 0 {
		id = id + int(math.Pow(2, float64(M)))
	}
	
	if between(n.ID, id, n.Successor.ID, true) {
		return n.Successor
	}

	// Otherwise, find the closest preceding node and ask it
	nprime := n.ClosestPrecedingNode(id)
	if nprime == n {
		return n.Successor
	}
	return nprime.FindSuccessor(id)
}

// Helper function to check if x is between a and b on the ring
func between(a, x, b int, inclusive bool) bool {
	if a < b {
		return (inclusive && a <= x && x <= b) || (!inclusive && a < x && x < b)
	}
	// Handle wrap-around
	return (inclusive && (x >= a || x <= b)) || (!inclusive && (x > a || x < b))
}

func (n *Node) ClosestPrecedingNode(id int) *Node {
	for i := M - 1; i >= 0; i-- {
		if n.FingerTable[i].Node != nil {
			fingerID := n.FingerTable[i].Node.ID
			if between(n.ID, fingerID, id, false) {
				return n.FingerTable[i].Node
			}
		}
	}
	return n
}

func (n *Node) Join(node *Node) {
	if node != nil {
		n.InitFingerTable(node)
		n.UpdateOthers()
	} else {
		n.Predecessor = n
		n.Successor = n
		for i := range n.FingerTable {
			n.FingerTable[i].Node = n
		}
	}
}

func (n *Node) InitFingerTable(node *Node) {
	// Find successor for first finger
	n.FingerTable[0].Node = node.FindSuccessor(n.FingerTable[0].Start)
	n.Successor = n.FingerTable[0].Node
	n.Predecessor = n.Successor.Predecessor
	n.Successor.Predecessor = n
	for i := 0; i < M-1; i++ {
		if between(n.ID, n.FingerTable[i+1].Start, n.FingerTable[i].Node.ID, true) {
			n.FingerTable[i+1].Node = n.FingerTable[i].Node
		} else {
			n.FingerTable[i+1].Node = node.FindSuccessor(n.FingerTable[i+1].Start)
		}
	}
}

func (n *Node) UpdateOthers() {
	for i := 0; i < M; i++ {
		// Find last node p whose i-th finger might be n
		power := int(math.Pow(2, float64(i)))
		id := (n.ID - power + int(math.Pow(2, float64(M)))) % int(math.Pow(2, float64(M)))
		p := n.FindPredecessor(id)
		p.UpdateFingerTable(n, i)
	}
}

func (n *Node) UpdateFingerTable(s *Node, i int) {
	if between(n.ID, s.ID, n.FingerTable[i].Node.ID, false) {
		n.FingerTable[i].Node = s
		p := n.Predecessor
		p.UpdateFingerTable(s, i)
	}
}

func (n *Node) FindPredecessor(id int) *Node {
	curr := n
	for !between(curr.ID, id, curr.Successor.ID, true) {
		curr = curr.ClosestPrecedingNode(id)
	}
	return curr
}

func main() {
	file, err := os.Open("nodes.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var nodes []*Node
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Error converting ID: %v\n", err)
			continue
		}
		nodes = append(nodes, NewNode(id))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if len(nodes) == 0 {
		fmt.Println("No nodes found in file")
		return
	}

	nodes[0].Join(nil) 
	for i := 1; i < len(nodes); i++ {
		nodes[i].Join(nodes[0])
	}

	for _, node := range nodes {
		node.UpdateFingers()
	}

	key := 200
	successor := nodes[0].FindSuccessor(key)
	fmt.Printf("Successor of key %d is node %d\n", key, successor.ID)
	for _, node := range nodes {
		fmt.Printf("\nNode %d finger table:\n", node.ID)
		for _, finger := range node.FingerTable {
			fmt.Printf("Start: %d, Node: %d\n", finger.Start, finger.Node.ID)
		}
	}
}