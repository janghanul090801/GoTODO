/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadIgnores(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var patterns []string
	for _, line := range strings.Split(string(data)+"\n*.exe", "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

func isIgnored(path string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := filepath.Match(pattern, filepath.Base(path))
		if match {
			return true
		}
		// 폴더 무시
		if strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
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

		err := filepath.Walk(pathInfo, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}

			ignores, err := loadIgnores(".gotodoignores")
			if err != nil {
				ignores = []string{}
			}

			if ext == "" {
				if !info.IsDir() && !isIgnored(path, ignores) {
					files = append(files, path)
				}
			} else {
				if !info.IsDir() && filepath.Ext(path) == ext && !isIgnored(path, ignores) {
					files = append(files, path)
				}
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
				if strings.Contains(line, "TODO") {
					fmt.Printf("%s col: %d, row: %d, detail: %s\n",
						filename,
						lineNumber,
						strings.Index(line, "TODO")+1, // 열 위치 (1부터 시작)
						line[strings.Index(line, "TODO"):len(line)-1],
					)
				}
				lineNumber++
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(whereismyfuckingtodoCmd)
	whereismyfuckingtodoCmd.Flags().String("path", "", "파일 검색 경로를 지정해주세요(default=./)")
	whereismyfuckingtodoCmd.Flags().String("ext", "", "검색할 파일 확장자를 지정해주세요(default=.go)")
}
