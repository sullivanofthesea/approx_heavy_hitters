//inputs file, estimates the top 10 most frequently accessed paths, outputs percentils of file sizes seen in universe
//preliminary submission by Acamar Orionis (Erica Stephens)
//Credit to: https://github.com/shenwei356/countminsketch for CMS data structure
//		 to: GoDS - Go Data Structures for Tree and Arraylist implementations

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"math"
	//library for countmin data structure for AHH
	"github.com/shenwei356/countminsketch"
	//library for tree solution for AHH, Percentile Trees
	myTrees "github.com/emirpasic/gods/maps/treemap"	
	//library for final AHH
	"github.com/emirpasic/gods/lists/arraylist"
)
/* global variable declaration */
//TODO Reconstruct to avoid global variables & encapsulate in privately accessed class
var pathcms string
var pathsize int

func main() {

	err := Main()
	if err != nil {
		panic(err)
	}
}

func Main() error {

	f, err := os.Open("path1.txt")
		if err != nil {
		    panic(err)
		}
		defer f.Close()
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = f

	return ReadInput(os.Stdin, func(path string, size int64) error {
		_, err :=  fmt.Printf("") 
		//_, err :=  fmt.Printf("path: %s (size: %d)\n", path, size) //DEBUG uncomment if all paths/sizes need to be output on console
		return err
	})
}

