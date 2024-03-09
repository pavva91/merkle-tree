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
	"strconv"
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
		fileOrderStr, _ := cmd.Flags().GetString("order")
		if fileOrderStr == "unassigned" {
			fmt.Printf("insert the file order, starting from one")
			return
		}
		fileOrder, err := strconv.Atoi(fileOrderStr)
		if err != nil {
			fmt.Printf("file order must be an integer")
			return
		}
		fileOrder--

		downloadFolder := viper.GetString("DEFAULT_DOWNLOAD_FOLDER")

		fileName := ""
		if len(args) == 1 && args[0] != "" {
			fileName = args[0]
		} else {
			fmt.Printf("insert the file name as only argument, starting from one")
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
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			fmt.Println(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		if res.StatusCode == 200 {

			// Create the file
			out, err := os.Create(fmt.Sprintf("%s/%s", downloadFolder, fileName))
			// out, err := os.Create(fileName)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			// Writer the body to file
			_, err = io.Copy(out, res.Body)
			if err != nil {
				fmt.Println(err)
			}

			// Get Merkle Proof
			mp := res.Header.Get("Merkle-Proof")
			// FIX:wrapperFunc: use strings.ReplaceAll method in `strings.Replace(mp, "[", "", -1)` (gocritic)
			merkleProofs := strings.SplitAfter(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(mp, "[", ""), "]", ""), "\"", ""), " ")

			for k, v := range merkleProofs {
				log.Printf("proof %v: %s", k, v)
			}

			fmt.Println("---------------------------------------------------------------------------------------")

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
			reconstructedRootHash, err := merkletree.ReconstructRootHash(file1, merkleProofs, fileOrder)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("wanted root hash:", rootHash)
			fmt.Println("reconstructed root hash:", reconstructedRootHash)

			isFileCorrect, err := merkletree.IsFileCorrect(file2, merkleProofs, rootHash, fileOrder)
			if err != nil {
				fmt.Println(err)
				return
			}

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
	getCmd.PersistentFlags().StringP("order", "o", "unassigned", "order of the file inside the merkle tree")
	getCmd.PersistentFlags().StringP("dir", "d", viper.GetString("DEFAULT_DOWNLOAD_FOLDER"), "output directory path where to store downloaded file")
	getCmd.PersistentFlags().StringP("store", "s", viper.GetString("DEFAULT_STORAGE_FOLDER"), "directory path where to find stored root-hash")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
