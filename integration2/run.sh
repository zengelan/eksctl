#!/bin/sh
exec $GOPATH/bin/ginkgo -tags integration -v -p  ./integration2
