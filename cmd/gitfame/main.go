//go:build !solution

package main

import (
	"encoding/csv"
	"encoding/json"
	flag "github.com/spf13/pflag"
	"gitlab.com/slon/shad-go/gitfame/internal"
	"gitlab.com/slon/shad-go/gitfame/internal/parse"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	flagOrderBy     = flag.StringP("order-by", "o", "lines", "ключ сортировки результатов")
	flagRepository  = flag.StringP("repository", "r", "./", "путь до Git репозитория")
	flagRevision    = flag.StringP("revision", "l", "HEAD", "указатель на коммит")
	flagUseCommiter = flag.StringP("use-committer ", "u", "lines", "булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера")
	flagFormat      = flag.StringP("format", "f", "csv", "формат вывода")
	flagExtensions  = flag.StringSlice("extensions", []string{}, "список расширений, сужающий список файлов в расчёте;")
	flagLanguages   = flag.StringSlice("languages", []string{}, "список языков (программирования, разметки и др.), сужающий список файлов в расчёте;")
	flagExclude     = flag.StringP("exclude", "e", "", "набор Glob паттернов, исключающих файлы из расчёта")
	flagRestict     = flag.StringP("restrict-to", "d", "", " набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")
)

func main() {
	flag.Parse()
	var languages []internal.Language
	file, err := os.ReadFile(`C:\Users\ADD-0\GolandProjects\shad-go\gitfame\configs\language_extensions.json`)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &languages)
	if err != nil {
		log.Fatal(err)
	}
	extensions := make(map[string]struct{})

	for _, ext := range *flagExtensions {
		extensions[ext] = struct{}{}
	}

	for _, langFlag := range *flagLanguages {
		upperName := strings.ToUpper(langFlag)
		upperNameFirst := strings.ToUpper(string(langFlag[0])) + langFlag[1:]
		for _, lang := range languages {
			if upperName == lang.Name || upperNameFirst == lang.Name {
				for _, ext := range lang.Extensions {
					extensions[ext] = struct{}{}
				}
			}
		}
	}

	statistics, _ := parse.CollectOfStatistic(*flagRevision, *flagRepository, extensions)

	names := make([]string, 0)
	for name, _ := range statistics {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		stat1 := statistics[names[i]]
		stat2 := statistics[names[j]]
		if stat1.Lines != stat2.Lines {
			return stat1.Lines > stat2.Lines
		} else {
			if stat1.Commits != stat2.Commits {
				return stat1.Commits > stat2.Commits
			} else {
				if stat1.Files != stat2.Files {
					return stat1.Lines > stat2.Lines
				}
			}
		}
		return names[i] < names[j]
	})

	output := make([][]string, len(names)+1)
	output[0] = []string{"Name", "Lines", "Commits", "Files"}
	i := 1
	for _, name := range names {
		stat := *statistics[name]
		stats := make([]string, 4)
		stats[0] = name
		stats[1] = strconv.Itoa(stat.Lines)
		stats[2] = strconv.Itoa(stat.Commits)
		stats[3] = strconv.Itoa(stat.Files)
		output[i] = stats
		i++
	}
	if *flagFormat == "csv" {
		w := csv.NewWriter(os.Stdout)
		for _, record := range output {
			if err = w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
		w.Flush()
	}
}
