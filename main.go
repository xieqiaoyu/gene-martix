package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/xieqiaoyu/gene-martix/metadata"
)

var usageStr = `
Usage:  [options] SVFILE
Options:
     -h --help                  show help
     -v --version               show version
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// Config 程序的配置参数
type Config struct {
	SourceFileName        string
	OutPutGeneFileName    string
	OutPutSpeciesFileName string
	ShowVersion           bool
	ShowHelp              bool
}

func getConfig() (config *Config) {
	config = &Config{
		SourceFileName:        "",
		OutPutGeneFileName:    "matrix_gene.tsv",
		OutPutSpeciesFileName: "matrix_species.tsv",
	}
	fs := flag.NewFlagSet("f", flag.ExitOnError)
	fs.Usage = usage

	fs.BoolVar(&config.ShowVersion, "v", false, "")
	fs.BoolVar(&config.ShowVersion, "version", false, "")
	fs.BoolVar(&config.ShowHelp, "h", false, "")
	fs.BoolVar(&config.ShowHelp, "help", false, "")
	fs.Parse(os.Args[1:])
	fileNames := fs.Args()
	if len(fileNames) > 0 {
		config.SourceFileName = fileNames[0]
	}
	return
}

func main() {
	config := getConfig()

	if config.ShowVersion {
		metadata.ShowVersion()
		return
	}
	if config.ShowHelp || config.SourceFileName == "" {
		usage()
		return
	}

	parsed, err := parseFile(config.SourceFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	nextNameIndex := 0
	nextGeneIndex := 0
	nextSpeciesIndex := 0

	nameIndexMap := map[string]int{}
	geneIndexMap := map[string]int{}
	speciesIndexMap := map[string]int{}

	nameCount := len(parsed)
	// 基因[名称...]的二维矩阵 (为了动态扩展,因为名称数量已经确定了）
	geneMatrix := [][]int{}
	speciesMatrix := [][]int{}

	genesHeader := []string{}
	nameHeader := []string{}
	speciesHeader := []string{}

	for name, keyGroup := range parsed {
		nameIndex, exists := nameIndexMap[name]
		if !exists {
			nameIndex = nextNameIndex
			nameIndexMap[name] = nameIndex
			nextNameIndex++
			nameHeader = append(nameHeader, name)
		}

		genes := keyGroup.genes
		for _, gene := range genes {
			geneIndex, exists := geneIndexMap[gene]
			if !exists {
				geneIndex = nextGeneIndex
				nextGeneIndex++
				genesHeader = append(genesHeader, gene)
				geneIndexMap[gene] = geneIndex
				// 动态添加新的数组到矩阵中
				geneMatrix = append(geneMatrix, make([]int, nameCount))
			}
			geneMatrix[geneIndex][nameIndex] = 1
		}

		manySpecies := keyGroup.species

		for _, species := range manySpecies {
			speciesIndex, exists := speciesIndexMap[species]
			if !exists {
				speciesIndex = nextSpeciesIndex
				nextSpeciesIndex++
				speciesHeader = append(speciesHeader, species)
				speciesIndexMap[species] = speciesIndex
				// 动态添加新的数组到矩阵中
				speciesMatrix = append(speciesMatrix, make([]int, nameCount))
			}
			speciesMatrix[speciesIndex][nameIndex]++
		}
	}

	// 对header 排序保证输入顺序不变
	sort.Strings(nameHeader)
	sort.Strings(genesHeader)
	sort.Strings(speciesHeader)

	geneOutString := "key gene\t" + strings.Join(genesHeader, "\t")
	for _, name := range nameHeader {
		geneOutString += fmt.Sprintf("\n%s", name)
		nameIndex := nameIndexMap[name]
		for _, gene := range genesHeader {
			geneIndex := geneIndexMap[gene]
			flag := geneMatrix[geneIndex][nameIndex]
			geneOutString += fmt.Sprintf("\t%d", flag)
		}
	}
	// 最后加一个换行符
	geneOutString += "\n"
	ioutil.WriteFile(config.OutPutGeneFileName, []byte(geneOutString), 0644)

	fmt.Printf("parsed %d names, %d genes\n", len(nameHeader), len(genesHeader))

	speciesOutString := "key species\t" + strings.Join(speciesHeader, "\t")
	for _, name := range nameHeader {
		speciesOutString += fmt.Sprintf("\n%s", name)
		nameIndex := nameIndexMap[name]
		for _, species := range speciesHeader {
			speciesIndex := speciesIndexMap[species]
			flag := speciesMatrix[speciesIndex][nameIndex]
			speciesOutString += fmt.Sprintf("\t%d", flag)
		}
	}
	// 最后加一个换行符
	speciesOutString += "\n"
	ioutil.WriteFile(config.OutPutSpeciesFileName, []byte(speciesOutString), 0644)

	fmt.Printf("parsed %d names, %d species\n", len(nameHeader), len(speciesHeader))
}
