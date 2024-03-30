#!/bin/bash

set -e
cd $(dirname "${BASH_SOURCE[0]}")

function usage() {
    echo "usage: add_entity.sh <Name>"
    exit 1
}

name=$1
[[ -z $name ]] && usage && exit 1;

go run -mod=mod entgo.io/ent/cmd/ent new $name

echo "After adding fields to new entity"
echo "Run \"make generate\" in root project folder"