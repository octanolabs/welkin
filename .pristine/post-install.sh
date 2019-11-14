#!/bin/bash

BLUE="\033[0;34m"
RED="\033[0;31m"
NC="\033[0m" # No Color

defaultBinary="pristine"
defaultOrg="octanolabs"
defaultRepo="pristine-go"

echo ""
echo "💎  Welcome Pristine Go Post-Install setup! 💎"
echo ""
echo ""

echo -e "${BLUE}Enter the github org/user that owns this projects repository (e.g ${defaultOrg}):${NC}"

read orginization

echo ""

echo -e "${BLUE}Enter the github repository name for this project (e.g ${defaultRepo}):${NC}"

read repository

echo ""

echo -e "${BLUE}Enter the executable/binary name (e.g ${defaultBinary}):${NC}"

read binary

echo ""

function replaceTextInFile() {
  # using ~ in place of / to avoid slashes in package names conflicting with sed
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i  "" -e "s~$1~$2~g" $3
  else
    sed -i  -e "s~$1~$2~g" $3
  fi
}

defaultGithub="${defaultOrg}/${defaultRepo}"
github="${orginization}/${repository}"

replaceTextInFile $defaultGithub $github  .circleci/config.yml
replaceTextInFile $defaultBinary $binary  .circleci/config.yml

replaceTextInFile $defaultGithub $github  internal/build/env.go

replaceTextInFile $defaultGithub $github  build/ci.go
replaceTextInFile $defaultBinary $binary  build/ci.go
# using ~ in place of / to avoid slashes in package names conflicting with sed

replaceTextInFile $defaultRepo $repository  build/env.sh
replaceTextInFile $defaultOrg $orginization  build/env.sh

replaceTextInFile $defaultBinary $binary  Makefile

replaceTextInFile $defaultGithub $github  README.md
replaceTextInFile $defaultRepo $repository  README.md

git mv "cmd/${defaultBinary}" "cmd/${binary}"
# build/env.sh

echo -e "${BLUE} 🚀  Project Setup Completed. 🚀"

echo ""
