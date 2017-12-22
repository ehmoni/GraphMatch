# __GraphMatch__ :  Powered by Concurrency & Parallelism
GraphMatch provides a Sub-graph Pattern Matching library for certain types of patterns. It uses naive searching algorithm to find the patterns in the Main-graph from the Sub-Graph. And also made good use of concurrency techniques using Golang's Channels and of course Parallelism in the processor. By the way ["Concurrency is not parallelism"](https://blog.golang.org/concurrency-is-not-parallelism) for those who have yet any confusion!

## **Architecture**
![Graph_Match](https://github.com/enamoni/GraphMatch/blob/master/img/GraphMatch.png)

In the very high level GraphMatch can be considered as composition of 5 modules. Let's go through them in brief:

For all modules the implementation of Graph Data Structured considered:

>A Graph is a set of Nodes and Edges

### Combination Generator: 
Combination generate all the possible combination of edges from the edge set of the graph. It also uses the channel to pass information to the goroutines for pattern matching (Pattern Finder Module). Point to be noted that, there is a single goroutine for combination generator (We can consider [Parallel Combination Generator Algorithm](http://www.sciencedirect.com/science/article/pii/0020019089901920) in the future versions.

### Pattern Finder: 
 PatternFinder takes the main graph and channels as input to receive combination data and send to Output Generator Module for JSON file creation. For Sub-Graph information it access the global Sub-Graph Profile (Created by Sub-Graph profiler function through main) as read-only mode. The naive algorithm exhaustively searches all the combinations which is by nature exponential, here we tried to see how such problems can be tackled (for which no efficient algorithm yet proposed) using concurrency and parallelism. Here we can create as many goroutines as we want to run concurrently using any number of cores. In the future versions, we can consider [DualIso: An Algorithm for Subgraph Pattern Matching](http://ieeexplore.ieee.org/document/6906821/?reload=true) for isomorphic pattern matching and of course efficiency. 

### Output Generator:
Output take the main graph and channel to receive the Output pattern for JSON file creation with the pattern marks. Also need to change the file name and path for using as a package. It can also be provided from main function as inputs. Point to be noted again, there are multiple Pattern Finder goroutines but only one Output goroutine, and as soon as they find any matched pattern, the report it to the Output module through the channel.

### Visual Graph:
This module, takes in the JSON file created by the Output module and creates [Forced Directed Graphs](https://en.wikipedia.org/wiki/Force-directed_graph_drawing) in the browser with Java Script. It creates nice visuals (mouse interactive) visual graphs inspired by physical particle simulation. This part is adapted from open source [JS library](https://gist.github.com/mbostock) by Mike Boston.

### Main and Other Programs:
Actually main is the holder of all these modules but the difference is that, when we use goroutines, the main process also runs in parallel with all other goroutines and if it's shorter and faster it may finish before other goroutines ends. So, results which we are expecting to be presented might not be visible from main, so there are few techniques (like [sync.WaitGroup](https://golang.org/pkg/sync/)) which should be used to handle such situation. About other procedures in a nut-shell:

<blockquote>
<li><b>Parallelism </b>:  Parallelism prints how many cores available in the processor and then changes the number of cores to be used based on the input provided by the user through main.
<li> <b>Random Graph Generator</b> RandomGraph generates random graph based on the input number of vertices. It uses random integer number to select 2 nodes randomly and connect them with an edge.
<li> <b>Sub-Graph Profiler</b>:  It takes input the Sub-graph and makes a profile for that in the global variable to be shared by others goroutines. The idea is that, you do not want to calculate the Sub-Graph features for every pattern finder goroutine, instead, build it once as global variable and shared by all processes as read-only mode.
<li> <b>Performance Calculator</b>:  Actually this is not an isolated function, but a section inside "main" program to calculate performance (mostly time) varying Core, Channel Buffer Size, Number of Goroutines etc. to get some insight about tuning such programs for optimum performance. 
<li><b> JSON String Producer</b>This takes input a Graph and produce string which can be used for JSON file creation. It takes into account the color and thickness of the edges (marking pattern of sub-graph in the graph) through JS Output.
<li><b>Set</b> Golang doesn't have generics and there was requirement for a Set data structure to identify unique items in a "bag" so it was adapted from this <a href="https://github.com/fatih/set">Set Library</a> by Fatih.

</blockquote>


## **Documentation**

Godoc for Golang is an excellent tool which I have used to produce the [**documentation for GraphMatch**](https://godoc.org/github.com/enamoni/GraphMatch) source code. If you use my code and want to extend it, this would be a good point to start from.

## **Walking Through GraphMatch!**
A step by step for running GraphMatch would be:

* You download everything and open the GraphMatch.go file in your favorite IDE/Editor (I have used [Goland](https://www.jetbrains.com/go/).
* Change the file path (in the main program) where you want to write the JSON files for Graph Visuals.
* Now in the main function, you can choose the number of Cores (1-4 in my case) as parameter to the Parallelism function. In the next section, you can modify number of nodes you want to for the main graph by gSize and for Sub-Graph by sgSize. These two works as parameters for RandomGraph function. Then you can fix the size of the buffer in the channel for Combination to PatternFinder by chSize and also can change the number of goroutines you want to run concurrently by grSize.
```go
	Parallelism(4)
	
	var gSize = 100
	var sgSize = 50
	var chSize = 10
	var grSize = 10
```
* After changing the variable values as desired, you can run the program. in my Goland IDE output it appeared as:

<p></p>
![Output](https://github.com/enamoni/GraphMatch/blob/master/img/Output.png)


* But of course this does not give the proper feelings, how it looks for an actual graph! To solve that, here comes the 

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone who's code was used
* Inspiration
* etc
