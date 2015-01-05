package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

extern char **environ;
EXTERN_C void boot_DynaLoader (pTHX_ CV* cv);

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

static void _interp_construct(PerlInterpreter *interp)
{
    char *args[] = { "", "-e", "0" };

    PERL_SET_CONTEXT(interp);
    perl_construct(interp);
    PERL_SET_CONTEXT(interp);
    perl_parse(interp, _xs_init, 3, args, NULL);
    perl_run(interp);
}

static void _interp_destruct(PerlInterpreter *interp)
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

static void _interp_set_context(PerlInterpreter *interp)
{
    PERL_SET_CONTEXT(interp);
}

static void _interp_eval(PerlInterpreter *interp, char *str)
{
    PERL_SET_CONTEXT(interp);
    eval_pv(str, 1);
}

static HV *_hv_new(PerlInterpreter *interp)
{
    PERL_SET_CONTEXT(interp);
    return newHV();
}

static void _hv_decref(PerlInterpreter *interp, HV *hv)
{
    PERL_SET_CONTEXT(interp);
    SvREFCNT_dec(hv);
}

static int _hv_store_str(PerlInterpreter *interp, HV *hv, char *key, char *val)
{
    SV *value;
    SV **sv_ptr;

    PERL_SET_CONTEXT(interp);
    value = newSVpv(val, strlen(val));
    if (value == NULL)
    {
        return -1;
    }

    // steals reference to value unless it returns NULL
    sv_ptr = hv_store(hv, key, strlen(key), value, 0);
    if (sv_ptr == NULL)
    {
        SvREFCNT_dec(value);
        return -1;
    }
    return 0;
}

*/
import "C"

import (
    "os"
    "runtime"
    "strconv"
    "unsafe"
)

type PerlInterpreter struct {
    interp *C.struct_interpreter
}

type PerlHV struct {
    hv *C.struct_hv
}

var doCache bool = true
var interpreters chan *PerlInterpreter

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
    interpreters = make(chan *PerlInterpreter, max_interps)
    C._perl_init()
    if (doCache) {
       doCache = C._perl_can_cache() != 0
    }
}

func doneWithInterpreter(interp *PerlInterpreter) {
    C._interp_destruct(interp.interp)
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
    C.perl_free(interp.interp)
}

func getInterpreter() *PerlInterpreter {
    for {
        select {

        case interp := <- interpreters:
            C._interp_construct(interp.interp)
            return interp

        default:
            return &PerlInterpreter{interp: C._interp_new()}
        }
    }
}

func WithInterpreter(fn func(*PerlInterpreter)) {
    interp := getInterpreter()
    defer doneWithInterpreter(interp)
    fn(interp)
}

func (interp *PerlInterpreter) Eval(s string) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    C._interp_eval(interp.interp, cs)
}

func (interp *PerlInterpreter) doneWithHV(hv *PerlHV) {
    C._hv_decref(interp.interp, hv.hv)
}

func (interp *PerlInterpreter) HVFromStringMap(m map[string]string) *PerlHV {
    chv := C._hv_new(interp.interp);

    for key, val := range(m) {
        ckey, cval := C.CString(key), C.CString(val)
        defer C.free(unsafe.Pointer(ckey))
        defer C.free(unsafe.Pointer(cval))
        if (C._hv_store_str(interp.interp, chv, ckey, cval) < 0) {
            C._hv_decref(interp.interp, chv)
            return nil
        }
    }

    hv := &PerlHV{hv: chv}
    runtime.SetFinalizer(hv, interp.doneWithHV)
    return hv
}
