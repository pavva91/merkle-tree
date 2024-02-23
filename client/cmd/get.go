/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/spf13/cobra"
)

const DOWNLOAD_FOLDER = "./downloads"

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a file and check its validity",
	Long:  `Get a file and check its validity with previously created and stored merkle tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
		var fileName = "dr-who.png"

		if len(args) >= 1 && args[0] != "" {
			fileName = args[0]
		} else {
			// TODO: Return error and stop execution
		}

		URL := "http://localhost:8080/files/" + fileName

		fmt.Println("Try to get '" + fileName + "' file...")

		// Get the data
		response, err := http.Get(URL)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			err := os.RemoveAll(DOWNLOAD_FOLDER)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = os.MkdirAll(DOWNLOAD_FOLDER, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Create the file
			out, err := os.Create(fmt.Sprintf("%s/%s", DOWNLOAD_FOLDER, fileName))

			// out, err := os.Create(fileName)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			// Writer the body to file
			_, err = io.Copy(out, response.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Perfect! Just saved in %s/%s!\n", DOWNLOAD_FOLDER, fileName)

			// Get Merkle Proof
			mp := response.Header.Get("Merkle-Proof")
			// NOTE: To reconstruct string[] from mp:
			merkleProofs := strings.SplitAfter(strings.Replace(strings.Replace(strings.Replace(mp, "[", "", -1), "]", "", -1), "\"", "", -1), " ")
			for k, v := range merkleProofs {
				log.Printf("proof %v: %s", k, v)
			}

			// TODO: Verify merkle-proof with file hash and root hash
			f, err := os.Open(DOWNLOAD_FOLDER + "/" + fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				log.Fatal(err)
			}

			hashFile := fmt.Sprintf("%x", h.Sum(nil))
			fmt.Println("hash file:", hashFile)
			// hashFile = "0dffefeae189629164f222e18c83883c1fd9b5b02eb55d5ca99bd207ebcf882d"

			// TODO: Verify file with merkle tree (library)
			// merkletree.ComputeMerkleProof(f, )
			// TODO: get rootHash from ./storage
			rootHash := "5880895435d8c5d8c8b549b520ef550882ab0245e1b241594c44ddffe5a6a8c0"
			reconstructedRootHash := merkletree.ReconstructRootHash(hashFile, merkleProofs)
			isNotCorrupted := merkletree.Verify(hashFile, merkleProofs, rootHash)

			fmt.Println("wanted root hash:", rootHash)
			fmt.Println("reconstructed root hash:", reconstructedRootHash)
			fmt.Println("file is not corrupted:", isNotCorrupted)

		} else {
			fmt.Println("Error: " + fileName + " not exists! :-(")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
