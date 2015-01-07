package main

import (
    "fmt"
    "perl"
    "runtime"
)

func perlHelloWorld(iter int, done_chan chan int) {
    perl.WithInterpreter(func (i *perl.Interpreter) {
        str := `use Data::Dumper; sub do_it { print Dumper(@_); @_ }`
        i.Eval(str)
        sv_iter := i.NewScalar(iter)
        defer sv_iter.Done()
        sv_d1 := i.NewScalar(float32(1.1))
        defer sv_d1.Done()
        sv_d2 := i.NewScalar(float64(2.2))
        defer sv_d2.Done()
        sv_str := i.NewScalar("wat")
        s := i.CallAsScalar("kjsfd")
        fmt.Printf("Scalar result: %v\n", s)
        s = i.CallAsScalar("do_it", sv_iter, sv_d1, sv_d2, sv_str, "scalar_mode")
        fmt.Printf("Scalar result: %v\n", s)
        res := i.CallAsArray("do_it", sv_iter, sv_d1, sv_d2, sv_str, "array_mode")
        fmt.Printf("Array result: %v\n", res)
    })
    done_chan <- 1
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
