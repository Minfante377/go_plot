package main

import (
	"bufio"
	"strconv"
	"fmt"
	"os"
	"strings"
	"github.com/go-echarts/go-echarts/charts"
)

type data struct {
	x []string
	y []string
}

func openCsv(path string) ([]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully opened CSV file")
	scanner := bufio.NewScanner(csvFile)
	var csvLines []string
	for scanner.Scan(){
		csvLines = append(csvLines, scanner.Text())
	}
	return csvLines,nil
}

func getData(lines []string, delimiter string) data {
	var d data
	for i, _ := range(lines){
		lines[i] = strings.Trim(lines[i], "\n")
		s := strings.Split(lines[i], delimiter)
		d.x = append(d.x, s[0])
		d.y = append(d.y, s[1])
	}
	return d
}


func renderBar(d data, name string) error{

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: name})
	bar.AddXAxis(d.x).AddYAxis("",d.y)
	f, err := os.Create("tmp/"+name+".html")
	if err != nil{
		return err
	}
	bar.Render(f)
	return nil
}

func plot(chartType int, d data, name string) error {	
	switch chartType {
	case 0:
		err := renderBar(d, name)
		return err
	}
	return nil
}

func readConfig() (int, error) {
	configFile, err := os.Open("service.config")
	if err != nil {
		return 8080, err
	}
	fmt.Println("Successfully opened Configuration file")
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan(){
		line := scanner.Text()
		if strings.Contains(line, "port") { 
			p := strings.Split(line, "=")
			p[1] = strings.TrimSpace(p[1])
			port, _ := strconv.Atoi(p[1])
			return port, nil
		}else{
			port := 8080
			fmt.Printf("Could not find port in config file. Listening on default")
			return port, nil
		}
	}
	return port, nil
}

var port int

func init() {
	port, _ = readConfig()
}

func main() {
	fmt.Printf("Listening on port %d...\n", port)

}
