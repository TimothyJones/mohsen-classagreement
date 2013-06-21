package main

import (
  "os"
  "io"
  "path"
  "log"
  "fmt"
  "encoding/csv"
  "strconv"
)


// Used to describe an input line from Mohsen's format.
// for example, 
//     name,"yes","no" 
// becomes
//     {a=true, b=false}
type AvB struct {
  a bool
  b bool
}

// This describes a whole Mohsen format file.
//
// name is the basename of the file
// sortkey is the value of the key used to sort at output time
// yesNo is an array of type AvB to represent the data lines
type SubcollectionRun struct {
  name string
  sortkey float64
  yesNo [120]AvB
}

// Used at sort time to hold all of the runs in the input folder
// Implements most of  sort.Interface
type SubcollectionRuns []*SubcollectionRun

//Returns the length of a SubCollectionRuns type.
// Part of sort.Interface. 
func (s SubcollectionRuns) Len() int      { return len(s) }

// Swaps two SubcollectionRun instances in a SubcollectionRuns type.
// Part of sort.Interface
func (s SubcollectionRuns) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Describes a sort by the sosrt key, and completes the implementation of
// sort.Interface
type BySortkey struct{ SubcollectionRuns }

// Part of sort.Interface. Just compares the sort keys of two SubcollectionRun instances.
func (s BySortkey) Less(i, j int) bool { return s.SubcollectionRuns[i].sortkey > s.SubcollectionRuns[j].sortkey }


// Returns 1 if a and b show SSA agreement, 0 otherwise.
//
// See the paper linked in the project description for details
func (a *AvB) ssa(b *AvB) int {
    if (a.a == true && b.a == true) ||
       (a.b == true && b.b == true) {
        return 1
    }
    return 0
}

// Returns 1 if a and b show SSD agreement, 0 otherwise.
//
// See the paper linked in the project description for details
func (a *AvB) ssd(b *AvB) int {
    if (a.a == true && b.b == true) ||
       (a.b == true && b.a == true) {
        return 1
    }
    return 0
}

// Returns 1 if a and b show SN agreement, 0 otherwise.
//
// See the paper linked in the project description for details
func (a *AvB) sn(b *AvB) int {
    if (a.a == true && b.b == false && b.a == false) ||
       (a.b == true && b.b == false && b.a == false) {
        return 1
    }
    return 0
}

// Returns 1 if a and b show NS agreement, 0 otherwise.
//
// See the paper linked in the project description for details
func (a *AvB) ns(b *AvB) int {
    if (b.a == true && a.b == false && a.a == false) ||
       (b.b == true && a.b == false && a.a == false) {
        return 1
    }
    return 0
}

// Returns 1 if a and b show NN agreement, 0 otherwise.
//
// See the paper linked in the project description for details
func (a *AvB) nn(b *AvB) int {
    if a.a == false && b.a == false &&
       a.b == false && b.b == false {
        return 1
    }
    return 0
}


// Calculates the SSA and NN class agreement between two Mohsen files.
// 
// Returns SSA and NN agreement as float64.
func CalcClassAgreement(colA , colB *SubcollectionRun) (agree, disagree float64) {
  var i, ssa, ssd, sn, ns, nn int

  if len(colA.yesNo) != len(colB.yesNo) {
    log.Fatal("Trying to calculate class agreement on non equal length lists");
  }

  for ; i < 120 ; i++ {
    ssa += colA.yesNo[i].ssa(&colB.yesNo[i])
    ssd += colA.yesNo[i].ssd(&colB.yesNo[i])
    sn += colA.yesNo[i].sn(&colB.yesNo[i])
    ns += colA.yesNo[i].ns(&colB.yesNo[i])
    nn += colA.yesNo[i].nn(&colB.yesNo[i])
  }
  if ssa + ssd + sn + ns + nn > len(colA.yesNo) {
    log.Fatal("Too many ssa, ssd, dn, nd, nn - check the tests")
  }
  return float64(2*ssa)/float64(2*ssa+sn+ns),float64(2*nn)/float64(2*nn+sn+ns)
}


// Stores the number of lines in an input file.
// Used to validate the input
var numLines int



// Reads in mohsen format files (see the project description for format details).
//
// Takes the filename to read and an index for the sort key in the header.
func CreateRun(filename string,sortkey int64) *SubcollectionRun {
  file, err := os.Open(filename) // For read access.
  if err != nil {
    log.Fatal(filename," ",  err)
  }

  defer file.Close()
  run := new(SubcollectionRun)
  run.name = path.Base(filename);

  reader := csv.NewReader(file)
  reader.TrailingComma = true;
  // the first record is the sort header
  record , err := reader.Read()
  if err != nil {
    log.Fatal(filename, " ",err)
  }

  run.sortkey, err = strconv.ParseFloat(record[sortkey],64)

  
  fmt.Println(run.name , ": num documents ",record[0], " num evaluated ", record[1], " percent evaluated " , record[2], " num rel " ,record[3], " percent rel ", record[4], " num non-rel " , record[5], " percent non-rel ", record[6]  );

  if err != nil {
    log.Fatal(filename, " ", err)
  }

  i := 0
  var l, r int
  for {
      record, err := reader.Read()
      if err == io.EOF {
          break
      } else if err != nil {
          log.Fatal(filename, " ", err)
      }
      if record[1] == "yes" {
        run.yesNo[i].a = true
        l++
      }
      if record[2] == "yes" {
        run.yesNo[i].b = true
        r++
      }
      i++
  }

  // remember the number of lines if this is our first run
  if numLines == 0 {
    numLines = i
  }

  // did we read the right number of lines?
  if numLines != i || numLines == 0  {
    log.Fatal("Read the incorrect number of lines (",i," vs ",numLines,") when reading '",run.name,"'")
  }

  //fmt.Println(run.name, " read with ",l , " left and ",r," right from ",i, " lines" );

  return run
}
