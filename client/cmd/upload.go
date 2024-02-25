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

	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/spf13/cobra"
)

// TODO: Idiomatic Go for constants
const (
	DEFAULT_STORAGE_FOLDER = "./storage"
	DEFAULT_UPLOAD_FOLDER  = "./testfiles"
	DEFAULT_SERVER_URL     = "http://localhost:8080/files"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	// TODO: client can choose server url
	// TODO: client can choose directory where to pick files to upload
	// TODO: client can choose directory where to store root hash
	Use:   "upload",
	Short: "Bulk Upload all files in a folder",
	Long: `Bulk upload all files inside a folder passed as input. 
		The function will also calculate the merkle tree of the files and store the "root-hash" of the merkle tree`,
	Run: func(cmd *cobra.Command, args []string) {
		// uploadFolder := "~/work/zama/client/testfiles"
		uploadFolder := "./testfiles"
		if len(args) >= 1 && args[0] != "" {
			uploadFolder = args[0]
		}

		// TODO: Check and Remove trailing "/"

		url := DEFAULT_SERVER_URL
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
			// TODO: Check if it is a file or a directory
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

		rootHash, err := merkletree.ComputeRootHash(rFiles...)

		err = os.RemoveAll(DEFAULT_STORAGE_FOLDER)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.MkdirAll(DEFAULT_STORAGE_FOLDER, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: RW mutex to have thread-safe access to root-hash
		rootHashPath := fmt.Sprintf("%s/%s", DEFAULT_STORAGE_FOLDER, "root-hash")
		err = os.WriteFile(rootHashPath, []byte(rootHash), 0666)
		fmt.Println("root hash stored in:", rootHashPath)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
