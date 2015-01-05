package main

import (
    "fmt"
    "perl"
    "runtime"
)

func perlHelloWorld(iter int, done_chan chan int) {
    perl.WithInterpreter(func (i *perl.PerlInterpreter) {
        str := fmt.Sprintf("my $iter = %d;\n", iter)
        str += `print "Hello world: $iter\n";`
        i.Eval(str)
    })
    done_chan <- 1
}

func main() {
    runtime.GOMAXPROCS(4);
    var done_chans []chan int
    for i := 0; i < 20 ; i++ {
        done_chan := make(chan int, 1)
        go perlHelloWorld(i, done_chan)
        done_chans = append(done_chans, done_chan)
    }
    for _, ch := range(done_chans) {
       <- ch
    }
    fmt.Println("...ending...")
}
