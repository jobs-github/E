#!/bin/bash

branch=$1

function main()
{
    if [ $# -lt 1 ]; then
        echo "no branch"
        return 0
    fi

    git pull --progress -v --no-rebase "origin" ${branch}
}

main $*
