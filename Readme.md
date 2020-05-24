# Shovel Data

A script that checks out all supported [scoop](https://scoop.sh/) buckets and collects the manifests in one searchable json file.

## Motivation

This can be prework for a SPA that allows to search for scoop apps without having to run a server. Also, I wanted to try out Go.

## Run

`go run shovel.go`

## Use

`https://mertd.github.io/shovel-data/manifests.json`

## Licence && Attribution

MIT

The manifests included in this repository were originally created and are maintained by contributors to the following scoop buckets:
* [main](https://github.com/ScoopInstaller/Main)
* [extras](https://github.com/lukesampson/scoop-extras)