// ReadInput calls cb with every new path and size read from the input
// TODO this function contains several "tarball" subfunctions that should be broken out into their own
//      i.e  initializeCMS(), updateCMS(), initializeAHHTree(), updateAHHTree(), initializePercentileTree(), updatePercentileTree(), printAHH()
func ReadInput(r io.Reader, cb func(path string, size int64) error) error {
	process := func(line string) error {
		
		fields := strings.Split(strings.TrimRight(line, "\r\n"), "\t")
		//TODO ENCAPSULATE IN PRIVATE CLASS
		pathcms = fields[0]

		if len(fields) != 2 {
			return fmt.Errorf("line %q: malformed", line)
		}
		size, err := strconv.ParseInt(fields[1], 10, 64)

		//TODO cleanup, privatize
		pathsize,_ =strconv.Atoi(fields[1])

		if err != nil {
			fmt.Printf("ERROR ELS- 1")
			return fmt.Errorf("line %q: %v", line, err)
		}
		return cb(fields[0], size)
	}
	br := bufio.NewReader(r)
	
	//CREATE INITIAL SKETCH- TODO separate function initializeCMS()
		var varepsilon, delta float64

		//FYI - 0.001, 0.0009 provides more accuracy, but also costs more space and processing time
		//TODO pull this varepsilon, delta values out into an external config file
		varepsilon, delta = 0.01, 0.9
		s, err := countminsketch.NewWithEstimates(varepsilon, delta)

		fmt.Printf("ε: %f, δ: %f -> d (#hashing functions): %d, w(size of every hash table): %d\n", varepsilon, delta, s.D(), s.W())
		//FYI - The error in responding to a sketch query is within a factor of ε with probability δ
		
		//TODO, JSON kept for debugging sketch, could be made into a viewSketch() function
		key := "placeholdercreationpath.txt"
		seedVal := "seedj1k2l3h4g5f6"
		s.UpdateString(key, 1)
		bytes, err := s.MarshalJSON()
		checkerr(err)
		err = s.UnmarshalJSON(bytes)
		checkerr(err)			
	//COUNT MIN FIRST SKETCH END

	//CREATE INITIAL AHH TREE MAP, Seeding it, TODO create separate function initializeAHHTree()
		m := myTrees.NewWithIntComparator()  
		for seed := 1; seed <= 12; seed++ {
				m.Put(seed, seedVal) 
			}
	//END INITIAL AHH TREE

	//CREATE INITIAL file size percentile tree, TODO create separate function initializePercentileTree()
		p := myTrees.NewWithIntComparator() //  TODO update library to handle uint64
	//END INITIAL file size percentile tree

	//READS in file line by line
	//CREATES new CMS data structure for each line, MERGES with previous CMS
	//UPDATES AHH treemap and prunes it to keep only top relevant 
	//UPDATES Percentile treemap to include size values
	//OUTPUTS final top AHH, percentile
	//TODO break these out into separate functions
	for {
		l, err := br.ReadString('\n')
	
			if err != nil {
				if err == io.EOF {
					if l != "" {
						err = process(l)

						if err != nil {
							return err
						}
					}
					//IF WE ARE AT THE END OF THE FILE, THESE ARE THE VERY LAST STEPS
					
					fmt.Printf( "\nTop 10 Paths:\n")
					// printTree(m) //TODO function to print tree for debugging
					
					//TODO Format top AHH, separate function printAHH()
					//grab greatest tree node (aka highest CMS value, parses to get top 10)
					//Consider efficiency of this as well
					itAHH := m.Iterator()
					itAHH.End()
					finalAHH := arraylist.New()

					for j := 1; j< m.Size() && finalAHH.Size() <10 ; j++ {
						itAHH.Prev()
						_, pathChainVal := itAHH.Key(), itAHH.Value()
						if pathChainVal != seedVal {
							//parse by |, loop through separate chained paths, add if doesnt already exist, not seedVal
							s := strings.Split(pathChainVal.(string), "|")	
						    for x, _:= range s {
						        if s[x] != seedVal && finalAHH.Contains(s[x])==false  && finalAHH.Size() <10 {
						        	finalAHH.Add(s[x])
						        	//here is where we print the top 10 AHH from highest to lowest
						        	fmt.Printf("%v. %v\n",finalAHH.Size(), s[x])
						        }
						    }
						}
					}

					//kth Percentile (.k)* # items = index in list (round up)
					//do this every batch N value (500 x)
					var p50f, p75f, p90f, p99f = 0.5, 0.75, 0.90, 0.99
					var p50key, p75key, p90key, p99key = 0, 0, 0, 0

					p50:=int(math.Ceil(p50f*float64(p.Size())) )
					p75:=int(math.Ceil(p75f*float64(p.Size())) )
					p90:=int(math.Ceil(p90f*float64(p.Size())) )
					p99:=int(math.Ceil(p99f*float64(p.Size())) )

					it := p.Iterator()
					    it.First()
					    szTotal:=0
					    //this loop obtains all the size values for this batch and accumulates a total
					    for i := 1; i< p.Size() ; i++ {
							key, _ := it.Key(), it.Value()

							//Save the tree keys that house the kth percentile values
							switch {
							    case i == p50:
							        p50key=key.(int)
							    case i == p75:
							        p75key=key.(int)
							    case i == p90:
							        p90key=key.(int)
							    case i == p99:
							        p99key=key.(int)
							        break								        
						    }
						    it.Next()
						    szTotal+=int(key.(int))
						}
					//PRINT PERCENTILES
					fmt.Printf( "\nPercentiles:\n")
					fmt.Printf("file_size_p50\t%v\t%v\t\n", p50key, szTotal )
					fmt.Printf("file_size_p75\t%v\t%v\t\n", p75key, szTotal )
					fmt.Printf("file_size_p90\t%v\t%v\t\n", p90key, szTotal )
					fmt.Printf("file_size_p99\t%v\t%v\t\n", p99key, szTotal )
					//EOF- THIS IS TRUE END TO READIN, REMAINING CODE IS LOOPING THROUGH FILE RECORD 1 through EOF-1
					return nil
				}
				return err
			}
			err = process(l)

		// CMS CREATE NEXT SKETCH TODO break into own function
			keynext := pathcms
			sketchnext, err := countminsketch.NewWithEstimates(varepsilon, delta)
			sketchnext.UpdateString(keynext, 1)
			//JSON for debugging
			bytesnext, err := sketchnext.MarshalJSON()
			checkerr(err)
			err = sketchnext.UnmarshalJSON(bytesnext)
			checkerr(err)

			//MERGE CMS SKETCH - TODO break into own function updateCMS()
			s.Merge(sketchnext)
			//JSON for debugging
			bytesmerge, err := s.MarshalJSON()
			checkerr(err)
			err = s.UnmarshalJSON(bytesmerge)
			checkerr(err)
		// END CMS NEXT SKETCH/MERGE

		//UPDATE TREEMAP/AHH - TODO break into separate function updateAHHTree()
			minCMS, _ := m.Min()
			thisCMS:= s.EstimateString(keynext)
			//TODO Cleanup conversions/data types here
			sminCMSTree := minCMS.(int)
			sthisCMS := strconv.FormatUint(thisCMS, 10)
			fileStreamCMS, _ := strconv.Atoi(sthisCMS)
			treePath, cmsAlreadyInTree := m.Get(fileStreamCMS)  

    		//AHH algorithm rules:
			//CASE1: if stream bigger than smallest tree value, and stream CMS DOES NOT already exist in tree
			if fileStreamCMS > sminCMSTree && cmsAlreadyInTree == false {
					//simply add the filestream value and path to the tree
					m.Put(fileStreamCMS, pathcms) 
			}
			//CASE2, 3: if CMS stream is equal to the minimum tree value (it's already in the tree) 
			// 			OR if CMS stream is greater but it's already somewhere else in the tree
			if fileStreamCMS == sminCMSTree || (fileStreamCMS > sminCMSTree && cmsAlreadyInTree == true) {				
				//if we already have the path, then just make sure there are no seeds left
				//if we DONT already have the path, remove seeds, append path as SEPARATE CHAINING to resolve CMS collisions
				if strings.Contains(treePath.(string), pathcms) {
					//remove seed values once you get something better
					treePath = strings.ReplaceAll(treePath.(string), seedVal, "") 
					//overide value in tree with new concatenated path
					m.Put(fileStreamCMS, treePath.(string)) 
				} else  {
					//remove seed values once you get something better
					treePath = strings.ReplaceAll(treePath.(string), seedVal, "") 
					newTreePath:=""

					//if stream path is not in treepath, append stream path
					if treePath.(string) == "" {
						newTreePath = pathcms	
					} else {
						newTreePath =  pathcms + "|"  + treePath.(string)
					}	
					//overide value in tree with new concatenated path
					m.Put(fileStreamCMS, newTreePath) 
				}
				    //TODO separate pruneAHHTree() function, maybe reduce to top 15 in separate config value
				    if m.Size()>30 {
				    m.Remove(sminCMSTree)
				    }
			}
		//END TREEMAP/AHH
			p.Put(pathsize,pathcms)
			if err != nil {
			return err
			}
		}//END READIN LOOP
}

//TODO merge with pre-existing error handling
func checkerr(err error) {
if err != nil {
fmt.Println(err)
os.Exit(1)
}
}
//TODO 	Break into smaller functions, switch back from file to stream, privatize global variables, 
//TODO  Create testing package, test all corner cases, create config file for settings
//TODO  Requirement #3, batch and memory
//TODO  Requirement #4, open and print to socket, currently printing to console
