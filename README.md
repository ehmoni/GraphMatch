# __GraphMatch__ :  Powered by Concurrency & Parallelism
GraphMatch provides a Sub-graph Pattern Matching library for certain types of patterns. It uses naive searching algorithm to find the patterns in the Main-graph from the Sub-Graph. And also made good use of concurrency techniques using Golang's Channels and of course Parallelism in the processor. By the way ["Concurrency is not parallelism"](https://blog.golang.org/concurrency-is-not-parallelism) for those who have yet any confusion!

## **Architecture**
![Graph_Match](https://github.com/enamoni/GraphMatch/blob/master/img/GraphMatch.png)

In the very high level GraphMatch can be considered as composition of 5 modules. Let's go through them in brief:

For all modules the implementation of Graph Data Structure we defined:

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

Godoc for Golang is an excellent tool which I have used to produce the [***documentation for GraphMatch***](https://godoc.org/github.com/enamoni/GraphMatch) source code. If you use my code and want to extend it, this would be a good point to start from.

## **Walking Through GraphMatch!**
A step by step for running GraphMatch would be:

* You download everything and open the GraphMatch.go file in your favorite IDE/Editor (I have used [Goland](https://www.jetbrains.com/go/).
* Change the file path (in the main program) where you want to write the JSON files for Graph Visuals.
* Now in the main function, you can choose the number of Cores (1-4 in my case) as parameter to the Parallelism function. In the next section, you can modify number of nodes you want to for the main graph by gSize and for Sub-Graph by sgSize. These two works as parameters for RandomGraph function. Then you can fix the size of the buffer in the channel for Combination to PatternFinder by chSize and also can change the number of goroutines you want to run concurrently by grSize.
```go
	Parallelism(4)
	
	var gSize = 500
	var sgSize = 30
	var chSize = 100
	var grSize = 100
```
* After changing the variable values as desired, you can run the program. in my Goland IDE output it appeared as:

<p></p>

![Output](https://github.com/enamoni/GraphMatch/blob/master/img/Output.png)


* But of course this does not give the proper feelings, how it looks for an actual graph! To solve that, here comes the Forced Directed Graph library. 
* Open the folder where the JS files and the htmls are located Open the index.html, indexSub.html and Pattern.html, you will see outputs like below:

<p></p>

![Out1](https://github.com/enamoni/GraphMatch/blob/master/img/Out1.PNG)
![Out2](https://github.com/enamoni/GraphMatch/blob/master/img/Out2.PNG)
![Out3](https://github.com/enamoni/GraphMatch/blob/master/img/Out3.PNG)

### Interactive Output:
Yes you can see them in action through the following links:


>[**Main Graph**](http://rawgit.com/enamoni/GraphMatch/master/VisualOutput/index.html) &nbsp;&nbsp;&nbsp; || &nbsp;&nbsp;
>[**Sub-Graph**   ](http://rawgit.com/enamoni/GraphMatch/master/VisualOutput/indexSub.html) &nbsp;&nbsp;&nbsp;|| &nbsp;&nbsp;
>[**Main Graph with the Pattern**](http://rawgit.com/enamoni/GraphMatch/master/VisualOutput/Pattern.html)



## **Insights**
Concurrency and Parallelism are not like straight forward programs. Here, what we might think is going to improve the performance, but due to lack of proper organization of statements, design and tuning it might result the opposite of what was expected before. So, there are some general principles which should be considered while designing as well as there are some special cases specific to the program which cannot be ignored. For this case, we checked some parameters and tried to figure out what works and what not. Just like any good experiment, while we are observing one parameter impact with the change of another parameter, we try to keep other parameters and environment variables fixed. So, for the below results and insights, we kept the main graph and sub-graph same for all of them. And the scope of the parameters consists of:

>1. Running Time of Main function (in ms)
>2. Running Time of Combination & (up to last) Pattern Finder goroutines (in s) 
>3. Number of Cores
>4. Number of Goroutines (As Pattern Finder)
>5. Size of the Channel (Combination => Pattern Finder) 
 
<p></p>



![Output](https://github.com/enamoni/GraphMatch/blob/master/img/Insight.PNG)


<blockquote>
<li> The plots above are the results of our measurement of (1,2) by changing (3,4,5) of the parameter list items mentioned above. Here the values are taken as average of 3 readings i.e. one configuration was run 3 times and taken their average timing.
<li> <b># Core Increase: </b> First we changed the number of Cores to get the essence of parallelism at processor level. Interesting findings is that, when we increased the number of cores from 1-2-3-4 the running time for the goroutines (orange bars) decreased significantly but the running time of the main function (blue bars) are not affected much. The reason could be main function here always ran in one core it was not shared among multiple cores whereas the goroutines must have been distributed among the 4 cores.
<li> <b># Goroutine Increase (in 1 Core):</b> As the number of goroutines increased altogether they took little bit more time but notable change was observed in increasing time of the main function. It is expected as in the one core now all the goroutines are working so main will be slower to finish.
<li> <b># Goroutine Increase: </b> It was done with all 4 cores and the interesting observation is that, as we increased the number of (Pattern Finder) goroutines it did not vary much after certain number. It is better of course than only 1 or few goroutines but after certain points there seems no impact of their numbers. To conclude firmly it has to be experimented more with varied input output ranges that should there be any best number of goroutines for any particular program. For our case with the inputs given, it seems like 20-30 would be the number of goroutines to run it efficiently.
<li> <b> # Channel Increase: </b> It is actually not the number of channels but the size of the channel i.e. the buffer size. I tried with 20 then increased to 100 and 500 other than the running time of the main function which could be of some other effect, I did not find any significant change in the time of the go routines. But when tried with synchronous channel (i.e. removed the buffer) then found a lag of time between the combination ending time and the pattern finder go routines. For the all other experiments their ending time was same (as the channel was closed as soon as the combination stops generating combinations) and I kept them together (orange bars). But overall, I found that asynchronous communication channels are bit faster than the synchronous approach for this case, as there are many goroutines (Pattern Finders) looking for the combination from the channel and as their number increases blocking the channel by a single goroutine for synchronous communication delays other access which might be ready to receive by then, besides it will also make slower the fast generation of the combinations send to the channel by the single combination generator module. 
</blockquote>

This GraphMatch is an attempt to tackle Combinatoric Problems from concurrency and parallelism point of view. There are many areas where this prototype can be improved and lot of research opportunities that can be explored. It is an open source program, so inviting all enthusiasts to use and contribute to this effort.

### **License**

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/) for details

## Acknowledgments

* Special Thanks to  [Prof Emil Sekerinski](http://www.cas.mcmaster.ca/~emil/Welcome.html) for introducing us to the amazing world of concurrency (and golang) and for all the inspiration, encouragement and guidance throughout the learning process!
