# Client

Will be a CLI app (https://dev.to/aurelievache/learning-go-by-examples-part-3-create-a-cli-app-in-go-1h43)[https://dev.to/aurelievache/learning-go-by-examples-part-3-create-a-cli-app-in-go-1h43]
with these functionalities:

- Bulk upload files:

  1. call merkle-tree library to create and locally store the root-hash
  2. call server "bulk-upload" http POST call (POST /files)

- Get a file:
  1. get file and merkle-proofs (http header)
  2. retrieve stored root-hash
  3. verify file with merkle-proofs and root hash (call library merkle-tree)

## Build

```bash
task build
```

## Run

### Upload files (by default is ./testfiles)

```bash
./bin/client-cli upload
```

```bash
./bin/client-cli upload ~/path/to/folder/with/files
```

### Get file and verify merkle tree

```bash
./bin/client-cli get f2
```

## Clean

```bash
task clean
```

### Cobra CLI

Cobra is both a library for creating powerful modern CLI applications and a program for generating applications and batch files.

Download the latest version:

```bash
go get -u github.com/spf13/cobra@latest
```

Then install the cobra CLI:

```bash
go install github.com/spf13/cobra-cli@latest
```

Initialize Cobra application

```bash
cobra-cli init
```

Add a cli command

```bash
cobra-cli add <name-command>
```

### Task

Task is a task runner / build tool that aims to be simpler and easier to use than, for example, GNU Make

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

#### Build

```bash
task build
```

#### Run

```bash
./bin/client-cli get f2
```

#### Clean

```bash
task clean
```

## Create Files

```bash
dd if=/dev/urandom of=./testfiles/f1 bs=1M count=1
dd if=/dev/urandom of=./testfiles/f2 bs=1M count=1
dd if=/dev/urandom of=./testfiles/f3 bs=1M count=1
```

## Upload Files

```bash
go run main.go upload testfiles
```

## Get and verify a File

```bash
go run main.go get f1
```
