#!/bin/bash
export VERSION=$(git describe --tag --dirty)
go build -ldflags="-X 'github.com/MetalBlueberry/chkb/cmd.Version=$VERSION'"