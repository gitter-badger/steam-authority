#!/bin/sh

echo $(curl --show-error --silent "http://localhost:8085/admin/rerank-levels")
