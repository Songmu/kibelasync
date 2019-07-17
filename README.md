kibelasync
=======

[![Build Status](https://travis-ci.org/Songmu/kibelasync.svg?branch=master)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/kibelasync/badge.svg?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/kibelasync?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/kibelasync
[coveralls]: https://coveralls.io/r/Songmu/kibelasync?branch=master
[license]: https://github.com/Songmu/kibelasync/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/kibelasync

kibelasync is a CLI for https://kibe.la

## Synopsis

```console
# set $KIBELA_TOKEN and $KIBELA_TEAM environment variable before using

% kibelasync pull
[kibelasync] saved to "notes/370.md"
[kibelasync] saved to "notes/381.md"
[kibelasync] saved to "notes/380.md"
...

% kibelasync push notes/370.md
[kibelasync] updated https://example.kibe.la/notes/370

% kibelasync publish < sample.md
[kibelasync] published https://songmu.kibe.la/@Songmu/382
```

## Description

kibela client to edit markdowns locally. It download markdowns with frontmatter.

## Installation

### Homebrew

```console
% brew install Songmu/tap/kibelasync
```

### go get

```console
% go get github.com/Songmu/kibelasync/cmd/kibelasync
```

## See Also

- [blogsync](https://github.com/motemen/blogsync)

## Author

[Songmu](https://github.com/Songmu)
