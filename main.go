/*
mohsen-classagreement generates class agreement tables similar to those found in the following paper:

        @inproceedings{Moffat:2012:MMI:2407085.2407092,
         author = {Moffat, Alistair and Scholer, Falk and Thomas, Paul},
         title = {Models and metrics: IR evaluation as a user process},
         booktitle = {Proceedings of the Seventeenth Australasian Document Computing Symposium},
         series = {ADCS '12},
         year = {2012},
         isbn = {978-1-4503-1411-4},
         location = {Dunedin, New Zealand},
         pages = {47--54},
         numpages = {8},
         url = {http://doi.acm.org/10.1145/2407085.2407092},
         doi = {10.1145/2407085.2407092},
         acmid = {2407092},
         publisher = {ACM},
         address = {New York, NY, USA},
         keywords = {evaluation, retrieval experiment, system measurement},
        } 

The tables differ in that instead of putting the SSA and NN tables straight in to the two diagonal halves
of one table, two tables are produced. This approach introduces a little bit of duplication in the table, 
but is slightly easier to read, as you can read down one column or across one row to see the information for
one class. Additionally, the output table is sorted.

It takes as input the tables produced by Mohsen Laali's code at https://github.com/Mohsen-Laali/TREC_COD ,
with a different header line.  The format of the input table is:

      <header line>
      <data line 1>
      ....
      <data line 120>

Where the header line contains S sort keys as floating point numbers, separated by commas. These sort keys are used to sort the output table.

Data lines are of the form:

      <test name>,<System A better than B>,<System B better than A>,,,, 

with enough empty columns to bring the file up to S entries.

The <System A better than B> entries should be "yes" if System A is statistically significantly better than System B, and "no" otherwise.
Similarly, the <System B better than A> entries should be "yes" if System B is statistically significantly better than System A, and "no" otherwise.

Test name entries are ignored by this program.

Usage:
     
     mohsen-classagreement <DIR> [index]

DIR is a path to a directory containing several csv files in the format described above, and index is the index in the header line to sort the results by (0 if unspecified).

All input files must be the same length, and this is confirmed at input time.

*/
package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "path"
  "sort"
  "log"
  "strconv"
)


func main() {
  if len(os.Args) < 2 || len(os.Args) > 3 {
    fmt.Println(os.Args[0], ": calculate class agreement between multiple subcollections")
    fmt.Println("Usage:")
    fmt.Println("   ",os.Args[0], " <path to directory containing subcollection significance csv files> [index for the sort key]")
    return
  }

  // Read the collection CSVs
  dirName :=  os.Args[1]
  var sortkey int64
  var err error 
  if len(os.Args) > 2 {
    sortkey, err = strconv.ParseInt(os.Args[2],10,0)
    if err != nil {
      log.Fatal(err)
    }
  }
  collectionName := path.Base(dirName)
  fileListing, err:= ioutil.ReadDir(dirName);
  if err != nil {
    log.Fatal(err)
  }
  cols := make([]*SubcollectionRun,0,len(fileListing))
  cols_ndcg := make([]*SubcollectionRun,0,len(fileListing))
  for _,file := range fileListing {
     if path.Ext(file.Name()) == ".csv"{
        cols = append(cols,CreateRun(path.Join(dirName,file.Name()),sortkey))
        cols_ndcg = append(cols_ndcg,CreateRun(path.Join(path.Dir(dirName),"results_ndcg",file.Name()),sortkey))
     }
  }

  sort.Sort(BySortkey{cols})
  sort.Sort(BySortkey{cols_ndcg})

  // create and fill the results table
  numCols := len(cols)

/*  for i :=0; i < numCols; i++ {
    if cols[i].sortkey < 500 {
      numCols = i
      break
    }
  }*/

  results := make([][]float64,numCols)

  means_ssa := make([]float64,numCols)
  means_nn := make([]float64,numCols)

  diagonal_ssa := make([]float64,numCols)
  diagonal_nn := make([]float64,numCols)


  for i := range results {
    results[i] = make([]float64,numCols)
  }
  for i := 0; i < numCols; i++ {
    for j := i+1 ; j < numCols; j++ {
      agree, disagree := CalcClassAgreement(cols[i],cols[j])
      results[i][j] = agree
      results[j][i] = disagree
      means_ssa[j] += agree
      means_ssa[i] += agree
      means_nn[j] += disagree
      means_nn[i] += disagree
    }
    diagonal_ssa[i] , diagonal_nn[i] = CalcClassAgreement(cols_ndcg[i],cols[i])
  }

 /* // Just print SSA
  for i := 0; i < numCols; i++ {
    for j := 0 ; j < numCols; j++ {
      if i != j {
        if i < j {
          fmt.Println(cols[i].name,",",cols[j].name,",",results[i][j])
        }
      }
    }
  }
  return*/


  // Print two tables. 
  // SSA first
  // Print the table header
  fmt.Print(collectionName," SSA,");
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name, "(",cols[i].sortkey,"),")
  }
  fmt.Println()

  // Print the table
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name, "(",cols[i].sortkey,"),")
    for j := 0 ; j < numCols; j++ {
      if i != j {
        if i < j {
          fmt.Print(results[i][j],",")
        } else {
          fmt.Print(results[j][i],",")
        }
      } else {
        fmt.Print(",");
      }
    }
    fmt.Println()
  }
  fmt.Print("Mean,")
  for i := 0; i < numCols; i++ {
    fmt.Print(means_ssa[i]/float64(numCols-1),",")
  }
  fmt.Println()
  fmt.Println()

  fmt.Print("map v ndcg,");
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name, "(",cols[i].sortkey,"),")
  }
  fmt.Println()
  fmt.Print("SSA,")
  for i := 0; i < numCols; i++ {
    fmt.Print(diagonal_ssa[i],",")
  }
  fmt.Println()
  fmt.Print("NN,")
  for i := 0; i < numCols; i++ {
    fmt.Print(diagonal_nn[i],",")
  }
  fmt.Println()
  fmt.Println()


  // Now NN
  // Print the table header
  fmt.Print(collectionName," NN,");
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name, "(",cols[i].sortkey,"),")
  }
  fmt.Println()

  // Print the table
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name, "(",cols[i].sortkey,"),")
    for j := 0 ; j < numCols; j++ {
      if i != j {
        if i > j {
          fmt.Print(results[i][j],",")
        } else {
          fmt.Print(results[j][i],",")
        }
      } else {
        fmt.Print(",");
      }
    }
    fmt.Println()
  }
  fmt.Print("Mean,")
  for i := 0; i < numCols; i++ {
    fmt.Print(means_nn[i]/float64(numCols-1),",")
  }
  fmt.Println()
}
