# SCONFIG - Simple configuration file management for the Linux Shell

## What is it
SCONFIG is a simple configuration file management tool for the Linux Shell (bash, ksh, csh, etc).
The purpose is to easily make changes in configuration files (TOML, JSON, YAML java.properties files).

Example:
```bash
sconfig config.toml fastsync.version=v1
sconfig genesis.json consensus_params.block.max_bytes=22020095 --type string
```
The first example will override (or add) the entry `version="v1"` into the `[fastsync]` section of `config.toml`.

The second example will override (or add) the entry `max_bytes="22020095"` in a json map. If the entry was an integer,
it will be converted into a string.

## How to get it
Check the [releases](https://github.com/freshautomations/sconfig/releases) page.

## How to build it from source
Install [Golang](https://golang.org/doc/install), then run:
```bash
go install github.com/freshautomations/sconfig@latest
```

## How to use
```
Usage:
  sconfig <filename> <key=value> [<key=value>] ... [flags]

Flags:
  -h, --help          help for sconfig
  -s, --strict        Only allow changes but not new entries.
  -t, --type string   Override value(s) type.
  -v, --version       version for sconfig
```

By default, sconfig will add a new entry to the config file, if it did not exist before. This can be limited using `-s`
in which case only existing entries can be updated.
For existing entries, sconfig will convert the input string into the type that of the configuration item. If it cannot
convert it, it will return with an error. This can be overridden by defining what type that new config parameter should be
with `-t`.
New entries are `string` type, unless overridden by `-t`.
If `-t` is present, all key=value pairs are going to be forced into that type.

### Accepted types
* bool, boolean
* float, float32, float64
* int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr
* intslice
* stringslice
* time
* duration
* string, str

Map types are not supported but you can use the dotted notation (`section.key=value`) to build maps one entry at a time.

## Caveats
File extensions are important.

The Golang Viper library that is used for sconfig will load and save a file based on its file extension.

## Examples:
Commands ran:
```bash
sconfig <file> myint=4 myfakeint=eleven
sconfig <file> myrealfloat=3.2 -t float
sconfig <file> myintlist=[3,4,5] -t intslice
```

File: example.toml
```toml
myint=3
myfakeint="1"
mystring="bye"
myintlist=[1,2,3]
mystringlist=["hello","goodbye"]
```

Result of commands:
```toml
myfakeint = "eleven"
myint = 4
myintlist = [3,4,5]
myrealfloat = 3.2
mystring = "bye"
mystringlist = ["hello","goodbye"]
```

File: example.json
```json
{
  "myint": 3,
  "myfakeint": "1",
  "mystring": "bye",
  "myintlist": [
    1,
    2,
    3
  ],
  "mystringlist": [
    "hello",
    "goodbye"
  ]
}
```

Result of commands:
```json
{
  "myfakeint": "eleven",
  "myint": 4,
  "myintlist": [
    3,
    4,
    5
  ],
  "myrealfloat": 3.2,
  "mystring": "bye",
  "mystringlist": [
    "hello",
    "goodbye"
  ]
}
```
