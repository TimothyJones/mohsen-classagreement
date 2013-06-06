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
  for _,file := range fileListing {
     if path.Ext(file.Name()) == ".csv"{
        cols = append(cols,CreateRun(path.Join(dirName,file.Name()),sortkey))
     }
  }

  sort.Sort(BySortkey{cols})

  // create and fill the results table
  numCols := len(cols)
  results := make([][]float64,numCols)
  for i := range results {
    results[i] = make([]float64,numCols)
  }
  for i := 0; i < numCols; i++ {
    for j := i+1 ; j < numCols; j++ {
      agree, disagree := CalcClassAgreement(cols[i],cols[j])
      results[i][j] = agree
      results[j][i] = disagree
    }
  }

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
}
