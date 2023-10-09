/*
Copyright © 2023 Muhammad Saad <github.5s32y@slmail.me>
This file is part of CLI application pix.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mxsaad/pix/internal"
	"github.com/spf13/cobra"
)

var (
	outputDirectory string
	replaceOriginal bool
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [flags] <format> <images/directories...>",
	Short: "Convert images to a specified format",
	Args:  cobra.MinimumNArgs(2), // At least two arguments are required (format and one or more image files/directories)
	Run: func(cmd *cobra.Command, args []string) {
		format := strings.ToLower(args[0])
		files := args[1:]

		for _, file := range files {
			err := processFile(file, format)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", file, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&outputDirectory, "output-directory", "o", "", "output directory for converted images")
	convertCmd.Flags().BoolVarP(&replaceOriginal, "replace-original", "r", false, "replace original images after conversion")
}

func processFile(file, format string) error {
	// Handle image conversion here using your internal/convert.go logic
	// Check if the file is a directory
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return convertImage(path, format)
			}
			return nil
		})
	}

	// If it's not a directory, just convert the single image
	return convertImage(file, format)
}

func convertImage(inputPath, format string) error {
	// Construct the output file path if an output directory is specified
	outputPath := inputPath
	if outputDirectory != "" {
		outputPath = filepath.Join(outputDirectory, filepath.Base(inputPath))
	}

	// Change the file extension to match the new format
	newExtension := "." + format
	outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + newExtension

	// Use internal/convert.go logic to perform the image conversion
	if err := internal.ConvertFormat(inputPath, outputPath, format); err != nil {
		return err
	}

	if replaceOriginal { // Check if the flag is set
		// Remove the original image if the flag is set
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}

	return nil
}
