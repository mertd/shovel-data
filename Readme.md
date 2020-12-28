# Shovel Data

A script that checks out all supported [scoop](https://scoop.sh/) buckets and collects the manifests in one searchable json file.

Use [Shovel](https://shovel.sh) to access the search.

## Motivation

This is the prework for the [Shovel](https://shovel.sh) SPA that allows to search for scoop apps without having to run a server. Also, I wanted to try out Go.

## Run

`go get && go run shovel.go`

## Use

`https://mertd.github.io/shovel-data/manifests.json`

# Licence

MIT

# Attribution

Refer to `Attribution.md`.

If a bucket you are looking for is not supported, it either is not a [known bucket](https://github.com/lukesampson/scoop/blob/master/buckets.json) or its licence is, in my amateur knowledge of the law (I am not a lawyer), incompatible with the MIT licence.