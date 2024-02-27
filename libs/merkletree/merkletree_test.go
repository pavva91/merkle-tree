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
	mu3              sync.Mutex
	mu10             sync.Mutex
	merkletree3files = [][]string{
		{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
		{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
		{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
	}
	rootHash3files    = merkletree3files[len(merkletree3files)-1][0]
	merkletree10files = [][]string{
		{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "6eb9943ff3fc1d55ee9a540dee5f31eb692c2320edade9dd9d52a12ba22a8dde", "808164e8f59a67d651a047e4f32958cc6b46fd838d0a685981c9837ce1dbbd31", "41b8573aaad9e17eef146d8f3ee8469c63434245c61c5b11c8076c79c03ac709", "d0fed8662e301578f18b134fa794d34d5e49a70c66eddea0d0c1a86f73013e43", "5af4f245b10bd35897a6edccdd16ae3771287ee9cd489de59549948b42d4732d", "c58dd405fd2cd4e7db6c81c9623ac56e801cbeb35fc5eb2707bc448b38e789dc", "f50f283167ab9138b43be9e35ae634dbd9ab92576c7f70b91881ddb7bd92e2cd"},
		{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa", "c87ca38f0fed6f32156293098c32ac6bb72dcddc1b440721c7d41e397f1a9e73", "9d687df7ae9acaa946f85c2e2e4865131b4e1738a2ad7d924396a8a093bce9a7", "304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830", "304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830"},
		{"1675119930739a21e3023ec740151d182e3072fd4fff0d5767d734524b72c826", "09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728", "ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0", "ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0"},
		{"aaacf5f89d65da043a505c043d38e010237da1c3d012088be61c3bfb3b340f6b", "d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771"},
		{"78835624cd2ad609ce362709bbaf5cc02f82afe7c87c498429fc8e13b8d861b6"}, // root hash
	}
	rootHash10files = merkletree10files[len(merkletree10files)-1][0]
)

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
				{"a96fa9a8ec57366442ecb3d70cc3039b7107543c0bf197828f80a9091c29491f"},
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
			[][]string{
				{"h1", "h2", "h3", "h3"},
				{"a96fa9a8ec57366442ecb3d70cc3039b7107543c0bf197828f80a9091c29491f", "f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff"},
				{"7728ad37f4999f80c92c7d45b32b63faf5934cf7549176cef0b0d15ecf0447f8"},
			},
		},
		"3 string real hashes": {
			args{
				hashLeaves: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f2
					"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f3
				},
			},
			[][]string{
				{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
				{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
				{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
			},
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
				{"a96fa9a8ec57366442ecb3d70cc3039b7107543c0bf197828f80a9091c29491f", "011676a9c0ad9df8158321c9f291f0ef3ffcd66ad93f2583a94c29d960b25d8f", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
				{"6c43537d3e6822f2e9e1b3d748c7bed78672f619e0e2383ce3e7483202cc9712", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
				{"670aef0f2c3152520e90594f6d0c3e44487fead6a5b3a17ef09aebb365b4d72e"},
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

func Test_createMerkleProof(t *testing.T) {
	type args struct {
		hashFile   string
		merkleTree [][]string
	}
	tests := map[string]struct {
		args args
		want []string
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
				merkleTree: [][]string{
					{"h1", "h2", "h3", "h3"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff"},
					{"4279484b826df5de36382d7cf13be9a59ea62f7bc986d257c038d0bd9df207e2"},
				},

				hashFile: "h4",
			},
			[]string{},
		},
		"3 string hashes, found f1": {
			args{
				merkleTree: [][]string{
					{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
					{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
					{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
				},

				hashFile: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f2": {
			args{
				merkleTree: [][]string{
					{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
					{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
					{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
				},

				hashFile: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", // f2
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"3 string hashes, found f3": {
			args{
				merkleTree: [][]string{
					{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
					{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
					{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
				},

				hashFile: "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", // f3
			},
			[]string{
				"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
				"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
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
			if got := createMerkleProof(tt.args.hashFile, tt.args.merkleTree); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMerkleProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateHashPair(t *testing.T) {
	type args struct {
		h1 string
		h2 string
	}
	tests := map[string]struct {
		args args
		want string
	}{
		"test reverse (ascii code) order": {
			args{
				h1: "aaaa",
				h2: "bbbb",
			},
			"bbbbaaaa",
		},
		"test in correct (ascii code) order": {
			args{
				h1: "bbbb",
				h2: "aaaa",
			},
			"bbbbaaaa",
		},
		"test in order, real hash": {
			args{
				h1: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				h2: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
			},
			"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe60dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
		},
		"test reverse order, real hash": {
			args{
				h1: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				h2: "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
			},
			"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe60dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := calculateHashPair(tt.args.h1, tt.args.h2); got != tt.want {
				t.Errorf("calculateHashPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconstructRootHash(t *testing.T) {
	type args struct {
		hashFile     string
		merkleProofs []string
	}
	tests := map[string]struct {
		args args
		want string
	}{
		"test hash f1, with correct proofs, reconstruct correct rootHash": {
			args{
				hashFile: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
			},
			rootHash3files,
		},
		"test hash f1, with not correct proofs, reconstruct not correct rootHash": {
			args{
				hashFile: "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", // f1
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe7",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
			},
			"d7b975d9510021f16925d48e12ad209ad64c178e8c6f930a4ff67bffd1ac177e",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := reconstructRootHash(tt.args.hashFile, tt.args.merkleProofs); got != tt.want {
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
				wantedRootHash: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
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
				wantedRootHash: "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0",
			},
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isHashFileCorrect(tt.args.hashFile, tt.args.merkleProofs, tt.args.wantedRootHash); got != tt.want {
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
		"3 files": {
			args{
				files: fFiles,
			},
			[][]string{
				{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52"},
				{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286"},
				{"5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"},
			},
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
func TestComputeMerkleTree_10files(t *testing.T) {
	mu10.Lock()
	defer mu10.Unlock()
	files, err := os.ReadDir(test10FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test10FilesPath, f.Name())
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
		"10 string hashes": {
			args{
				files: fFiles,
			},
			[][]string{
				{"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d", "f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6", "34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52", "6eb9943ff3fc1d55ee9a540dee5f31eb692c2320edade9dd9d52a12ba22a8dde", "808164e8f59a67d651a047e4f32958cc6b46fd838d0a685981c9837ce1dbbd31", "41b8573aaad9e17eef146d8f3ee8469c63434245c61c5b11c8076c79c03ac709", "d0fed8662e301578f18b134fa794d34d5e49a70c66eddea0d0c1a86f73013e43", "5af4f245b10bd35897a6edccdd16ae3771287ee9cd489de59549948b42d4732d", "c58dd405fd2cd4e7db6c81c9623ac56e801cbeb35fc5eb2707bc448b38e789dc", "f50f283167ab9138b43be9e35ae634dbd9ab92576c7f70b91881ddb7bd92e2cd"},
				{"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc", "4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa", "c87ca38f0fed6f32156293098c32ac6bb72dcddc1b440721c7d41e397f1a9e73", "9d687df7ae9acaa946f85c2e2e4865131b4e1738a2ad7d924396a8a093bce9a7", "304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830", "304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830"},
				{"1675119930739a21e3023ec740151d182e3072fd4fff0d5767d734524b72c826", "09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728", "ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0", "ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0"},
				{"aaacf5f89d65da043a505c043d38e010237da1c3d012088be61c3bfb3b340f6b", "d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771"},
				{"78835624cd2ad609ce362709bbaf5cc02f82afe7c87c498429fc8e13b8d861b6"}, // root hash
			},
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
		"3 files": {
			args{
				files: fFiles,
			},
			rootHash3files,
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

func TestComputeRootHash_10files(t *testing.T) {
	mu10.Lock()
	defer mu10.Unlock()
	files, err := os.ReadDir(test10FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	var fFiles []*os.File
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test10FilesPath, f.Name())
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
		"10 files": {
			args{
				files: fFiles,
			},
			"78835624cd2ad609ce362709bbaf5cc02f82afe7c87c498429fc8e13b8d861b6",
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
				merkleTree: merkletree3files,
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"find second file (f2)": {
			args{
				file:       fFiles[1],
				merkleTree: merkletree3files,
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
			},
		},
		"find third file (f3)": {
			args{
				file:       fFiles[2],
				merkleTree: merkletree3files,
			},
			[]string{
				"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
				"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
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

func TestComputeMerkleProof_10files(t *testing.T) {
	mu10.Lock()
	defer mu10.Unlock()
	files, err := os.ReadDir(test10FilesPath)
	if err != nil {

		fmt.Println(err)
		return
	}

	fFiles := make(map[string]*os.File)
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test10FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {

			fmt.Println(err)
			return
		}
		fFiles[f.Name()] = ff
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
				file:       fFiles["f01"],
				merkleTree: merkletree10files,
			},
			[]string{
				"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
				"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
				"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
				"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
			},
		},
		"find second file (f2)": {
			args{
				file:       fFiles["f02"],
				merkleTree: merkletree10files,
			},
			[]string{
				"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
				"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
				"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
				"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
			},
		},
		"find third file (f3)": {
			args{
				file:       fFiles["f03"],
				merkleTree: merkletree10files,
			},
			[]string{
				"6eb9943ff3fc1d55ee9a540dee5f31eb692c2320edade9dd9d52a12ba22a8dde",
				"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
				"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
				"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
			},
		},
		"find seventh file (f7)": {
			args{
				file:       fFiles["f07"],
				merkleTree: merkletree10files,
			},
			[]string{
				"5af4f245b10bd35897a6edccdd16ae3771287ee9cd489de59549948b42d4732d",
				"c87ca38f0fed6f32156293098c32ac6bb72dcddc1b440721c7d41e397f1a9e73",
				"1675119930739a21e3023ec740151d182e3072fd4fff0d5767d734524b72c826",
				"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771"},
		},
		"find tenth file (f10)": {
			args{
				file:       fFiles["f10"],
				merkleTree: merkletree10files,
			},
			[]string{
				"c58dd405fd2cd4e7db6c81c9623ac56e801cbeb35fc5eb2707bc448b38e789dc",
				"304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830",
				"ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0",
				"aaacf5f89d65da043a505c043d38e010237da1c3d012088be61c3bfb3b340f6b"},
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
			},
			rootHash3files,
			false,
		},
		"test f2, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles[1],
				merkleProofs: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
			},
			rootHash3files,
			false,
		},
		"test f3, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles[2],
				merkleProofs: []string{
					"34575cdd0f12f999e0fc36ef7d70bbd5d302b9bca1a24a0712f505f490cf7a52",
					"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
				},
			},
			rootHash3files,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ReconstructRootHash(tt.args.file, tt.args.merkleProofs)
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

func TestReconstructRootHash_10files(t *testing.T) {
	mu10.Lock()
	defer mu10.Unlock()
	// defer
	files, err := os.ReadDir(test10FilesPath)
	if err != nil {

		fmt.Println(err)
		return
	}

	// var fFiles []*os.File
	fFiles := make(map[string]*os.File)
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", test10FilesPath, f.Name())
		ff, err := os.Open(filePath)
		if err != nil {

			fmt.Println(err)
			return
		}
		fFiles[f.Name()] = ff
		defer ff.Close()
	}

	type args struct {
		file         *os.File
		merkleProofs []string
	}
	tests := map[string]struct {
		args    args
		want    string
		wantErr bool
	}{
		"test f1, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles["f01"],
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
					"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
				},
			},
			rootHash10files,
			false,
		},
		"test f2, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles["f02"],
				merkleProofs: []string{
					"0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
					"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
					"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
				},
			},
			rootHash10files,
			false,
		},
		"test f3, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles["f03"],
				merkleProofs: []string{
					"6eb9943ff3fc1d55ee9a540dee5f31eb692c2320edade9dd9d52a12ba22a8dde",
					"26b28d79c60bda9bbec02d214d5defe3e21075276927239729cb2c01d9931acc",
					"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
				},
			},
			rootHash10files,
			false,
		},
		"test f7, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles["f07"],
				merkleProofs: []string{
					"5af4f245b10bd35897a6edccdd16ae3771287ee9cd489de59549948b42d4732d",
					"c87ca38f0fed6f32156293098c32ac6bb72dcddc1b440721c7d41e397f1a9e73",
					"1675119930739a21e3023ec740151d182e3072fd4fff0d5767d734524b72c826",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771"},
			},
			rootHash10files,
			false,
		},
		"test f10, with correct proofs, reconstruct correct rootHash": {
			args{
				file: fFiles["f10"],
				merkleProofs: []string{
					"c58dd405fd2cd4e7db6c81c9623ac56e801cbeb35fc5eb2707bc448b38e789dc",
					"304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830",
					"ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0",
					"aaacf5f89d65da043a505c043d38e010237da1c3d012088be61c3bfb3b340f6b",
				},
			},
			rootHash10files,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ReconstructRootHash(tt.args.file, tt.args.merkleProofs)
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
				wantedRootHash: rootHash3files,
			},
			true,
			false,
		},
		// FIX: I can't use fFiles["f1"] 2 times
		// works fine when running the test alone
		// It will fail sometimes with concurrent access to the file from other tests
		// "test f1, with not correct proofs, return false": {
		// 	args{
		// 		file: fFiles["f1"],
		// 		merkleProofs: []string{
		// 			"f9addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
		// 			"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
		// 		},
		// 		wantedRootHash: rootHash3files,
		// 	},
		// 	false,
		// 	false,
		// },
		"test f2, with not correct proofs, return false": {
			args{
				file: fFiles["f2"],
				merkleProofs: []string{
					"f9addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"dfa84bc707cd740d3551233bfda2cfa6df519d1e7e7174882efa7dc3cdab2286",
				},
				wantedRootHash: rootHash3files,
			},
			false,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := IsFileCorrect(tt.args.file, tt.args.merkleProofs, tt.args.wantedRootHash)
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

func TestIsFileCorrect_10files(t *testing.T) {
	mu10.Lock()
	defer mu10.Unlock()
	files, err := os.ReadDir(test10FilesPath)
	if err != nil {
		fmt.Println(err)

		return
	}

	fFiles := make(map[string]*os.File)
	for k, f := range files {
		filePath := fmt.Sprintf("%s/%s", test10FilesPath, f.Name())
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
	}
	tests := map[string]struct {
		args    args
		want    bool
		wantErr bool
	}{
		"test f1, with correct proofs, return true": {
			args{
				file: fFiles["f01"],
				merkleProofs: []string{
					"f8addeff4cc29a9a55589ae001e2230ecd7a515de5be7eeb27da1cabba87fbe6",
					"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
					"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
				},

				wantedRootHash: rootHash10files,
			},
			true,
			false,
		},
		"test f10, with correct proofs, return true": {
			args{
				file: fFiles["f10"],
				merkleProofs: []string{
					"c58dd405fd2cd4e7db6c81c9623ac56e801cbeb35fc5eb2707bc448b38e789dc",
					"304407ffd72159a2ba0e6fcd5d5f4abe581bb34b16a1ef5b854e66d44502c830",
					"ab0a11a9192d5ded7e65da7a6de3aa7cbab0ca4dd3a895e02ccd2ba05d1276f0",
					"aaacf5f89d65da043a505c043d38e010237da1c3d012088be61c3bfb3b340f6b",
				},
				wantedRootHash: rootHash10files,
			},
			true,
			false,
		},
		"test f2, with not correct proofs, return false": {
			args{
				file: fFiles["f02"],
				merkleProofs: []string{
					"1dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d",
					"4a7493d6c9fab2b2a0cc82ef847d1e9142c7195b2f340fc4e67bc8b10aa808fa",
					"09c294d0970eab9be9462e98dcda4c9094f67fd583fb50e1d47950dccacc5728",
					"d264659f4944fa857ebdf54b66c9e445ca9a364c0a749a9fbfe873871e229771",
				},

				wantedRootHash: rootHash10files,
			},
			false,
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := IsFileCorrect(tt.args.file, tt.args.merkleProofs, tt.args.wantedRootHash)
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
