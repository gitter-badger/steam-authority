#!/usr/bin/env bash

# go get github.com/zackslash/goviz

goviz -i github.com/steam-authority/steam-authority -p | dot -Tpng -o ./dep-graph.png
