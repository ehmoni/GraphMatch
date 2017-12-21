/*=============================================================================
 |   Assignment:  Software Design (CAS-703) - Final Project
 |
 |       Author:  Enamul Haque
 |     Language:  Go (Golang)
 |   To Compile:  GoLand (by JetBrains) is used.
 |
 |   Instructor:  Emil Sekerinski
 |  Last Update:  20 Dec 2017
 *===========================================================================*/


//Package GraphMatch provides a Sub-graph Pattern Matching library for certain types of patterns.
//It uses naive searching algorithm to find the patterns in the Main-graph. And also made good use of concurrency techniques using Golang's Channels.
package GraphMatch

import (
	"fmt"
	"math/rand"
	"encoding/json"
	"math"
	"io/ioutil"
	"sort"
	"runtime"
	"time"
	"log"
)


//Node is the basic graph vertex.
//Here Pattern is a special field used to identify in the main graph the nodes which matched with the sub-graph.
type Node struct{
	ID int
	nEdge int
	Pattern int
	Edges []*Edge
}


//Edge is the basic graph edge which have at most and at least 2 vertices.
//EdgeRank is used for sorting edges and comparision with sub-graph edges.
type Edge struct{
	ID int
	EdgeRank int
	Start *Node
	End *Node
}


//Graph is the set of nodes and edges.
type Graph struct{
	Nodes []Node
	Edges []Edge
}


//sgProfile is used to make a profile of the sub-graph which is shared with all the pattern finder functions
// to match the subset of nodes from the main graph and the sub-graph.
var sgProfile struct{
	nNodes int
	nEdges int
	sortedERank []int
	sortedNEdge []int
}


//These variables are used for performance measurement and tuning.
var (
	proc int
	allsent = false
	start time.Time
	elapCom time.Duration
	elapPf time.Duration
	elapOut time.Duration
	elapMain time.Duration
)




//main function is the entry point and also used for for performance dashboard and tuning.
func main(){

	//Defines the number of cores will be used for this program.
	Parallelism(4)

	start = time.Now()

	//Graph, Sub-Graph and Channel size with number of Goroutines for pattern matching specified.
	var gSize = 50
	var sgSize = 3
	var chSize = 100
	var grSize = 100


	//Main graph and Sub-graph generated.
	maingraph := RandomGraph(gSize)
	subgraph := RandomGraph(sgSize)


	//Channel ch declared for communicating with the pattern matching and combination generation.
	//Channel cj declared for communicating with the json Output maker.
	ch := make(chan []int,chSize)
	cj := make(chan []int)


	//Writing the JSON file both the main-graph and sub-graph
	bbb := []byte(StrJSON(maingraph))
	bb := []byte(StrJSON(subgraph))

	err1 := ioutil.WriteFile("C:\\SourceCode\\Go\\Graphdrawing\\Graph.json", bbb, 0644)
	check(err1)
	err2 := ioutil.WriteFile("C:\\SourceCode\\Go\\Graphdrawing\\SubGraph.json", bb, 0644)
	check(err2)


	//Generation of Combinations for the Pattern Finder
	edgeIDs := make([]int,len(maingraph.Edges))
	for i,e:=range maingraph.Edges{edgeIDs[i] = e.ID}

	go Combinations(edgeIDs,sgSize,ch)


	//Sub-Graph Global Information Profiler
	SgProfiler(subgraph)

	//Channeled patten searching
	for i:=0;i<grSize;i++ {
		//wg.Add(1)
		go PatternFinder(maingraph, ch, cj)
	}

	//For Output writing into json */
	go Output(*maingraph, cj)


	//Performance Measurement

	elapMain = time.Since(start)


	//Waiting for all others to finish
	var in string
	fmt.Scanln(&in)

	//Performance log printing
	log.Println("Channel Size: ", chSize)
	log.Println("Goroutines: ", proc)
	log.Println("Main : ", elapMain)
	log.Println("Combination : ", elapCom)
	log.Println("Pattern Finders : ", elapPf)
	log.Println("Output Producer : ", elapOut)
}




