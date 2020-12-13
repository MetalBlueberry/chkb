#!/bin/bash
export VERSION=$(git describe --dirty)
go install -ldflags="-X 'github.com/MetalBlueberry/chkb/cmd.Version=$VERSION'"