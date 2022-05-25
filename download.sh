#!/bin/bash

curl -X POST http://127.0.0.1:8080/download?hash=$1
echo ""
