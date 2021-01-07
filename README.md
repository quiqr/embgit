# EMBGIT

Embgit is a [PoppyGo](https://poppygo.io) project.

Embgit is an minimal git client made in go. It's main goals are:

- fat binary, dependancy free
- support for all main platforms
- clone, add, commit, push
- use ssh-keys for identification

# Features

- [x] git clone
- [x] git add
- [x] git commit
- [x] git push
- [x] option for alternative ssh-key

## Build

make build

## Cross platform builds

go get github.com/mitchellh/gox

make buildx

