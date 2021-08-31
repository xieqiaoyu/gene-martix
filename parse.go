package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type KeyGroup struct {
	genes   []string
	species []string
}

// parseFile 根据文件名后缀尝试解析 sv 类型的文件
func parseFile(fileName string) (map[string]*KeyGroup, error) {
	fileExt := filepath.Ext(fileName)
	var delimiter rune
	switch fileExt {
	case ".tsv":
		delimiter = '\t'
	case ".csv":
		delimiter = ','
	default:
		return nil, fmt.Errorf("unsupport file ext %s", fileExt)
	}
	// 经过效率测试手写的解析速度更快
	return parseSVFileDirty(fileName, delimiter)

}

// parseSVFile 解析 Separated Value 文件 比如csv ,tsv  将文件内容解析为 按名字 - {genes,species}分组的map
// delimiter 为分割符
// 使用csv库实现解析的版本
func parseSVFile(fileName string, delimiter rune) (map[string]*KeyGroup, error) {
	svFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(svFile)
	csvReader.Comma = delimiter
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	records = records[1:] // First line is header
	out := map[string]*KeyGroup{}
	for _, record := range records {
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

// parseSVFileDirty 解析 Separated Value 文件 比如csv ,tsv  将文件内容解析为 按名字 - {genes,species}分组的map
// delimiter 为分割符
// 不使用csv库实现解析的版本
func parseSVFileDirty(fileName string, delimiter rune) (map[string]*KeyGroup, error) {
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
		record := strings.Split(line, string(delimiter))
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
