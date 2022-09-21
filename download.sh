#!/bin/bash

curl -X POST http://127.0.0.1:7882/download?hash=$1
echo ""
