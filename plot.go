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


func main() {
	lines, err := openCsv("test.csv")
	if err != nil{
		fmt.Printf("Could not read CSV")
	}
	d:= getData(lines, ",")
	plot(0, d, "Test")
}
