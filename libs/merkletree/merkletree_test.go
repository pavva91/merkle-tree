package merkletree

import (
	"os"
	"reflect"
	"testing"
)

func TestComputeRootHash(t *testing.T) {
	type args struct {
		files []*os.File
	}
	tests := map[string]struct {
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ComputeRootHash(tt.args.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeRootHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComputeRootHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMerkleTreeAsMatrix(t *testing.T) {
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
			[][]string{
				{"h1", "h2", "h3", "h3"},
				{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff"},
				{"4279484b826df5de36382d7cf13be9a59ea62f7bc986d257c038d0bd9df207e2"},
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
				{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "55c82a1f310283eefe23c4e02d409428fb0e768551eb4845291ed67ac2b16ec3", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd", "77c1e48fb22fdbd534fd76bc0b8fa98745e113634512b7171ceaae33b097e6fd"},
				{"c14124471a06847b5042b48aa94ece8030e2a21fbcf2927e2741ef2602f37363", "3a31f2d07d9715bc3e48106ecbdcb02e2081e96d3537a47b54975a0c52be66b7"},
				{"12761d3647c296c8a6e39bb363652479da1e95382128ba28c6eb0e79ee74a97a"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := createMerkleTreeAsMatrix(tt.args.hashLeaves); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
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
		"3 string hashes, found": {
			args{
				merkleTree: [][]string{
					{"h1", "h2", "h3", "h3"},
					{"dac079ce8e97c5434424c28112b96e601aa4ff36ba0377619b9e38f473310cf3", "f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff"},
					{"4279484b826df5de36382d7cf13be9a59ea62f7bc986d257c038d0bd9df207e2"},
				},

				hashFile: "h2",
			},
			[]string{
				"h1",
				"f858727943465ed9759534a90713a9da630425509eaf13633d9c3229434490ff",
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