//Parallelism prints how many cores available in the processor and then changes the number of cores to be used based on the input for usage.
func Parallelism(n int){
	numCPU := runtime.NumCPU()
	log.Println("Available Cores: ", numCPU)

	if n > numCPU{n = numCPU}
	if n < 0 {n = 0}
	maxProcs := runtime.GOMAXPROCS(n)

	log.Println("Was Using: ", maxProcs)

	maxProcs = runtime.GOMAXPROCS(n)
	log.Println("Now Using: ", maxProcs)
}



//SgProfiler take input the Sub-graph and makes a profile for that in the global variable to be shared by others goroutines.
func SgProfiler(graph *Graph){

	sgProfile.nNodes = len(graph.Nodes)
	sgProfile.nEdges = len(graph.Edges)

	ser := make([]int, len(graph.Edges))
	for i,e:=range graph.Edges{ser[i] = e.EdgeRank}
	sort.Sort(sort.Reverse(sort.IntSlice(ser)))
	sgProfile.sortedERank = ser

	sne := make([]int, len(graph.Nodes))
	for i,n:=range graph.Nodes{sne[i] = n.nEdge}
	sort.Sort(sort.Reverse(sort.IntSlice(sne)))
	sgProfile.sortedNEdge = sne
}



//RandomGraph generates random graph based on the input number of vertices.
//It uses random integer number to select 2 nodes randomly and connect them with an edge.
func RandomGraph(size int) *Graph{

	graph := new(Graph)
	graph.Nodes = make([]Node, size)
	graph.Edges = make([]Edge, 0)


	if size<2 {return graph}

	var ri = 1
	var e  = 0


	//Random graph creation.
	for i:=0;i<size; i++{

		graph.Nodes[i].ID = i

		for {
			ri = rand.Intn(size)
			if ri != i {break}
		}

		edge := Edge{ID:e, EdgeRank: 0, Start: &graph.Nodes[i], End:&graph.Nodes[ri]}
		e++

		graph.Edges = append(graph.Edges, edge)
	}


	//Inside the nodes stuffing the edge information.
	for i:=0;i<size; i++{

		id := graph.Nodes[i].ID

		for j:=0;j<len(graph.Edges);j++{

			edge := graph.Edges[j]

			if (edge.Start.ID == id) || (edge.End.ID == id){
				graph.Nodes[i].Edges = append(graph.Nodes[i].Edges, &edge)
				graph.Nodes[i].nEdge++
			}
		}
	}

	//Calculating EdgeRank as the total number of edges the associated nodes contain.
	for i,edge:=range graph.Edges{

		graph.Edges[i].EdgeRank = edge.Start.nEdge + edge.End.nEdge
	}

	return graph
}



//combination generate all the possible combination of edges from the edge set.
//It also uses the channel to pass information to the goroutines for pattern matching.
func Combinations(iterable []int, r int, ch chan []int){

	pool := iterable
	n := len(pool)

	if r>n{return}

	indices := make([]int, r)
	for i := range indices{indices[i] = i}

	result := make([]int, r)

	for i, el := range indices{result[i] = pool[el]}

	var tmp1 = make([]int, len(result))
	copy(tmp1, result)
	ch <- tmp1

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1{}


		//As the end signal it updates the global variable and closes the channel.
		if i < 0 {
			fmt.Println("Finished Sending!")
			allsent = true
			close(ch)
			elapCom = time.Since(start)
			return
		}

		indices[i] += 1
		for j:=i+1; j<r; j+=1{indices[j] = indices[j-1] + 1}

		for;i<len(indices);i+=1{result[i] = pool[indices[i]]}

		var tmp2 = make([]int, len(result))
		copy(tmp2, result)
		ch <- tmp2
	}
}



