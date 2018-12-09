package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fnames := checkExt(".exp")
	allGenes := make(map[string]bool)
	data := make(map[string]map[string]string)
	for i, fname := range fnames {
		fmt.Printf("Processing %v/%v: %v\n", i+1, len(fnames), fname)
		fnameTokens := strings.Split(fname, ".")
		data[fnameTokens[0]] = make(map[string]string)
		file, err := os.Open("Expression/" + fname)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		genes := []string{}
		tmps := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, "\t")
			if !strings.Contains(tokens[0], "ENSG") {
				continue
			}
			genes = append(genes, tokens[0])
			tmps = append(tmps, tokens[8])
			data[fnameTokens[0]][tokens[0]] = tokens[8]
		}
		out, err := os.Create("PreparedData/" + fnameTokens[0] + ".csv")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Total %v genes\n", len(genes)-1)
		fmt.Fprintln(out, strings.Join(genes[1:], ","))
		fmt.Fprintln(out, strings.Join(tmps[1:], ","))
		for _, gene := range genes[1:] {
			allGenes[gene] = true
		}
		out.Close()
		file.Close()
	}
	out, err := os.Create("expression.csv")
	defer out.Close()
	if err != nil {
		panic(err)
	}
	allGenesSlice := []string{}
	allGenesSlice = append(allGenesSlice, "Patient")
	for gene := range allGenes {
		allGenesSlice = append(allGenesSlice, gene)
	}
	fmt.Fprintln(out, strings.Join(allGenesSlice, ","))
	for patient := range data {
		fmt.Println("Create data for patient: " + patient)
		line := []string{patient}
		for _, gene := range allGenesSlice[1:] {
			if tmp, ok := data[patient][gene]; ok {
				line = append(line, tmp)
			} else {
				line = append(line, "0")
			}
		}
		fmt.Fprintln(out, strings.Join(line, ","))
	}
}

func checkExt(ext string) []string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(path.Join(wd, "Expression"), func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}
