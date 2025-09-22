/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ig *ignore.GitIgnore

func loadIgnores(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultIgnores := []string{
			".git/",
			".idea/",
			"*.exe",
			"gotodo",
		}

		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		for _, line := range defaultIgnores {
			_, _ = file.WriteString(line + "\n")
		}
	}

	var err error
	ig, err = ignore.CompileIgnoreFile(filename)
	if err != nil {
		panic(err)
	}
}

func isIgnored(path string) bool {
	if ig == nil {
		return false
	}
	return ig.MatchesPath(path)
}

// whereismyfuckingtodoCmd represents the whereismyfuckingtodo command
var whereismyfuckingtodoCmd = &cobra.Command{
	Use:   "whereismyfuckingtodo",
	Short: "내가 싸질러 놓은 TODO 들을 찾아줍니다",
	Long:  `내가 여기저기 싸질러 놓은 TODO 들을 찾아줍니다. 고맙죠?`,
	Run: func(cmd *cobra.Command, args []string) {
		var files []string

		var pathInfo string

		term, _ := cmd.Flags().GetString("path")
		if term != "" {
			pathInfo = term
		} else {
			pathInfo = "./"
		}

		ext, _ := cmd.Flags().GetString("ext")

		loadIgnores(".gotodoignores")

		err := filepath.Walk(pathInfo, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}

			ignoreBase := ".gotodoignores"
			ignoreDir := filepath.Dir(ignoreBase)
			relPath, err := filepath.Rel(ignoreDir, path)

			if isIgnored(relPath) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if !info.IsDir() && (ext == "" || filepath.Ext(path) == ext) {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				idx := strings.Index(line, "TODO")
				if idx != -1 {
					fmt.Printf("%s col: %d, row: %d, detail: %s\n",
						filename,
						idx+1,
						lineNumber,
						line[idx:],
					)
				}
				lineNumber++
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			file.Close()
		}
	},
}

func init() {
	rootCmd.AddCommand(whereismyfuckingtodoCmd)
	whereismyfuckingtodoCmd.Flags().String("path", "", "파일 검색 경로를 지정해주세요(default=./)")
	whereismyfuckingtodoCmd.Flags().String("ext", "", "검색할 파일 확장자를 지정해주세요(default=.go)")
}
