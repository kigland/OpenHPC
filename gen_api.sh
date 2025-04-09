#!/bin/bash

rm -fr tmp

mkdir -p tmp

rm -fr coodinator/models/apimod

openapi-generator-cli generate \
    -i OpenHPC.openapi.json \
    -g go-gin-server \
    -o ./tmp \
    --additional-properties=packageName=apimod \
    --global-property models,modelDocs=false \
    --skip-validate-spec

mv tmp/go coodinator/models/apimod