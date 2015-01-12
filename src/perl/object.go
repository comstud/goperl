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
    return newSVuv(uv);
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

static SV *_newSV_ref(_PIARG, SV *sv)
{
    _PSETC;
    return newRV_inc(sv);
}

static int _SV_type(SV *sv)
{
    return SvTYPE(sv);
}

static int _SV_is_string(_PIARG, SV *sv)
{
    _PSETC;
    return SvPOK(sv);
}

static int _SV_is_int(_PIARG, SV *sv)
{
    _PSETC;
    return SvIOK(sv);
}

static int _SV_is_uint(_PIARG, SV *sv)
{
    _PSETC;
    return SvIOK(sv) && SvIsUV(sv);
}

static int _SV_is_double(_PIARG, SV *sv)
{
    _PSETC;
    return SvNOK(sv);
}

static int _SV_is_ref(_PIARG, SV *sv)
{
    _PSETC;
    return SvROK(sv);
}

static SV *_SV_deref(_PIARG, SV *sv)
{
    _PSETC;
    return SvRV(sv);
}

static char *_SV_as_string(_PIARG, SV *sv, int *len)
{
    STRLEN l;
    char *str;

    _PSETC;
    str = SvPV(sv, l);
    if (len)
        *len = (int)l;
    return str;
}

static int _SV_as_int(_PIARG, SV *sv)
{
    _PSETC;
    return SvIV(sv);
}

static unsigned int _SV_as_uint(_PIARG, SV *sv)
{
    _PSETC;
    return SvUV(sv);
}

static double _SV_as_double(_PIARG, SV *sv)
{
    _PSETC;
    return SvNV(sv);
}

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

static int _hv_store_ent(_PIARG, HV *hv, SV *key, SV *val)
{
    _PSETC;

    // steals references unless it returns NULL
    if (hv_store_ent(hv, key, val, 0) != NULL)
    {
        SvREFCNT_inc(key);
        SvREFCNT_inc(val);
        return 0;
    }
    return -1;
}

static AV *_av_new(_PIARG)
{
    _PSETC;
    return newAV();
}

static void _av_incref(_PIARG, AV *av)
{
    _PSETC;
    SvREFCNT_inc(av);
}

static void _av_decref(_PIARG, AV *av)
{
    _PSETC;
    SvREFCNT_dec(av);
}

static int _av_store(_PIARG, AV *av, SSize_t index, SV *val)
{
    _PSETC;
    if (av_store(av, index, val) != NULL)
    {
        SvREFCNT_inc(val);
        return 0;
    }
    return -1;
}

static SSize_t _av_top_index(_PIARG, AV *av)
{
    _PSETC;
    // These are the same. Apparently because av_len is confusing in
    // that it returns the top index (and not the actual length), they
    // added av_top_index. Prefer that.
#ifdef av_top_index
    return av_top_index(av);
#else
    return av_len(av);
#endif
}

static SV *_av_fetch(_PIARG, AV *av, SSize_t i)
{
    _PSETC;
    return *av_fetch(av, i, 0);
}

static void _hv_iterinit(_PIARG, HV *hv)
{
    _PSETC;
    hv_iterinit(hv);
}

static HE *_hv_iternext(_PIARG, HV *hv)
{
    _PSETC;
    return hv_iternext(hv);
}

static void _he_get_keyval(_PIARG, HE *he, SV **key_ret, SV **val_ret)
{
    _PSETC;
    *key_ret = newSVsv(HeSVKEY_force(he));
    *val_ret = newSVsv(HeVAL(he));
}

*/
import "C"

import (
    "fmt"
    "reflect"
    "runtime"
    "unsafe"
)

type Obj struct {
    interp *Interpreter
    sv     *C.struct_sv
}

func (obj *Obj) String() string {
    return fmt.Sprintf("%v", obj.Val())
}

func doneWithObj(obj *Obj) {
    C._SV_decref(obj.interp.my_perl, obj.sv)
}

