# Vocabulary
A Merkle tree is a non-linear, binary, hash tree-like data structure.
Each leaf node of the tree stores the hash value of a data element.
Calculate Merkle Tree 4 files:

- root = hash(hash(hash(a) + hash(b)), hash(hash(c) + hash(d)))
  leaf node = file
  inner node
  top hash = root hash = master hash

Demonstrating that a leaf node is a part of a given binary hash tree requires computing a number of hashes proportional to the logarithm of the number of leaf nodes in the tree

A Merkle tree is therefore an efficient example of a cryptographic commitment scheme, in which the root of the tree is seen as a commitment and leaf nodes may be revealed and proven to be part of the original commitment

Hash trees can be used to verify any kind of data stored, handled and transferred in and between computers.

The main difference from a hash list is that one branch of the hash tree can be downloaded at a time and the integrity of each branch can be checked immediately, even though the whole tree is not available yet.

## Odd merkle tree

https://bitcoin.stackexchange.com/questions/46767/merkle-tree-structure-for-9-transactions
https://bitcoin.stackexchange.com/questions/79364/are-number-of-transactions-in-merkle-tree-always-even
Merkle Tree := Is a binary tree where each node is a hash over its child nodes.
If there are an odd number of nodes on any level of the merkle tree, the last node is duplicated and hashed with itself.
If there were a Tx4 (fith element), the diagram would look like this:

```bash
                  Root (Hash01234444)
               /                      \
        Hash0123                 Hash4444
        /      \                    /     \
   Hash01      Hash23           Hash44     Hash44
   /   \        /   \             /   \
Hash0  Hash1  Hash2  Hash3     Hash4   Hash4
```
```bash
                            ROOT
                            /  \   
                         /        \  
                      /              \
                   /                    \
               z1                            z2
             /  \                           / \
           /      \                       /     \
        /            \                /            \
      y1              y2             y3             y3
     /  \            /  \           / \
   /      \        /      \       /     \
  x1      x2      x3      x4     x5    x5
 / \     / \     / \     / \     / \ 
a   b   c   d   e   f   g   h   i   i
```

The merkle root (root hash in the diagram) is stored. That is the only hash stored. It is a fixed size, regardless of the rest of the merkle tree. The hashes in the merkle tree are not stored.
