#!/bin/bash
cat test.json | json2msgpack | go run crud.go --chains chains-cass.yaml
