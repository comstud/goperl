#!/bin/sh

file=$1

if [ -z "$file" ] ; then
    echo "No filename specified"
    exit 1
fi

cflags=`perl -MExtUtils::Embed -e ccopts`
ldflags=`perl -MExtUtils::Embed -e ldopts`

gopath=`pwd`
GOPATH=$gopath CGO_CFLAGS="$cflags" CGO_LDFLAGS="$ldflags" go build -race $file
