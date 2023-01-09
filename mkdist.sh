#!/bin/bash

rm dist.zip
zip -r dist.zip ./ -x .git/\* solution/\* mkdist.sh .devcontainer/\*