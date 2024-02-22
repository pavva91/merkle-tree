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

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
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

		url := "http://localhost:8080/files"
		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)

		// cycle all files in the folder
		files, err := os.ReadDir(uploadFolder)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, f := range files {
			filePath := uploadFolder + "/" + f.Name()
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
			_, err = io.Copy(part1, file)
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
