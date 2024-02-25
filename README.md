# Coding Challenge

[![Go](https://github.com/pavva91/merkle-tree/actions/workflows/go.yml/badge.svg)](https://github.com/pavva91/merkle-tree/actions/workflows/go.yml)

[![Go Report Card Merkle Tree Library](https://goreportcard.com/badge/github.com/pavva91/merkle-tree/libs/merkletree)](https://goreportcard.com/report/github.com/pavva91/merkle-tree/libs/merkletree)

[![Go Report Card Server](https://goreportcard.com/badge/github.com/pavva91/merkle-tree/server)](https://goreportcard.com/report/github.com/pavva91/merkle-tree/server)

[![Go Report Card Client](https://goreportcard.com/badge/github.com/pavva91/merkle-tree/client)](https://goreportcard.com/report/github.com/pavva91/merkle-tree/client)

## Description

Imagine a client has a large set of potentially small files {F0, F1, …, Fn} and wants to upload them to a server and then delete its local copies.
The client wants, however, to later download an arbitrary file from the server and be convinced that the file is correct and is not corrupted in any way (in transport, tampered with by the server, etc.).

You should implement the client, the server and a Merkle tree to support the above (we expect you to implement the Merkle tree rather than use a library, but you are free to use a library for the underlying hash functions).

The client must compute a single Merkle tree root hash and keep it on its disk after uploading the files to the server and deleting its local copies. The client can request the i-th file Fi and a Merkle proof Pi for it from the server. The client uses the proof and compares the resulting root hash with the one it persisted before deleting the files - if they match, file is correct.

You can use any programming language you want (we use Go and Rust internally). We would like to see a solution with networking that can be deployed across multiple machines, and as close to production-ready as you have time for. Please describe the short-coming your solution have in a report, and how you would improve on them given more time.

We expect you to send us within 7 days:

- a demo of your app that we can try (ideally using eg Docker Compose)
- the code of the app
- a report (max 2-3 pages) explaining your approach, your other ideas, what went well or not, etc..

## Solution

For the solution I created this monorepo with the 3 codebases:

- [merkletree library](./libs/merkletree/) : is the shared library that implements the Merkle Tree logic.
- [server](./server/) : Is the server that stores the files and linked Merkle Tree.
- [client](./client/) : Is the client that bulk uploads files into the server, computes and stores their "root-hash" and then downloads one of them and checks its integrity with the locally stored "root-hash".

### Run

To run the demo these 2 steps are needed to spin-up the server and the client:

#### 1. Run Server

```bash
docker compose up
```

Check the Swagger API in [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

#### 2. Run Client

##### Build

```bash
cd client
task build
```

Then the two actions that the client does:

##### a) Bulk upload files from a folder

The client uploads a large set of potentially small files {F0, F1, …, Fn} to a server.

**_NOTE:_** By default the client bulk uploads the files from `./client/testfiles`

```bash
./bin/client-cli upload
```

##### b) Get a file from the server and checks its integrity

The client downloads an arbitrary file from the server and checks that the file is correct and is not corrupted in any way (in transport, tampered with by the server, etc.).

**_NOTE:_** By default the client downloads the file into `./client/downloads`

```bash
./bin/client-cli get f1
```
