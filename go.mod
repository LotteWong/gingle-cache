module gingle-cache

go 1.13

// Try to use go build instead of go run
require ginglecache v0.0.0

// From go 1.11 version, import packages with relative path
// by using replace `ailas` => `path`
replace ginglecache => ./ginglecache