func newObj(interp *Interpreter, sv *C.SV) *Obj {
    obj := &Obj{interp: interp, sv: sv}
    runtime.SetFinalizer(obj, doneWithObj)
    C._SV_incref(interp.my_perl, sv)
    return obj
}

func avToGo(my_perl *C.struct_interpreter, av *C.AV) []interface{} {
    iarr := []interface{}{}
    top_index := C._av_top_index(my_perl, av)
    for i := C.SSize_t(0); i <= top_index; i++ {
        elem := svToGo(my_perl, C._av_fetch(my_perl, av, i))
        iarr = append(iarr, elem)
    }
    return iarr
}

func hvToGo(my_perl *C.struct_interpreter, hv *C.HV) map[interface{}]interface{} {
    imap := make(map[interface{}]interface{})
    C._hv_iterinit(my_perl, hv)
    for {
        he := C._hv_iternext(my_perl, hv)
        if he == nil {
            break
        }

        var key, val *C.SV

        C._he_get_keyval(my_perl, he, &key, &val)
        defer C._SV_decref(my_perl, key)
        defer C._SV_decref(my_perl, val)

        keyobj := svToGo(my_perl, key)
        valobj := svToGo(my_perl, val)

        imap[keyobj] = valobj
    }

    return imap
}

func svToGo(my_perl *C.struct_interpreter, sv *C.SV) interface{} {
    var i interface{}
    var is_ref bool

    if C._SV_is_ref(my_perl, sv) != 0 {
        sv = C._SV_deref(my_perl, sv)
        is_ref = true
    } else {
        is_ref = false
    }
    switch C._SV_type(sv) {
        case C.SVt_NULL:
            i = nil
        case C.SVt_IV, C.SVt_PVIV:
            if C._SV_is_uint(my_perl, sv) != 0 {
                i = uint(C._SV_as_int(my_perl, sv))
            } else {
                i = int(C._SV_as_int(my_perl, sv))
            }
        case C.SVt_NV, C.SVt_PVNV:
            i = float64(C._SV_as_double(my_perl, sv))
        case C.SVt_PV:
            var len C.int
            tmp := C._SV_as_string(my_perl, sv, &len)
            i = C.GoStringN(tmp, len)
        case C.SVt_PVHV:
            hv := (*C.HV)(unsafe.Pointer(sv))
            return hvToGo(my_perl, hv)
        case C.SVt_PVAV:
            av := (*C.AV)(unsafe.Pointer(sv))
            return avToGo(my_perl, av)
        default:
            fmt.Printf("Unsupported SV type: %v\n", sv)
            return nil
    }

    if is_ref {
        i = &i
    }

    return i
}

func (obj *Obj) Val() interface{} {
    my_perl := obj.interp.my_perl
    return svToGo(my_perl, obj.sv)
}

func (obj *Obj) AsInt() int {
    my_perl := obj.interp.my_perl
    if C._SV_is_int(my_perl, obj.sv) != 0 {
        return int(C._SV_as_int(my_perl, obj.sv))
    }
    panic("perl.Obj type not compatible with int")
}

func (obj *Obj) AsUInt() uint {
    my_perl := obj.interp.my_perl
    if C._SV_is_uint(my_perl, obj.sv) != 0 {
        return uint(C._SV_as_uint(my_perl, obj.sv))
    }
    panic("perl.Obj type not compatible with uint")
}

func (obj *Obj) AsFloat() float64 {
    my_perl := obj.interp.my_perl
    if C._SV_is_double(my_perl, obj.sv) != 0 {
        return float64(C._SV_as_double(my_perl, obj.sv))
    }
    panic("perl.Obj type not compatible with float")
}

func (obj *Obj) AsString() string {
    my_perl := obj.interp.my_perl
    if C._SV_is_string(my_perl, obj.sv) != 0 {
        var len C.int
        tmp := C._SV_as_string(my_perl, obj.sv, &len)
        return C.GoStringN(tmp, len)
    }
    panic("perl.Obj type not compatible with string")
}

