package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "path"
  "log"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println(os.Args[0], ": calculate class agreement between multiple subcollections")
    fmt.Println("Usage:")
    fmt.Println("   ",os.Args[0], " <path to directory containing subcollection significance csv files> ")
    return
  }

  // Read the collection CSVs
  dirName :=  os.Args[1]
  collectionName := path.Base(dirName)
  fileListing, err:= ioutil.ReadDir(dirName);
  if err != nil {
    log.Fatal(err)
  }
  cols := make([]*SubcollectionRun,0,len(fileListing))
  for _,file := range fileListing {
     if path.Ext(file.Name()) == ".csv"{
        cols = append(cols,CreateRun(path.Join(dirName,file.Name())))
     }
  }

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
    fmt.Print(cols[i].name,",")
  }
  fmt.Println()

  // Print the table
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name,",")
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
    fmt.Print(cols[i].name,",")
  }
  fmt.Println()

  // Print the table
  for i := 0; i < numCols; i++ {
    fmt.Print(cols[i].name,",")
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
