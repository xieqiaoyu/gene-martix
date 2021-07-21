package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type KeyGroup struct {
	genes   []string
	species []string
}

// parseSVFile 解析 Separated Value 文件 比如csv ,tsv  将文件内容解析为 按名字 - {genes,species}分组的map
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

// parseFile 根据文件名后缀尝试解析 sv 类型的文件
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
