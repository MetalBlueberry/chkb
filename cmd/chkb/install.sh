#!/bin/bash
export VERSION=$(git describe --dirty)
go install -ldflags="-X 'main.Version=$VERSION'"