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
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a file and check its validity",
	Long:  `Get a file and check its validity with previously created and stored merkle tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadFolder := viper.GetString("DEFAULT_DOWNLOAD_FOLDER")

		fileName := ""
		if len(args) == 1 && args[0] != "" {
			fileName = args[0]
		} else {
			fmt.Printf("insert the file name as only argument")
			return
		}

		userDownloadFolder, _ := cmd.Flags().GetString("dir")
		if userDownloadFolder != "" {
			if userDownloadFolder == viper.GetString("DEFAULT_DOWNLOAD_FOLDER") {
				if !utils.DirPathIsValid(downloadFolder) {
					err := os.MkdirAll(downloadFolder, os.ModePerm)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			var err error
			downloadFolder, err = utils.PathValidation(userDownloadFolder)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if downloadFolder == viper.GetString("DEFAULT_DOWNLOAD_FOLDER") {
			if !utils.DirPathIsValid(downloadFolder) {
				err := os.MkdirAll(downloadFolder, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		storageFolder := viper.GetString("DEFAULT_STORAGE_FOLDER")
		userStorageFolder, _ := cmd.Flags().GetString("store")
		fmt.Println("STORE:", userStorageFolder)
		if userStorageFolder != "" {
			var err error
			storageFolder, err = utils.PathValidation(userStorageFolder)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if storageFolder == viper.GetString("DEFAULT_STORAGE_FOLDER") {
			if !utils.DirPathIsValid(storageFolder) {
				err := os.MkdirAll(storageFolder, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		URL := viper.GetString("SERVER_URL") + "/files/" + fileName

		fmt.Println("Try to get '" + fileName + "' file...")

		// Get the data
		response, err := http.Get(URL)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {

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

			rootHashFile, err := os.Open(storageFolder + "/root-hash")
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

			// Verify file with merkle tree (library)
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
	getCmd.PersistentFlags().StringP("dir", "d", viper.GetString("DEFAULT_DOWNLOAD_FOLDER"), "output directory path where to store downloaded file")
	getCmd.PersistentFlags().StringP("store", "s", viper.GetString("DEFAULT_STORAGE_FOLDER"), "directory path where to find stored root-hash")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
