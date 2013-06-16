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

type AvB struct {
  a bool
  b bool
}

type SubcollectionRun struct {
  name string
  sortkey float64
  yesNo [120]AvB
}

type SubcollectionRuns []*SubcollectionRun

func (s SubcollectionRuns) Len() int      { return len(s) }
func (s SubcollectionRuns) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type BySortkey struct{ SubcollectionRuns }

func (s BySortkey) Less(i, j int) bool { return s.SubcollectionRuns[i].sortkey > s.SubcollectionRuns[j].sortkey }


func (a *AvB) ssa(b *AvB) int {
    if (a.a == true && b.a == true) ||
       (a.b == true && b.b == true) {
        return 1
    }
    return 0
}

func (a *AvB) ssd(b *AvB) int {
    if (a.a == true && b.b == true) ||
       (a.b == true && b.a == true) {
        return 1
    }
    return 0
}

func (a *AvB) sn(b *AvB) int {
    if (a.a == true && b.b == false && b.a == false) ||
       (a.b == true && b.b == false && b.a == false) {
        return 1
    }
    return 0
}

func (a *AvB) ns(b *AvB) int {
    if (b.a == true && a.b == false && a.a == false) ||
       (b.b == true && a.b == false && a.a == false) {
        return 1
    }
    return 0
}

func (a *AvB) nn(b *AvB) int {
    if a.a == false && b.a == false &&
       a.b == false && b.b == false {
        return 1
    }
    return 0
}

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


var numLines int


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
