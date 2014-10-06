package main

import (
	l4g "code.google.com/p/log4go"
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"strings"
)

const LOG_CONFIGURATION_FILE = "logging-conf.xml"

func init() {
	l4g.LoadConfiguration(LOG_CONFIGURATION_FILE)
}

func main() {
	l4g.Info("Process ID: %d", os.Getpid())
	usage := `geiger.

Usage:
  geiger list --path=<path_location>
  geiger -h | --help
  geiger -v | --version

Options:
  -h --help     Show this screen.
  -v --version  Show version.
  --rfc=<rfc>   RFC filter to reduce the search space [default: *].
  --date=<date> Date filter to reduce the search space [default: today]. Format: YYYY-MM-DD (zero-padding!)`

	arguments, _ := docopt.Parse(usage, nil, true, "geiger 0.0.0", false)
	l4g.Debug(arguments)
	c := make(<-chan int)

	if arguments["list"].(bool) {
		fmt.Println(strings.Join([]string{"NoFactura","FechaGeneracion","razonSocial","claveprod", "cant","unidad","descripcion","Valorunit","importe"}, "\t"))
		c = WriteCount(Gen(arguments))
	}
	l4g.Info("geiger stopped")
	<-c
}

func Gen(options map[string]interface{}) <-chan string {
	out := make(chan string)
	go func() {
		globPatternList := GetGlobPatternList(options)
		l4g.Info("Directorios encontrados: %d", len(globPatternList))

		for _, globPatternTuple := range globPatternList {
			files, _ := ListFiles(globPatternTuple)
			l4g.Info("%d archivos en directorio %s", len(files), globPatternTuple)
			for _, filePath := range files {
				out <- filePath
				for _, row := range EncodeAsRows(filePath) {
					fmt.Println(row)
				}
			}
		}
		close(out)
	}()
	return out
}

func WriteCount(in <-chan string) <-chan int {
	out := make(chan int)
	go func() {
		for record := range in {
			l4g.Debug(record)
		}
		out <- 1
		close(out)
	}()
	return out
}
