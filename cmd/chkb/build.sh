#!/bin/bash
export VERSION=$(git describe --dirty)
go build -ldflags="-X 'main.Version=$VERSION'"