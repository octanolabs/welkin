#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
wsorgdir="$workspace/src/github.com/octanolabs"
if [ ! -L "$wsorgdir/pristine-go" ]; then
    mkdir -p "$wsorgdir"
    cd "$wsorgdir"
    ln -s ../../../../../. pristine-go
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$wsorgdir/pristine-go"
PWD="$wsorgdir/pristine-go"

# Launch the arguments with the configured environment.
exec "$@"
