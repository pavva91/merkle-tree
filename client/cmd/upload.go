/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pavva91/merkle-tree/client/internal/utils"
	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Bulk Upload all files in a folder",
	Long: `Bulk upload all files inside a folder passed as input. 
	The function will also calculate the merkle tree of the files and store the "root-hash" of the merkle tree`,
	Run: func(cmd *cobra.Command, _ []string) {
		uploadFolder := viper.GetString("DEFAULT_UPLOAD_FOLDER")
		userUploadFolder, _ := cmd.Flags().GetString("dir")
		if userUploadFolder != "" {
			var err error
			uploadFolder, err = utils.PathValidation(userUploadFolder)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		storageFolder := viper.GetString("DEFAULT_STORAGE_FOLDER")
		userStorageFolder, _ := cmd.Flags().GetString("store")
		if userStorageFolder != "" {
			if userStorageFolder == viper.Get("DEFAULT_STORAGE_FOLDER") {
				if !utils.DirPathIsValid(storageFolder) {
					err := os.MkdirAll(storageFolder, os.ModePerm)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			var err error
			storageFolder, err = utils.PathValidation(userStorageFolder)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if storageFolder == viper.Get("DEFAULT_STORAGE_FOLDER") {
			if !utils.DirPathIsValid(storageFolder) {
				err := os.MkdirAll(storageFolder, os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		url := viper.GetString("SERVER_URL") + "/files"
		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)

		// cycle all files in the folder
		files, err := os.ReadDir(uploadFolder)
		if err != nil {
			fmt.Println(err)
			return
		}

		var rFiles []*os.File

		for _, f := range files {
			if f.IsDir() {
				fmt.Println("use a folder with only the files you want to upload, without subfolders")
				return
			}
		}

		for _, f := range files {
			filePath := uploadFolder + "/" + f.Name()
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			file2, err := os.Open(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			rFiles = append(rFiles, file2)
			defer file2.Close()

			part, err := writer.CreateFormFile("file", filepath.Base(filePath))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = io.Copy(part, file)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "multipart/form-data")

		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("server response:", string(body))
		fmt.Println("---------------------------------------------------------------------------------------")

		rootHash, err := merkletree.ComputeRootHash(rFiles...)
		if err != nil {
			fmt.Println(err)
			return
		}

		rootHashPath := fmt.Sprintf("%s/%s", storageFolder, "root-hash")
		err = os.WriteFile(rootHashPath, []byte(rootHash), 0o600)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("root hash stored in:", rootHashPath)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")
	uploadCmd.PersistentFlags().StringP("dir", "d", viper.GetString("DEFAULT_UPLOAD_FOLDER"), "directory path where to bulk upload files")
	uploadCmd.PersistentFlags().StringP("store", "s", viper.GetString("DEFAULT_STORAGE_FOLDER"), "output directory path where to store root-hash")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
