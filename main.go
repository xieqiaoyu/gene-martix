package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type KeyGroup struct {
	genes   []string
	species []string
}

// parseSVFile 解析 Separated Value 文件 比如csv ,tsv  将文件内容解析为 按名字[基因...]分组的map
// delimiter 为分割符
func parseSVFile(fileName string, delimiter string) (map[string]*KeyGroup, error) {
	svFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	scn := bufio.NewScanner(svFile)

	var lines []string

	for scn.Scan() {
		line := scn.Text()
		lines = append(lines, line)
	}

	if err := scn.Err(); err != nil {
		return nil, err
	}

	lines = lines[1:] // First line is header
	out := map[string]*KeyGroup{}

	for _, line := range lines {
		record := strings.Split(line, delimiter)
		name := record[0]
		gene := record[1]
		species := record[2]
		keyGroup, exists := out[name]
		if exists {
			keyGroup.genes = append(keyGroup.genes, gene)
			keyGroup.species = append(keyGroup.species, species)
		} else {
			keyGroup = &KeyGroup{
				genes:   []string{gene},
				species: []string{species},
			}
			out[name] = keyGroup
		}
	}

	return out, nil
}

func parseFile(fileName string) (map[string]*KeyGroup, error) {
	fileExt := filepath.Ext(fileName)
	var delimiter string
	switch fileExt {
	case ".tsv":
		delimiter = "\t"
	case ".csv":
		delimiter = ","
	default:
		return nil, fmt.Errorf("unsupport file ext %s", fileExt)
	}
	return parseSVFile(fileName, delimiter)

}

const sourceFileName = "./gene-flow.csv"

const outGeneFileName = "matrix_gene.tsv"
const outSpeciesFileName = "matrix_species.tsv"

func main() {
	parsed, err := parseFile(sourceFileName)
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
	ioutil.WriteFile(outGeneFileName, []byte(geneOutString), 0644)

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
	ioutil.WriteFile(outSpeciesFileName, []byte(speciesOutString), 0644)

	fmt.Printf("parsed %d names, %d species\n", len(nameHeader), len(speciesHeader))
}
