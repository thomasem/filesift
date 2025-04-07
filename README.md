# FileSift

FileSift is a set of tools for filtering down old files, such as pictures or videos, to clean up old duplicates or similar files. Currently it uses a SHA256 checksum to identify duplicates and store a unique copy of each
file in a specified output directory.

## Build

```
go build -o dedupef ./cmd/dedupef
```

## Basic Usage

```
./dedupef -dir /path/to/source -out /path/to/output
```

## Example with Test Data

```
./dedupef -dir testdata
```

## Options

`-dir string`: source directory to dedupe (default ".")
`-out string`: output directory to store deduplicated copies (default "output")
`-dry-run`: print unique files
