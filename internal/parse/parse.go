package parse

import (
	"gitlab.com/slon/shad-go/gitfame/internal"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/capture"
)

type FileEmptiness int

const (
	FileEmpty FileEmptiness = iota
	FileNotEmpty
)

func criterionCheck(fileName string, extensions map[string]struct{}) bool {
	if len(extensions) == 0 {
		return true
	}
	partsOfName := strings.Split(fileName, ".")
	ext := partsOfName[len(partsOfName)-1]
	_, ok := extensions["."+ext]
	return ok
}

func parsePorcelainNotEmpty(statistics map[string]*internal.Statistic, line string, authorOfCommit map[string]string) {
	lines := strings.Split(line, "\n")
	setOfAuthors := make(map[string]struct{})
	for ind, lin := range lines {
		components := strings.Split(lin, " ")
		if components[0] == "author" {
			author := strings.Join(components[1:], " ")
			componentsPrev := strings.Split(lines[ind-1], " ")
			_, authorExist := setOfAuthors[author]
			if !authorExist {
				setOfAuthors[author] = struct{}{}
			}
			_, commitExist := authorOfCommit[componentsPrev[0]]
			if !commitExist {
				authorOfCommit[componentsPrev[0]] = author
			}
			linesNumb, _ := strconv.Atoi(componentsPrev[3])
			addAuthorInStatistic(statistics, author, linesNumb, authorExist, commitExist)
		}
		authorName, ok := authorOfCommit[components[0]]
		if ok && len(components) == 4 {
			stat, _ := statistics[authorName]
			linesNumb, _ := strconv.Atoi(components[3])
			(*stat).Lines += linesNumb
		}
	}
}

func addAuthorInStatistic(statistics map[string]*internal.Statistic, author string, lineNumber int, authorExist, commitExist bool) {
	_, ok := statistics[author]
	if !ok {
		statistics[author] = &internal.Statistic{
			Commits: 0,
			Lines:   0,
			Files:   0,
		}
	}
	stat, _ := statistics[author]
	if !authorExist {
		(*stat).Files++
	}

	if !commitExist {
		(*stat).Commits++
	}

	(*stat).Lines += lineNumber
}

func IterateBySystemTree(revision, files, names, path string, statistics map[string]*internal.Statistic, extensions map[string]struct{}, authorOfCommit map[string]string) {
	filesSplit := strings.Split(files, "\n")
	namesSplit := strings.Split(names, "\n")
	for ind, file := range filesSplit {
		fileInfo := strings.Split(file, " ")
		if len(fileInfo) > 1 {
			path1 := filepath.Join(path, namesSplit[ind]) + `\`
			if fileInfo[1] == "blob" && criterionCheck(namesSplit[ind], extensions) {
				line, _ := capture.MakeListOfLastCommitsOfFile(revision, namesSplit[ind], path)
				if len(line) == 0 {
					nameWithHash, _ := capture.MakeListOfEmptyFileChangers(revision, namesSplit[ind], path)
					nameWithHash = nameWithHash[1 : len(nameWithHash)-1]
					parts := strings.Split(nameWithHash, " ")
					hash := parts[len(parts)-1]
					name := strings.Join(parts[:len(parts)-1], " ")
					_, commitExist := authorOfCommit[hash]
					if !commitExist {
						authorOfCommit[hash] = name
					}
					addAuthorInStatistic(statistics, name, 0, false, commitExist)
				} else {
					parsePorcelainNotEmpty(statistics, line, authorOfCommit)
				}
			} else if fileInfo[1] == "tree" {
				names1, _ := capture.MakeListOfFileNames(revision, path1)
				directories, _ := capture.MakeListOfFilesAndDirectories(revision, path1)
				IterateBySystemTree(revision, directories, names1, path1, statistics, extensions, authorOfCommit)
			}
		}
	}
}

func CollectOfStatistic(revision, path string, extensions map[string]struct{}) (map[string]*internal.Statistic, error) {
	statistics := make(map[string]*internal.Statistic)
	authorOfCommit := make(map[string]string)

	directories, _ := capture.MakeListOfFilesAndDirectories(revision, path)
	filenames, _ := capture.MakeListOfFileNames(revision, path)

	IterateBySystemTree(revision, directories, filenames, path, statistics, extensions, authorOfCommit)

	return statistics, nil
}
