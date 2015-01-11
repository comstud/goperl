package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

#define _PIARG PerlInterpreter *interp
#define _PSETC PERL_SET_CONTEXT(interp)

typedef struct _call_response
{
    SV *err_sv;
    int count;
    SV **results;
} _CallResponse;

static void _free_callresponse(_CallResponse *cr)
{
    int i;
    return;

    if (cr == NULL)
    {
        return;
    }

    if (cr->err_sv != NULL)
    {
        SvREFCNT_dec(cr->err_sv);
    }

    for (i = 0;i < cr->count;i++)
    {
        SvREFCNT_dec(cr->results[i]);
    }

    free(cr);
}

static _CallResponse *_interp_call(_PIARG, char *name, SV **svs, int flags)
{
    SV **sv_ptr;
    SV *sv;
    int ret;
    int count;
    _CallResponse *cr;

    _PSETC;
    dSP;
    ENTER;
    SAVETMPS;
    PUSHMARK(SP);

    for (sv_ptr=svs;sv = *sv_ptr;sv_ptr++)
    {
       if (SvTYPE(sv) == SVt_PVHV)
       {
           // expand hashes
           HE *he;
           HV *hv = (HV *)sv;
           SV *he_sv;

           hv_iterinit(hv);
           while ((he = hv_iternext(hv)) != NULL)
           {
               XPUSHs(sv_mortalcopy(HeSVKEY_force(he)));
               XPUSHs(sv_mortalcopy(HeVAL(he)));
           }
       }
       else if (SvTYPE(sv) == SVt_PVAV)
       {
           AV *av = (AV *)sv;
           SSize_t i;
           SSize_t top_index = av_top_index(av);

           for(i=0;i<=top_index;i++)
           {
               XPUSHs(sv_mortalcopy(*av_fetch(av, i, 0)));
           }
       }
       else
       {
           XPUSHs(sv);
       }
    }

    PUTBACK;

    count = perl_call_pv(name, flags | G_EVAL);

    // malloc enough memory to hold the CallResponse struct as well as
    // all of the SV pointers
    cr = malloc(sizeof(*cr) + ((count > 0 ? count + 1 : 1) * sizeof(SV *)));
    if (cr == NULL)
    {
        return NULL;
    }

    cr->err_sv = NULL;
    cr->results = (SV **)(cr + 1); // point to after struct
    cr->count = count;

    SPAGAIN;

    // TODO: handle returning errors and data

    if (SvTRUE(ERRSV))
    {
        cr->err_sv = newSVsv(ERRSV);
    }

    // Results are popped off stack in reverse order
    sv_ptr = cr->results + count;
    // Technically don't need to NULL terminate
    *sv_ptr-- = NULL;

    for(;count > 0;--count,--sv_ptr)
    {
        *sv_ptr = newSVsv(POPs);
    }

    PUTBACK;
    FREETMPS;
    LEAVE;

    return cr;
}

static void _interp_eval(_PIARG, char *str)
{
    _PSETC;
    eval_pv(str, 1);
}

*/
import "C"

import (
    "fmt"
    "unsafe"
)

func (interp *Interpreter) call(name string, mode C.int, args []interface{}) []*Obj {
    my_perl := interp.my_perl

    var sv_arr []*C.SV

    for _, arg := range(args) {
        obj := interp.ObjFromGo(arg)
        sv_arr = append(sv_arr, obj.sv)
    }

    sv_arr = append(sv_arr, nil)

    cs := C.CString(name) 
    defer C.free(unsafe.Pointer(cs))

    res := C._interp_call(my_perl, cs, &sv_arr[0], mode)
    if res == nil {
        panic("perl call failed miserably")
    }

    defer C._free_callresponse(res)

    if res.err_sv != nil {
        obj := interp.ObjFromPerl(res.err_sv)
        fmt.Printf("Got error from perl: %v\n", obj)
        return nil
    }

    count := int(res.count)
    ptrSz := unsafe.Sizeof(*res.results)
    results := uintptr(unsafe.Pointer(res.results))

    var objects []*Obj
    var sv_ptr **C.SV
    
    for i := 0; i < count; i++ {
        sv_ptr = (**C.SV)(unsafe.Pointer(results))
        objects = append(objects, interp.ObjFromPerl(*sv_ptr))
        results += uintptr(ptrSz)
    }

    return objects
}

func (interp *Interpreter) CallAsScalar(name string, args ...interface{}) *Obj {
    results := interp.call(name, C.G_SCALAR, args)
    if len(results) > 0 {
        return results[0]
    }
    return nil
}

func (interp *Interpreter) CallAsArray(name string, args ...interface{}) []*Obj {
    return interp.call(name, C.G_ARRAY, args)
}

func (interp *Interpreter) Eval(s string) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    C._interp_eval(interp.my_perl, cs)
}
