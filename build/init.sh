#!/bin/sh

set -e

# check if go-version is less than 1.13
echo "Checking Go version"
GO_VERSION=$(go version | cut -d ' ' -f 3 | sed -n 's/go//p' | cut -d '.' -f 2)
if [ "$GO_VERSION" -lt 13 ]; then
  echo "Go version is less than 1.13. Please upgrade"
  exit 1
fi

export BASE=${PWD##*/}

replace() {
    echo "Replacing go-template with $BASE in $1"
    sed -i "s/go\-template/$BASE/g" "$1"
}

echo "Updating project BIN and PKG with '$BASE'"
replace envfile

echo "Replacing paths in source-files with '$BASE'"
replace cmd/go-template/main.go
replace cmd/go-template/main_test.go
replace internal/ping/v1/handler.go
replace go.mod

echo "Moving go-files in cmd/go-template to cmd/$BASE"
mv cmd/go-template "cmd/$BASE"

echo "installing tools and updating lint-rules"
make install-tools lint-update

echo "Building README"
printf "# %s\n\n" "$BASE" > README.md
cat README.in.md >> README.md
rm README.in.md

echo "Replacing doc title"
echo "$BASE" > md/_title

echo "Removing .git"
rm -rf .git

echo "Initialising new git repo"
git init

echo "Setting git remote to git@bitbucket.org:tktip/$BASE.git"
git remote add origin "git@bitbucket.org:tktip/$BASE.git"

echo "removing init-functionality"
sed -i "/#__start__/,/#__end__/c\\" Makefile
rm build/init.sh

echo "Creating initial commit"
git add .
git commit -m "init commit"

if ! git push --set-upstream origin master; then
    echo "Initialisation partially completed: Was unable to push and set remote repository"
    echo "Create the remote repository and ensure the git-command has the correct access rights, and run \"git push --set-upstream origin master\""
    echo "Or open your git-client of choice and complete the setup from there"
else
    echo "Shit's good yo. Happy hacking!"
fi
