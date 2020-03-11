#!/bin/sh
# installs tools required for makefile to work

set -e

echo "Installing revive"
GO111MODULE=off go get github.com/mgechev/revive

echo "Installing reflex"
GO111MODULE=off go get github.com/cespare/reflex
