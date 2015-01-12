package main

import (
    "fmt"
    "perl"
    "runtime"
)

func testArray(interp *perl.Interpreter) {
    return
//    ref := interp.ObjFromGo("hi there, this is a ref!").Ref()
    res_arr := interp.CallAsArray("do_it", 1, uint(2), float32(1.1),
                                  float64(2.2), "array_mode")
    fmt.Printf("Array result: %v\n", res_arr)
}

func testScalar(interp *perl.Interpreter) {
    return
    res := interp.CallAsScalar("do_it", 1, float32(1.1), float64(2.2),
                               "scalar_mode")
    fmt.Printf("Scalar result: %v\n", res.AsInt())
}

func testConvertHash(interp *perl.Interpreter) {
    m := map[string]*map[string]int{
        "cat": &map[string]int{"val": 1, "wat": 69},
        "cow": &map[string]int{"val2": 2, "wat": 69},
        "dog": &map[string]int{"val3": 3, "wat": 69},
    }
    obj := interp.ObjFromGo(m).Ref()
    res_arr := interp.CallAsArray("do_it2", obj)
    fmt.Printf("Array result: %v\n", res_arr)
    fmt.Printf("cat result: %v\n", res_arr[0].AsHash()["cat"])
}

func testConvertArray(interp *perl.Interpreter) {
    return
    iarr := &[4]int{4, 2, 1, 5}
    obj := interp.ObjFromGo(iarr)
    res_arr := interp.CallAsArray("do_it", obj)
    elem := res_arr[0].AsArray()
    for _, e := range(elem) {
        fmt.Printf("elem %v\n", e)
    }
}

func perlHelloWorld(iter int, done_chan chan int) {
    defer func() { done_chan <- 1 }()
    interp := perl.GetInterpreter()

    str := `sub do_it { use Data::Dumper; print Dumper(@_); @_ }`
    interp.Eval(str)

    str = `sub do_it2 { @_ }`
    interp.Eval(str)

    testScalar(interp)
    testArray(interp)
    testConvertHash(interp)
    testConvertArray(interp)
}

func main() {
    runtime.GOMAXPROCS(4);
    var done_chans []chan int
    for i := 0; i < 10 ; i++ {
        done_chan := make(chan int, 1)
        go perlHelloWorld(i, done_chan)
        done_chans = append(done_chans, done_chan)
    }
    for _, ch := range(done_chans) {
       <- ch
    }
    fmt.Println("...ending...")
}
