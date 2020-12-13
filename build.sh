#!/bin/bash
export VERSION=$(git describe --dirty)
go build -ldflags="-X 'github.com/MetalBlueberry/chkb/cmd.Version=$VERSION'"