#!/bin/sh

file=$1

if [ -z "$file" ] ; then
    echo "No filename specified"
    exit 1
fi

cflags=`perl -MExtUtils::Embed -e ccopts`
ldflags=`perl -MExtUtils::Embed -e ldopts`

GOPATH=`pwd` CGO_CFLAGS="$cflags" CGO_LDFLAGS="$ldflags" go run -race $file
