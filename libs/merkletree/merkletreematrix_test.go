package merkletree

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
)

const (
	test3FilesPath  = "./testfiles/3files"
	test10FilesPath = "./testfiles/10files"
)

var (
	mu3  sync.Mutex
	mu10 sync.Mutex
)

var merkleTree3HashesM = [][]string{
	{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
	{"e2abf2fd16b981a59d8ecf4a9a0ac0498c715e45801e29f0152cefad8c6f87f4", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
	{hashRoot3Hashes},
}

var hashRoot3Hashes = "b2a1f9e0a30ae91c0b6b70eacb673e4e030f3d5199ec0a0f0ed64ad45c0ca7f4"

var merkleTree3StringsM = [][]string{
	{"h1", "h2", "h3", "h3"},
	{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff"},
	{hashRoot3Strings},
}

var hashRoot3Strings = "4279484b826df5de36382d7cf13be9a59ea62f7bc986d257c038d0bd9df207e2"

func Test_createMerkleTree(t *testing.T) {
	type args struct {
		hashLeaves []string
	}
	tests := map[string]struct {
		args args
		want [][]string
	}{
		"2 string hashes": {
			args{
				hashLeaves: []string{
					"h1",
					"h2",
				},
			},
			[][]string{
				{"h1", "h2"},
				{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3"},
			},
		},
		"3 string hashes": {
			args{
				hashLeaves: []string{
					"h1",
					"h2",
					"h3",
				},
			},
			merkleTree3StringsM,
		},
		"3 string real hashes": {
			args{
				hashLeaves: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f2
					"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f3
				},
			},
			merkleTree3HashesM,
		},
		"5 string hashes": {
			args{
				hashLeaves: []string{
					"h1",
					"h2",
					"h3",
					"h4",
					"h5",
				},
			},
			[][]string{
				{"h1", "h2", "h3", "h4", "h5", "h5"},
				{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
				{"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
				{"12761d3647c296c8a6e39bb363652479da1e95382128ba28c6eb0e79ee74a97a"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := createMerkleTree(tt.args.hashLeaves); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMerkleTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMerkleProofM(t *testing.T) {
	type args struct {
		hashFile   string
		merkleTree [][]string
	}
	tests := map[string]struct {
		args             args
		wantMerkleProofs []string
	}{
		"2 string hashes, not found": {
			args{
				merkleTree: [][]string{
					{"h1", "h2"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3"},
				},
				hashFile: "h3",
			},
			[]string{},
		},
		"2 string hashes, found": {
			args{
				merkleTree: [][]string{
					{"h1", "h2"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3"},
				},
				hashFile: "h2",
			},
			[]string{
				"h1",
			},
		},
		"3 string hashes, not found": {
			args{
				merkleTree: merkleTree3HashesM,

				hashFile: "h4",
			},
			[]string{},
		},
		"3 string hashes, found f1": {
			args{
				merkleTree: merkleTree3HashesM,
				hashFile:   "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f2": {
			args{
				merkleTree: merkleTree3HashesM,
				hashFile:   "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f2
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f3": {
			args{
				merkleTree: merkleTree3HashesM,
				hashFile:   "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f3
			},
			[]string{
				"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
				"e2abf2fd16b981a59d8ecf4a9a0ac0498c715e45801e29f0152cefad8c6f87f4",
			},
		},
		"5 string hashes, found h2": {
			args{
				merkleTree: [][]string{
					{"h1", "h2", "h3", "h4", "h5", "h5"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
					{"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
					{"12761d3647c296c8a6e39bb363652479da1e95382128ba28c6eb0e79ee74a97a"},
				},

				hashFile: "h2",
			},
			[]string{
				"h1",
				"55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3",
				"3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7",
			},
		},
		"5 string hashes, found h5": {
			args{
				merkleTree: [][]string{
					{"h1", "h2", "h3", "h4", "h5", "h5"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
					{"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
					{"12761d3647c296c8a6e39bb363652479da1e95382128ba28c6eb0e79ee74a97a"},
				},

				hashFile: "h5",
			},
			[]string{
				"h5",
				"77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd",
				"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363",
			},
		},
		"5 string hashes, not found h6": {
			args{
				merkleTree: [][]string{
					{"h1", "h2", "h3", "h4", "h5", "h5"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
					{"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
					{"12761d3647c296c8a6e39bb363652479da1e95382128ba28c6eb0e79ee74a97a"},
				},

				hashFile: "h6",
			},
			[]string{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotMerkleProofs := createMerkleProofMatrix(tt.args.hashFile, tt.args.merkleTree)

			if len(gotMerkleProofs) != len(tt.wantMerkleProofs) {
				t.Errorf("createMerkleProof(): len = %v, want %v", len(gotMerkleProofs), len(tt.wantMerkleProofs))
			}

			for k, v := range gotMerkleProofs {
				if v != tt.wantMerkleProofs[k] {
					t.Errorf("createMerkleProof(): element %v = %v, want %v", k, v, tt.wantMerkleProofs[k])
				}
			}
		})

		// FIX: Why reflect.DeepEqual is now giving error
		// 	if got := createMerkleProofMatrix(tt.args.hashFile, tt.args.merkleTree); !reflect.DeepEqual(got, tt.want) {
		// 		t.Errorf("createMerkleProof() = %v, want %v", got, tt.want)
		// 	}
		// })
	}
}

// 0 : right
// 1 : left
func Test_reconstructRootHash(t *testing.T) {
	type args struct {
		hashLeaf     string
		merkleProofs []string
		fileOrder    int
	}
	tests := map[string]struct {
		args args
		want string
	}{
		"test hash f1, (3 string hashes), with correct proofs, reconstruct correct rootHash": {
			args{
				hashLeaf: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				fileOrder: 0,
			},
			hashRoot3Hashes,
		},
		"test hash f1, with not correct proofs, reconstruct not correct rootHash": {
			args{
				hashLeaf: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe7",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				fileOrder: 0,
			},
			"2a6820a743f3d68416e37cc1d0088b39556ad4c63b23a2e7044a7eea287e5ee5",
		},
		"test f2 (3 string hashes), with correct proofs, reconstruct correct rootHash": {
			args{
				hashLeaf: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f1
				merkleProofs: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				fileOrder: 1,
			},
			hashRoot3Hashes,
		},
		"test f3 (3 string hashes), with correct proofs, reconstruct correct rootHash": {
			args{
				hashLeaf: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f1
				merkleProofs: []string{
					"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					"e2abf2fd16b981a59d8ecf4a9a0ac0498c715e45801e29f0152cefad8c6f87f4",
				},
				fileOrder: 2,
			},
			hashRoot3Hashes,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := reconstructRootHash(tt.args.hashLeaf, tt.args.merkleProofs, tt.args.fileOrder); got != tt.want {
				t.Errorf("ReconstructRootHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHashFileCorrect(t *testing.T) {
	type args struct {
		hashFile       string
		merkleProofs   []string
		wantedRootHash string
		fileOrder      int
	}
	tests := map[string]struct {
		args args
		want bool
	}{
		"test f1, with correct proofs, return true": {
			args{
				hashFile: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				wantedRootHash: hashRoot3Hashes,
				fileOrder:      0,
			},
			true,
		},
		"test f1, with not correct proofs, return false": {
			args{
				hashFile: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe7",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				wantedRootHash: hashRoot3Hashes,
				fileOrder:      0,
			},
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isHashFileCorrect(tt.args.hashFile, tt.args.merkleProofs, tt.args.wantedRootHash, tt.args.fileOrder); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeMerkleTree_3files(t *testing.T) {
	mu3.Lock()
	defer mu3.Unlock()
	files, err := os.ReadDir(test3FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test3FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)

			return
		}
		fFiles = append(fFiles, ff)
		defer ff.Close()
	}
	//

	type args struct {
		files []*os.File
	}
	tests := map[string]struct {
		args    args
		want    [][]string
		wantErr bool
	}{
		"3 string hashes": {
			args{
				files: fFiles,
			},
			merkleTree3HashesM,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ComputeMerkleTree(tt.args.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeMerkleTree() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeMerkleTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeRootHash_3files(t *testing.T) {
	mu3.Lock()
	defer mu3.Unlock()
	files, err := os.ReadDir(test3FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test3FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fFiles = append(fFiles, ff)
		defer ff.Close()
	}
	//

	type args struct {
		files []*os.File
	}
	tests := map[string]struct {
		args    args
		want    string
		wantErr bool
	}{
		"3 string hashes": {
			args{
				files: fFiles,
			},
			hashRoot3Hashes,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ComputeRootHash(tt.args.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeMerkleTree() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeMerkleTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeMerkleProof_3files(t *testing.T) {
	mu3.Lock()
	defer mu3.Unlock()

	files, err := os.ReadDir(test3FilesPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test3FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fFiles = append(fFiles, ff)
		defer ff.Close()
	}

	type args struct {
		file       *os.File
		merkleTree [][]string
	}
	tests := map[string]struct {
		args args
		want []string
	}{
		"find first file (f1)": {
			args{
				file:       fFiles[0],
				merkleTree: merkleTree3HashesM,
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"find second file (f2)": {
			args{
				file:       fFiles[1],
				merkleTree: merkleTree3HashesM,
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"find third file (f3)": {
			args{
				file:       fFiles[2],
				merkleTree: merkleTree3HashesM,
			},
			[]string{
				"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
				"e2abf2fd16b981a59d8ecf4a9a0ac0498c715e45801e29f0152cefad8c6f87f4",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := ComputeMerkleProof(tt.args.file, tt.args.merkleTree); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeMerkleProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconstructRootHash_3files(t *testing.T) {
	mu3.Lock()
	defer mu3.Unlock()
	// defer
	files, err := os.ReadDir(test3FilesPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test3FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fFiles = append(fFiles, ff)
		defer ff.Close()
	}

	type args struct {
		file         *os.File
		merkleProofs []string
		fileOrder    int
	}
	tests := map[string]struct {
		args    args
		want    string
		wantErr bool
	}{
		"test f1, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles[0],
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				fileOrder: 0,
			},
			hashRoot3Hashes,
			false,
		},
		"test f2, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles[1],
				merkleProofs: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				fileOrder: 1,
			},
			hashRoot3Hashes,
			false,
		},
		"test f3, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles[2],
				merkleProofs: []string{
					"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					"e2abf2fd16b981a59d8ecf4a9a0ac0498c715e45801e29f0152cefad8c6f87f4",
				},
				fileOrder: 2,
			},
			hashRoot3Hashes,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ReconstructRootHash(tt.args.file, tt.args.merkleProofs, tt.args.fileOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconstructRootHash() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("ReconstructRootHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileCorrect_3files(t *testing.T) {
	// TODO: breaks with concurrent run
	mu3.Lock()
	defer mu3.Unlock()
	files, err := os.ReadDir(test3FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	fFiles := make(map[string]*os.File)
	for k, f := range files {
		filePath := fmt.Sprintf("%s/%s", test3FilesPath, f.Name())
		fmt.Printf("filepath %v: %s\n", k, filePath)
		ff, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)

			return
		}
		// fFiles = append(fFiles, ff)
		fFiles[f.Name()] = ff
		defer ff.Close()
	}
	//

	type args struct {
		file           *os.File
		merkleProofs   []string
		wantedRootHash string
		fileOrder      int
	}
	tests := map[string]struct {
		args    args
		want    bool
		wantErr bool
	}{
		"test f1, with correct proofs, return true": {
			args{
				file: fFiles["f1"],
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				wantedRootHash: hashRoot3Hashes,
				fileOrder:      0,
				// FIX:
				// scenario error 1
				// sometimes it gets: f814d46d6b2225d6188a8483b3ba7f53f903a911442e4a6eefb6514fd0afa7db
				// hash 0 is: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855

				// FIX:
				// scenario error 2
				// sometimes root-hash gets: 649e75ba58832108c03e2cd841d2e43b722759c4a644ea0bd5e8a31b42be791e
				// hash 0 is: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
			},
			true,
			false,
		},
		// FIX:
		"test f1, with not correct proofs, return false": {
			args{
				file: fFiles["f2"],
				merkleProofs: []string{
					"f9addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				wantedRootHash: hashRoot3Hashes,
				fileOrder:      0,
			},
			false,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := IsFileCorrect(tt.args.file, tt.args.merkleProofs, tt.args.wantedRootHash, tt.args.fileOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFileCorrect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsFileCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}
