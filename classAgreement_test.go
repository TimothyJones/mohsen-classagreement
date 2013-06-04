package main

import (
  "testing"
)

func shouldSSD(a,b *AvB, t *testing.T) {
  if a.ssd(b) != 1 || b.ssd(a) != 1 || a.ssd(a) != 0 || b.ssd(b) != 0{
    t.Error("SSD isn't 1 ",a,b);
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 || a.ssa(a) != 1 || b.ssa(b) != 1{
    t.Error("SSD test case shouldn't be SSA");
  }
  if a.sn(b) == 1 || b.sn(a) == 1 || a.sn(a) == 1 || b.sn(b) == 1{
    t.Error("SSD test case shouldn't be SN");
  }
  if a.ns(b) == 1 || b.ns(a) == 1 || a.ns(a) == 1 || b.ns(b) == 1{
    t.Error("SSD test case shouldn't be NS");
  }
  if a.nn(b) == 1 || b.nn(a) == 1 || a.nn(a) == 1 || b.nn(b) == 1{
    t.Error("SSD test case shouldn't be NN");
  }
}
func shouldNN(a,b *AvB, t *testing.T) {
  if a.nn(b) != 1 || b.nn(a) != 1 || a.nn(a) != 1 || b.nn(b) != 1{
    t.Error("NN isn't 1");
  }
  if a.ssd(b) == 1 || b.ssd(a) == 1 || a.ssd(a) == 1 || b.ssd(b) == 1{
    t.Error("NN test case shouldn't be SSD");
  }
  if a.sn(b) == 1 || b.sn(a) == 1 || a.sn(a) == 1 || b.sn(b) == 1{
    t.Error("NN test case shouldn't be SN");
  }
  if a.ns(b) == 1 || b.ns(a) == 1 || a.ns(a) == 1 || b.ns(b) == 1{
    t.Error("NN test case shouldn't be NS");
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 || a.ssa(a) == 1 || b.ssa(b) == 1{
    t.Error("NN test case shouldn't be SSA");
  }
}

func shouldSSA(a,b *AvB, t *testing.T) {
  if a.ssa(b) != 1 || b.ssa(a) != 1 || a.ssa(a) != 1 || b.ssa(b) != 1{
    t.Error("SSA isn't 1");
  }
  if a.ssd(b) == 1 || b.ssd(a) == 1 || a.ssd(a) == 1 || b.ssd(b) == 1{
    t.Error("SSA test case shouldn't be SSD");
  }
  if a.sn(b) == 1 || b.sn(a) == 1 || a.sn(a) == 1 || b.sn(b) == 1{
    t.Error("SSA test case shouldn't be SN");
  }
  if a.ns(b) == 1 || b.ns(a) == 1 || a.ns(a) == 1 || b.ns(b) == 1{
    t.Error("SSA test case shouldn't be NS");
  }
  if a.nn(b) == 1 || b.nn(a) == 1 || a.nn(a) == 1 || b.nn(b) == 1{
    t.Error("SSA test case shouldn't be NN");
  }
}

func shouldNS(a,b *AvB, t *testing.T) {
  if a.ns(b) != 1 {
    t.Error("NS No ns on b v a ", a ,b);
  }
  if b.sn(a) != 1 {
    t.Error("NS No sn on b v a ", a ,b);
  }
  if a.nn(a) != 1 {
    t.Error("NS No nn on a ", a ,b);
  }
  if b.ssa(b) != 1 {
    t.Error("NS No ssa on b ", a ,b);
  }
  if a.ssd(b) == 1 || b.ssd(a) == 1 || a.ssd(a) == 1 || b.ssd(b) == 1{
    t.Error("NS test case shouldn't be SSD");
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 {
    t.Error("NS test case shouldn't be SSA");
  }
  if a.sn(b) == 1 || a.sn(a) == 1 {
    t.Error("NS test case shouldn't be SN");
  }
  if a.nn(b) == 1 || b.nn(a) == 1 || b.nn(b) == 1{
    t.Error("NS test case shouldn't be NN");
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 {
    t.Error("NS test case shouldn't be SSA");
  }
}


func shouldSN(a,b *AvB, t *testing.T) {
  if a.sn(b) != 1 {
    t.Error("SN No sn on b v a ", a ,b);
  }
  if b.ns(a) != 1 {
    t.Error("SN No ns on b v a ", a ,b);
  }
  if a.ssa(a) != 1 {
    t.Error("SN No ssa on a ", a ,b);
  }
  if b.nn(b) != 1 {
    t.Error("SN No nn on b ", a ,b);
  }
  if a.ssd(b) == 1 || b.ssd(a) == 1 || a.ssd(a) == 1 || b.ssd(b) == 1{
    t.Error("SN test case shouldn't be SSD");
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 {
    t.Error("SN test case shouldn't be SSA");
  }
  if a.ns(b) == 1 || a.ns(a) == 1 {
    t.Error("SN test case shouldn't be NS");
  }
  if a.nn(b) == 1 || b.nn(a) == 1 || a.nn(a) == 1 {
    t.Error("SN test case shouldn't be NN");
  }
  if a.ssa(b) == 1 || b.ssa(a) == 1 {
    t.Error("SN test case shouldn't be SSA");
  }
}





func makeAvB(a,b bool) *AvB {
  var ret AvB
  ret.a = a
  ret.b = b
  return &ret
}

var falseFalse = makeAvB(false,false)
var falseTrue = makeAvB(false,true)
var trueFalse = makeAvB(true,false)

func Test_ssa (t *testing.T) {
  shouldSSA(falseTrue,falseTrue,t)
  shouldSSA(trueFalse,trueFalse,t)
}

func Test_ssd (t *testing.T) {
  shouldSSD(falseTrue,trueFalse,t)
  shouldSSD(trueFalse,falseTrue,t)
}

func Test_sn (t *testing.T) {
  shouldSN(falseTrue,falseFalse,t)
  shouldSN(trueFalse,falseFalse,t)
}

func Test_ns (t *testing.T) {
  shouldNS(falseFalse,falseTrue,t)
  shouldNS(falseFalse,trueFalse,t)
}
func Test_nn (t *testing.T) {
  shouldNN(falseFalse,falseFalse,t)
}


