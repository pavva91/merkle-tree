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

## Cobra CLI

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

## Create Files

```bash
dd if=/dev/urandom of=./testfiles/f1 bs=1M count=1
dd if=/dev/urandom of=./testfiles/f2 bs=1M count=1
dd if=/dev/urandom of=./testfiles/f3 bs=1M count=1
```

## Upload Files
