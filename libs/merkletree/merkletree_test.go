package merkletree

import (
	"testing"
)

func Test_calcMT(t *testing.T) {
	type args struct {
		hashNodes []*BinaryNode
	}
	tests := map[string]struct {
		args args
		want []*BinaryNode
	}{
		// TODO: Add test cases.
		"2 string hashes": {
			args{
				hashNodes: []*BinaryNode{
					{
						Value: "h1",
					},
					{
						Value: "h2",
					},
				},
			},
			[]*BinaryNode{
				{
					Value: "a96fa9a8ec57366442ecb3d70cc3039b7107543c0bf197828f80a9091c29491f",
					LeftNode: &BinaryNode{
						Value: "h1",
					},
					RightNode: &BinaryNode{
						Value: "h2",
					},
				},
			},
		},
		"3 string real hashes": {
			args{
				hashNodes: []*BinaryNode{
					{
						Value: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
					},
					{
						Value: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f2
					},
					{
						Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f3
					},
				},
			},
			[]*BinaryNode{
				{
					Value: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
					LeftNode: &BinaryNode{
						Value: "26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
						LeftNode: &BinaryNode{
							Value: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
						},
						RightNode: &BinaryNode{
							Value: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
						},
					},
					RightNode: &BinaryNode{
						Value: "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
						LeftNode: &BinaryNode{
							Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
						},
						RightNode: &BinaryNode{
							Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
						},
					},
				},
				// {"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
				// {"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
				// {"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := calcMT(tt.args.hashNodes)
			if len(got) != 1 {
				t.Errorf("size got = %v, want %v", len(got), 1)
			}
			if got[0].Value != tt.want[0].Value {
				t.Errorf("calcMT() = %v, want %v", got[0].Value, tt.want[0].Value)
			}
		})
	}
}

func Test_createMerkleProof(t *testing.T) {
	type args struct {
		hashLeaf   string
		merkleTree MerkleTree
	}
	tests := map[string]struct {
		args             args
		wantMerkleProofs []string
	}{
		// TODO: Add test cases.

		"2 string hashes, not found": {
			args{
				merkleTree: MerkleTree{
					RootHashNode: &BinaryNode{
						Value: "dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3",
						LeftNode: &BinaryNode{
							Value: "h1",
						},
						RightNode: &BinaryNode{
							Value: "h2",
						},
					},
					HashLeaves: []string{
						"h1",
						"h2",
					},
				},
				hashLeaf: "h3",
			},
			[]string{},
		},
		"2 string hashes, found h2": {
			args{
				merkleTree: MerkleTree{
					RootHashNode: &BinaryNode{
						Value: "dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3",
						LeftNode: &BinaryNode{
							Value: "h1",
						},
						RightNode: &BinaryNode{
							Value: "h2",
						},
					},
					HashLeaves: []string{
						"h1",
						"h2",
					},
				},
				hashLeaf: "h2",
			},
			[]string{
				"h1",
			},
		},

		"3 string hashes, found f1": {
			args{
				merkleTree: MerkleTree{
					RootHashNode: &BinaryNode{
						Value: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
						LeftNode: &BinaryNode{
							Value: "26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
							LeftNode: &BinaryNode{
								Value: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
							},
							RightNode: &BinaryNode{
								Value: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
							},
						},
						RightNode: &BinaryNode{
							Value: "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
							LeftNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
							RightNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
						},
					},
					HashLeaves: []string{
						"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
						"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
						"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					},
				},
				hashLeaf: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f2": {
			args{
				merkleTree: MerkleTree{
					RootHashNode: &BinaryNode{
						Value: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
						LeftNode: &BinaryNode{
							Value: "26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
							LeftNode: &BinaryNode{
								Value: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
							},
							RightNode: &BinaryNode{
								Value: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
							},
						},
						RightNode: &BinaryNode{
							Value: "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
							LeftNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
							RightNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
						},
					},
					HashLeaves: []string{
						"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
						"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
						"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					},
				},
				hashLeaf: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f3": {
			args{
				merkleTree: MerkleTree{
					RootHashNode: &BinaryNode{
						Value: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
						LeftNode: &BinaryNode{
							Value: "26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
							LeftNode: &BinaryNode{
								Value: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
							},
							RightNode: &BinaryNode{
								Value: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
							},
						},
						RightNode: &BinaryNode{
							Value: "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
							LeftNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
							RightNode: &BinaryNode{
								Value: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
							},
						},
					},
					HashLeaves: []string{
						"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
						"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
						"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					},
				},
				hashLeaf: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
			},
			[]string{
				"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
				"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotMerkleProofs := createMerkleProof(tt.args.hashLeaf, tt.args.merkleTree)

			if len(gotMerkleProofs) != len(tt.wantMerkleProofs) {
				t.Errorf("createMerkleProof(): len = %v, want %v", len(gotMerkleProofs), len(tt.wantMerkleProofs))
			}

			for k, v := range gotMerkleProofs {
				if v != tt.wantMerkleProofs[k] {
					t.Errorf("createMerkleProof(): element %v = %v, want %v", k, v, tt.wantMerkleProofs[k])
				}
			}
		})
	}
}