func (obj *Obj) AsArray() []interface{} {
    my_perl := obj.interp.my_perl
    var sv *C.SV

    if C._SV_is_ref(my_perl, obj.sv) != 0 {
        sv = C._SV_deref(my_perl, obj.sv)
    } else {
        sv = obj.sv
    }

    if C._SV_type(sv) != C.SVt_PVAV {
        panic("perl.Obj type not compatible with array")
    }

    av := (*C.AV)(unsafe.Pointer(sv))
    return avToGo(my_perl, av)
}

func (obj *Obj) AsHash() map[interface{}]interface{} {
    my_perl := obj.interp.my_perl
    var sv *C.SV

    if C._SV_is_ref(my_perl, obj.sv) != 0 {
        sv = C._SV_deref(my_perl, obj.sv)
    } else {
        sv = obj.sv
    }

    if C._SV_type(sv) != C.SVt_PVHV {
        panic("perl.Obj type not compatible with hash")
    }

    hv := (*C.HV)(unsafe.Pointer(sv))
    return hvToGo(my_perl, hv)
}

func (obj *Obj) Ref() *Obj {
    my_perl := obj.interp.my_perl
    sv := C._newSV_ref(my_perl, obj.sv)
    defer C._SV_decref(my_perl, sv)
    return newObj(obj.interp, sv)
}

func (interp *Interpreter) objFromReflectValue(val reflect.Value) *Obj {
    my_perl := interp.my_perl
    var sv *C.SV

    switch val.Kind() {
        case reflect.Ptr:
            i := val.Interface()
            switch ival := i.(type) {
                case *Obj:
                    return ival
            }
            return interp.objFromReflectValue(val.Elem()).Ref()
        case reflect.String:
            cs := C.CString(val.String())
            defer C.free(unsafe.Pointer(cs))
            sv = C._newSV_from_string(my_perl, cs)
        case reflect.Int, reflect.Int8, reflect.Int16,
             reflect.Int32, reflect.Int64:
            sv = C._newSV_from_int(my_perl, C.IV(val.Int()))
        case reflect.Uint, reflect.Uint8, reflect.Uint16,
             reflect.Uint32, reflect.Uint64:
            sv = C._newSV_from_uint(my_perl, C.UV(val.Uint()))
        case reflect.Bool:
            if val.Bool() {
                sv = C._newSV_from_int(my_perl, 1)
            } else {
                sv = C._newSV_from_int(my_perl, 0)
            }
        case reflect.Float32, reflect.Float64:
            sv = C._newSV_from_double(my_perl, C.double(val.Float()))
        case reflect.Map:
            hv := C._hv_new(my_perl)
            for _, key := range val.MapKeys() {
               val := val.MapIndex(key)
               keyobj := interp.objFromReflectValue(key)
               valobj := interp.objFromReflectValue(val)
               if C._hv_store_ent(my_perl, hv, keyobj.sv, valobj.sv) < 0 {
                   C._hv_decref(my_perl, hv)
                   panic("failed to add entry to perl hash")
               }
               C._SV_incref(my_perl, keyobj.sv)
               C._SV_incref(my_perl, valobj.sv)
            }
            sv = (*C.SV)(unsafe.Pointer(hv))
        case reflect.Array:
            av := C._av_new(my_perl)
            for i := 0 ; i < val.Len(); i++ {
                valobj := interp.objFromReflectValue(val.Index(i))
                if C._av_store(my_perl, av, C.SSize_t(i), valobj.sv) < 0 {
                    C._av_decref(my_perl, av)
                    panic("failed to add entry to perl array")
                }
                C._SV_incref(my_perl, valobj.sv)
            }
            sv = (*C.SV)(unsafe.Pointer(av))
        default:
            panic(fmt.Sprintf("Unsupported type for NewSV: %v\n", val))
    }
    defer C._SV_decref(my_perl, sv)
    return newObj(interp, sv)
}

func (interp *Interpreter) ObjFromGo(arg interface{}) *Obj {
    return interp.objFromReflectValue(reflect.ValueOf(arg))
}

func (interp *Interpreter) ObjFromPerl(sv *C.struct_sv) *Obj {
    return newObj(interp, sv)
}
