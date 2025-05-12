# ECWID interview [test task](https://github.com/Ecwid/new-job/blob/master/IP-Addr-Counter-GO.md)

## Content

- [Build](#build)
- [Run](#run)
- [Test](#test)
- [Results](#results)

## Build

To build the executable file just execute make command

```bash
    make build
```

## Run

To file you need initialize environment variables.

```bash
    export FILE_PATH="path/to/your/ip_addresses"
    export STRATEGY="linear|concurrent"
```

after environment initialization you can execute command with make.

```bash
    make run
```

You can execute command in the bin/ directory after [build](#build) stage to get more information.

```bash
    ./ip_extractor --help
```

## Test

To run test execute the command.

```bash
    make test
```

## Results

Started on Machine with 20 CPUs 16gb RAM on arch amd64. With 10mb chunk size.
File size 107gb downloaded from [here](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip).

| Strategy  | Time elapsed | Total allocated mb |
|-----------|--------------|--------------------|
| concurrent| 51s          | 692mb              |
| linear    | 4m39s        | 522mb              |
