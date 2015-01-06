package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

extern char **environ;
EXTERN_C void boot_DynaLoader (pTHX_ CV* cv);

#define _PIARG PerlInterpreter *interp
#define _PSETC PERL_SET_CONTEXT(interp)

static HV *_hv_new(_PIARG)
{
    _PSETC;
    return newHV();
}

static void _hv_incref(_PIARG, HV *hv)
{
    _PSETC;
    SvREFCNT_inc(hv);
}

static void _hv_decref(_PIARG, HV *hv)
{
    _PSETC;
    SvREFCNT_dec(hv);
}

static int _hv_store_str(_PIARG, HV *hv, char *key, char *val)
{
    SV *value;
    SV **sv_ptr;

    _PSETC;
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
    "unsafe"
)

type Hash struct {
    my_perl *C.struct_interpreter
    hv *C.struct_hv
}

func (hash *Hash) Done() {
    C._hv_decref(hash.my_perl, hash.hv)
}

func (interp *Interpreter) HashFromStringMap(m map[string]string) *Hash {
    my_perl := interp.my_perl
    hv := C._hv_new(my_perl)

    for key, val := range(m) {
        ckey, cval := C.CString(key), C.CString(val)
        defer C.free(unsafe.Pointer(ckey))
        defer C.free(unsafe.Pointer(cval))
        if (C._hv_store_str(interp.my_perl, hv, ckey, cval) < 0) {
            C._hv_decref(interp.my_perl, hv)
            return nil
        }
    }

    return &Hash{my_perl: my_perl, hv: hv}
}