//PatternFinder takes the main graph and channels to receive combination data and send to Output for JSON file creation.
//For Sub-Graph information it access the global Sub-Graph Profile as read-only mode.
func PatternFinder(mg *Graph, ch chan []int, cj chan []int){

	//Tracking the number of Goroutine created.
	proc++

	for {
		msg := <-ch
		if msg == nil {
			break}
		match := true


		//From the message get the edge information and build it.
		edges := make([]Edge, sgProfile.nEdges)
		for i, e := range msg {
			edges[i] = mg.Edges[e]
		}


		//Using Set to identify unique nodes for all the edges.
		allnodes := make([]Node, 0)
		s := new()
		for _, v := range edges{
			s.insert((*v.Start).ID)
			s.insert((*v.End).ID)
			allnodes = append(allnodes, *v.Start)
			allnodes = append(allnodes, *v.End)
		}


		//If the number of unique nodes matches with the number of nodes in the Sub-graph then it enters for checking.
		if s.len() == sgProfile.nNodes{

			fmt.Println("Entered!")

			//After entering based on the sorted descending order if any node has less rank/edge from the corresponding
			//nodes those are discarded and no match.
			mer := make([]int, sgProfile.nEdges)
			for i,e:=range edges{
				mer[i] = e.EdgeRank
			}
			sort.Sort(sort.Reverse(sort.IntSlice(mer)))

			nodes := make([]Node, 0)
			mne := make([]int, 0)
			for _,v:=range allnodes{
				if s.has(v.ID){
					nodes = append(nodes, v)
					mne = append(mne, v.nEdge)
					s.remove(v.ID)
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(mne)))

			for i:=0;i<sgProfile.nEdges;i++{
				if mer[i]<sgProfile.sortedERank[i]{
					match = false
					break
				}
			}

			fmt.Println(mer)

			for i:=0;i<sgProfile.nNodes;i++{
				if mne[i] < sgProfile.sortedNEdge[i]{
					match = false
					break
				}
			}

			fmt.Println(mne)
		}else{match = false}


		//If the match is found it send the information to the channel of the Output for making JSON file.
		if match{
			fmt.Println("Found a match")
			fmt.Println(msg)
			cj <- msg
		}

		if allsent{elapPf = time.Since(start)}
	}
}



//Output take the main graph and channel to receive the Output pattern for JSON file created with the pattern marks.
//Also need to change the file name and path for using as a package. It can also be provided from main inputs.
func Output(mg Graph, cj chan []int){

	for{
		msg := <-cj

		for _,e:=range msg{

			mg.Edges[e].EdgeRank = 25
			mg.Nodes[mg.Edges[e].Start.ID].Pattern = 15
			mg.Nodes[mg.Edges[e].End.ID].Pattern = 15
		}

		bb := []byte(StrJSON(&mg))

		err := ioutil.WriteFile("C:\\SourceCode\\Go\\Graphdrawing\\PattGraph.json", bb, 0644)
		check(err)
		elapOut = time.Since(start)
	}
}



//check checks for error in file handling.
func check(e error) {
	if e != nil {
		panic(e)
	}
}



//StrJSON file takes input a Graph and produce string which can be used for JSON file creation.
//It takes into account the color and thickness of the edges (marking pattern of sub-graph in the graph) through JS Output.
func StrJSON(graph *Graph) string{

	type NodesJ struct{
		Id int `json:"id"`
		Group int `json:"group"`
	}

	type LinksJ struct{
		Source int `json:"source"`
		Target int `json:"target"`
		Value int `json:"value"`
	}

	nNodes := make([]NodesJ, len(graph.Nodes))
	for i,n:=range graph.Nodes{
		nNodes[i] = NodesJ{Id:n.ID, Group:n.Pattern}
		//nNodes[i] = NodesJ{Id:n.ID, Group:n.nEdge}
	}

	nLinks := make([]LinksJ, len(graph.Edges))
	for i,n:=range graph.Edges{
		nLinks[i] = LinksJ{Source:n.Start.ID, Target:n.End.ID, Value: int(math.Ceil((float64(n.EdgeRank))/2.0))}
	}

	b, _ := json.Marshal(nNodes)
	bb, _ := json.Marshal(nLinks)

	var msg string
	msg = "{ \"nodes\": " + string(b) + "," + "\"links\": " + string(bb) + " }"

	return msg
}





//Set and nothing structures are used making Set generics using the map and interface feature.
type (
	Set struct{hash map[interface{}]nothing}
	nothing struct{}
)


// Creates a new set
func new(initial ...interface{}) *Set{

	s := &Set{make(map[interface{}]nothing)}

	for _, v:=range initial{s.insert(v)}

	return s
}

// Test to see whether or not the element is in the set
func (this *Set) has(element interface{}) bool{_, exists:=this.hash[element]
	return exists
}

// Add an element to the set
func (this *Set) insert(element interface{}){this.hash[element] = nothing{}}

// remove an element from the set
func (this *Set) remove(element interface{}){delete(this.hash, element)}

// Return the number of items in the set
func (this *Set) len() int{return len(this.hash)}
