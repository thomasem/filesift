# FileSift

FileSift is a simple tool for filtering down files to clean up duplicates or find similarities. One primary use-case
is cleaning up old pictures and videos. Currently it uses a SHA256 checksum to identify duplicates and store a unique
copy of each file in a specified output directory.

## Build (without Installing)

```
go build -o filesift ./cmd/filesift
./filesift --help
```

## Install

```
go install ./cmd/filesift
filesift --help
```

## Basic Usage

```
filesift <command> [options]
```

### Dedupe

Deduplicates files and stores unique copies in a separate directory to preserve the source files.

```
filesift dedupe [options]
```

#### Options

`-src`: Directory where source files exist (default: ".")

`-out`: Directory where deduplicated files will be stored (default: './output')

`-dry-run`: Print unique files and exit

#### Example

```
$ filesift dedupe -src testdata -dry-run
testdata/text_a
testdata/text_b
testdata/bumble.png
testdata/elsa_sploot.jpg
testdata/elsa_walk.jpg
testdata/flower.jpg
```
