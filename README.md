# Coding Challenge

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

- ./libs/merkletree/ : is the shared library that implements the Merkle Tree logic.
- ./server/ : Is the server that stores the files and linked Merkle Tree.
- ./client/ : Is the client that bulk uploads files, computes and stores their "root-hash" and then downloads one of them and checks its integrity with the locally stored "root-hash".

### Run

#### Run Server

```bash
docker compose up
```

#### Run Client

##### Build

```bash
cd client
task build
```

##### 1) Bulk upload files from a folder (by default from ./client/testfiles/)

```bash
./bin/client-cli upload
```

##### 2) Get a file from the server and checks its integrity

```bash
./bin/client-cli get f1
```
