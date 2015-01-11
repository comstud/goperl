package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

extern char **environ;
EXTERN_C void boot_DynaLoader (pTHX_ CV* cv);

#define _PIARG PerlInterpreter *interp
#define _PSETC PERL_SET_CONTEXT(interp)

static void _xs_init(PTHX)
{
    char *file = __FILE__;

    newXS("DynaLoader::boot_DynaLoader", boot_DynaLoader, file);
}

static void _perl_init(void)
{
    int num_args = 1;
    char *args[] = { "perl", NULL };

    PERL_SYS_INIT3(&num_args, (char ***)&args, &environ);
}

static int _perl_can_cache(void)
{
    //
    // I haven't tested 5.16, but 5.14 is certainly broken. When
    // trying to call perl_parse() a 2nd time on the same interpreter,
    // you get a warning: "Attempt to free unreferenced scalar"
    // and then a later crash.
    //
    return PERL_VERSION > 14;
}

static void _interp_construct(_PIARG)
{
    char *args[] = { "", "-e", "0" };

    _PSETC;
    perl_construct(interp);
    perl_parse(interp, _xs_init, 3, args, NULL);
    perl_run(interp);
}

static void _interp_destruct(_PIARG)
{
    PERL_SET_CONTEXT(interp);
    perl_destruct(interp);
}

static PerlInterpreter *_interp_new(void)
{
    PerlInterpreter *interp = perl_alloc();

    if (interp != NULL)
    {
        _interp_construct(interp);
    }
    return interp;
}

static void _interp_set_context(_PIARG)
{
    _PSETC;
}

*/
import "C"

import (
    "os"
    "runtime"
    "strconv"
)

type Interpreter struct {
    my_perl *C.struct_interpreter
}

var doCache bool = true
var interpreters chan *Interpreter

func init() {
    var max_interps int

    max_interps_env := os.Getenv("GO_MAX_PERL_INTERPRETERS")
    if max_interps_env != "" {
       var err error
       max_interps, err = strconv.Atoi(max_interps_env)
       if err != nil || max_interps <= 0 {
           doCache = false
       }
    }
    if max_interps <= 0 {
       max_interps = runtime.NumCPU()
    }
    interpreters = make(chan *Interpreter, max_interps)
    C._perl_init()
    if (doCache) {
       doCache = C._perl_can_cache() != 0
    }
}

func doneWithInterpreter(interp *Interpreter) {
    C._interp_destruct(interp.my_perl)
    if (doCache) {
        for {
            select {
                case interpreters <- interp:
                    return
                default:
            }
            break
        }
    }
    C.perl_free(interp.my_perl)
}

func GetInterpreter() *Interpreter {
    var interp *Interpreter
    select {
        case interp = <- interpreters:
            C._interp_construct(interp.my_perl)
        default:
            interp = &Interpreter{my_perl: C._interp_new()}
    }
    runtime.SetFinalizer(interp, doneWithInterpreter)
    return interp
}
