package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

#define _PIARG PerlInterpreter *interp
#define _PSETC PERL_SET_CONTEXT(interp)

static SV **_interp_call(_PIARG, char *name, SV **svs, int flags)
{
    SV **sv_ptr;
    SV *sv;
    int count;

    _PSETC;
    dSP;
    ENTER;
    SAVETMPS;
    PUSHMARK(SP);

    for (sv_ptr=svs;sv = *sv_ptr;sv_ptr++)
    {
       XPUSHs(sv);
    }

    PUTBACK;

    count = perl_call_pv(name, flags | G_EVAL);

    SPAGAIN;

    // TODO: handle returning errors and data

    if (SvTRUE(ERRSV))
    {
        printf("Got error: %s\n", SvPV_nolen(ERRSV));
    }

    for(;count > 0;--count)
    {
        // NOTE: Must INCREF sv if we keep it
        sv = POPs;
    }

    PUTBACK;
    FREETMPS;
    LEAVE;

    return NULL;
}

static void _interp_eval(_PIARG, char *str)
{
    _PSETC;
    eval_pv(str, 1);
}

*/
import "C"

import (
    "unsafe"
)

func (interp *Interpreter) call(name string, mode C.int, args []interface{}) []*Scalar {
    my_perl := interp.my_perl

    var sv_arr []*C.SV

    for _, arg := range(args) {
        scalar := interp.NewScalar(arg)
        defer scalar.Done()
        sv_arr = append(sv_arr, scalar.sv)
    }

    sv_arr = append(sv_arr, nil)

    cs := C.CString(name) 
    defer C.free(unsafe.Pointer(cs))

    C._interp_call(my_perl, cs, &sv_arr[0], mode)
    return []*Scalar{}
}

func (interp *Interpreter) CallAsScalar(name string, args ...interface{}) *Scalar {
    results := interp.call(name, C.G_SCALAR, args)
    if len(results) > 0 {
        return results[0]
    }
    return nil
}

func (interp *Interpreter) CallAsArray(name string, args ...interface{}) []*Scalar {
    return interp.call(name, C.G_ARRAY, args)
}

func (interp *Interpreter) Eval(s string) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    C._interp_eval(interp.my_perl, cs)
}
