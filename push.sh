#!/bin/bash

go mod tidy
git add .
git commit -m "$1"
git push
