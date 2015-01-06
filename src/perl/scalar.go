package perl

/*
#include <EXTERN.h>
#include <perl.h>
#include <XSUB.h>

#define _PIARG PerlInterpreter *interp
#define _PSETC PERL_SET_CONTEXT(interp)

static void _SV_incref(_PIARG, SV *sv)
{
    _PSETC;
    SvREFCNT_inc(sv);
}

static void _SV_decref(_PIARG, SV *sv)
{
    _PSETC;
    SvREFCNT_dec(sv);
}

static SV *_newSV_from_int(_PIARG, IV iv)
{
    _PSETC;
    return newSViv(iv);
}

static SV *_newSV_from_uint(_PIARG, UV uv)
{
    _PSETC;
    return newSViv(uv);
}

static SV *_newSV_from_string(_PIARG, char *pv)
{
    _PSETC;
    return newSVpv(pv, 0);
}

static SV *_newSV_from_double(_PIARG, double nv)
{
    _PSETC;
    return newSVnv(nv);
}

static SV *_newSV_from_sv(_PIARG, SV *sv)
{
    _PSETC;
    return newSVsv(sv);
}

*/
import "C"

import (
    "fmt"
    "unsafe"
)

type ScalarType int

const (
    SCALAR_TYPE_INT ScalarType = iota
    SCALAR_TYPE_UINT
    SCALAR_TYPE_STRING
    SCALAR_TYPE_DOUBLE
)

func (stype ScalarType) String() string {
    m := []string{ "int", "uint", "string", "double" }
    return m[stype]
}

type Scalar struct {
    my_perl *C.struct_interpreter
    val interface{}
    stype ScalarType
    sv *C.struct_sv
}

func (scalar *Scalar) String() string {
    return fmt.Sprintf("(%s: %v)", scalar.stype, scalar.val)
}

func (scalar *Scalar) Done() {
    C._SV_decref(scalar.my_perl, scalar.sv)
}

func (interp *Interpreter) NewScalar(arg interface{}) *Scalar {
    my_perl := interp.my_perl
    var sv *C.struct_sv
    var stype ScalarType

    switch val := arg.(type) {
        case *Scalar:
            C._SV_incref(my_perl, val.sv)
            return val
        case int:
            stype = SCALAR_TYPE_INT
            sv = C._newSV_from_int(my_perl, C.IV(val))
        case uint:
            stype = SCALAR_TYPE_UINT
            sv = C._newSV_from_uint(my_perl, C.UV(val))
        case bool:
            stype = SCALAR_TYPE_INT
            if arg.(bool) {
                sv = C._newSV_from_int(my_perl, 1)
            } else {
                sv = C._newSV_from_int(my_perl, 0)
            }
        case string:
            cs := C.CString(arg.(string))
            defer C.free(unsafe.Pointer(cs))
            stype = SCALAR_TYPE_STRING
            sv = C._newSV_from_string(my_perl, cs)
        case float32:
            stype = SCALAR_TYPE_DOUBLE
            sv = C._newSV_from_double(my_perl, C.double(val))
        case float64:
            stype = SCALAR_TYPE_DOUBLE
            sv = C._newSV_from_double(my_perl, C.double(val))
        default:
            panic(fmt.Sprintf("Unsupported type for NewSV: %v\n", val))
    }

    return &Scalar{my_perl: my_perl, val: arg, stype: stype, sv: sv}
}
