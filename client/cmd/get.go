/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pavva91/merkle-tree/client/internal/utils"
	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/spf13/cobra"
)

// TODO: Idiomatic Go for constants
const DEFAULT_DOWNLOAD_FOLDER = "./downloads"

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a file and check its validity",
	Long:  `Get a file and check its validity with previously created and stored merkle tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadFolder := DEFAULT_DOWNLOAD_FOLDER

		userDownloadFolder, _ := cmd.Flags().GetString("dir")
		if userDownloadFolder != "" {
			home, _ := os.UserHomeDir()
			fmt.Println("HOME:", home)

			firstChar := userDownloadFolder[:1]
			if firstChar != "." {
				containsHome := strings.Contains(userDownloadFolder, home)
				if !containsHome {
					fmt.Println("folder with absolute path must be inside home")
					return
				}
			}
			if firstChar == "~" {
				userDownloadFolder = home + userDownloadFolder[1:]
			}
			if !utils.DirPathIsValid(userDownloadFolder) {
				fmt.Printf("folder %s does not exist", userDownloadFolder)
				return
			}
			downloadFolder = userDownloadFolder
		}

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
			err := os.RemoveAll(downloadFolder)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = os.MkdirAll(downloadFolder, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Create the file
			out, err := os.Create(fmt.Sprintf("%s/%s", downloadFolder, fileName))

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

			// Get Merkle Proof
			mp := response.Header.Get("Merkle-Proof")
			// NOTE: To reconstruct string[] from mp:
			merkleProofs := strings.SplitAfter(strings.Replace(strings.Replace(strings.Replace(mp, "[", "", -1), "]", "", -1), "\"", "", -1), " ")
			for k, v := range merkleProofs {
				log.Printf("proof %v: %s", k, v)
			}

			// TODO: Use RWMutex for root hash
			// TODO: Verify merkle-proof with file1 hash and root hash
			file1, err := os.Open(downloadFolder + "/" + fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file1.Close()

			file2, err := os.Open(downloadFolder + "/" + fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file2.Close()

			rootHashFile, err := os.Open(DEFAULT_STORAGE_FOLDER + "/root-hash")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rootHashFile.Close()

			rootHashBytes, err := io.ReadAll(rootHashFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			// print human readable
			rootHash := fmt.Sprint(string(rootHashBytes))
			fmt.Println("root hash retrieved:", rootHash)

			//NOTE: Verify file with merkle tree (library)
			reconstructedRootHash, err := merkletree.ReconstructRootHash(file1, merkleProofs)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("wanted root hash:", rootHash)
			fmt.Println("reconstructed root hash:", reconstructedRootHash)

			isFileCorrect, err := merkletree.IsFileCorrect(file2, merkleProofs, rootHash)
			if err != nil {
				fmt.Println(err)
				return
			}

			iscorrect := rootHash == reconstructedRootHash
			fmt.Println("IS CORRECT:", iscorrect)

			if isFileCorrect {
				fmt.Printf("file %s is not corrupted\n", fileName)
				fmt.Printf("downloaded file saved in %s/%s!\n", downloadFolder, fileName)

			} else {
				fmt.Printf("file %s is corrupted\n", fileName)
				err := os.Remove(downloadFolder + "/" + fileName)
				if err != nil {
					fmt.Println("Error during removal of " + fileName)
					return
				}
			}

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
	getCmd.PersistentFlags().String("dir", "", "Output directory path where to store downloaded file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
