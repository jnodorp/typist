# typist

Emulate typing standard input (or the contents of a file) to standard output. Useful for (recording) demos.

## Usage

Type from file to standard output:

```
typist type README.md
```

Type from standard input to standard output:

```
typist type < README.md
```

Change accuracy:

```
typist --accuracy 0.5 type README.md
# or
typist -a 0.5 type README.md
# or
TYPIST_ACCURACY=0.5 typist type README.md
```

Change words per minute:

```
typist --wpm 100 type README.md
# or
typist -w 100 type README.md
# or
TYPIST_WPM=100 typist type README.md
```

## Installation

```
go install github.com/jnodorp/typist@latest
```
