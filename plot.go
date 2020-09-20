package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/charts"
)

type axe struct{
	values []string
	header string
}

type data struct {
	xs []axe
	ys []axe
}

type Message struct{
	Delimiter string
	Path string  
	PlotType string
	Title string
	HaveTitles string
	Format string
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

func getFormat(format string, delimiter string) []int{
	var parsed []int
	s := strings.Split(format, delimiter)
	for i := range(s){
		if s[i] == "x"{
			parsed = append(parsed,0)
			fmt.Printf("Appending x...\n")
		}else if(s[i] == "y"){
			parsed = append(parsed, 1)
			fmt.Printf("Appending y...\n")
		}
	}
	return parsed
}
		

func getData(lines []string, delimiter string, format string, headers bool) data {
	var d data
	var j, k, header_flag int
	parsed := getFormat(format, delimiter)
	fmt.Printf("Parsed format is %v\n", parsed)
	j = 0
	k = 0
	header_flag = 0
	if headers{
		hds := strings.Trim(lines[0], "\n")
		s := strings.Split(hds, delimiter)
		for i := range(s){
			var a axe
			if parsed[i] == 0 {
				d.xs = append(d.xs, a)
				d.xs[j].header = s[i]
				j = j+1
			}else{
				d.ys = append(d.ys, a)
				d.ys[k].header = s[i]
				k = k+1
			}
		}
		header_flag = 1
	}else{
		for i := range(parsed){
			var a axe
			if parsed[i] == 0{
				d.xs = append(d.xs, a)
				d.xs[j].header = " "
				j = j+1
			}else{
				d.ys = append(d.ys, a)
				d.ys[k].header = " "
				k = k+1
			}
		}
	}
	j = 0
	k = 0
	for i := header_flag; i < len(lines); i++ {
		lines[i] = strings.Trim(lines[i], "\n")
		s := strings.Split(lines[i], delimiter)
		for l := range(s){
			if parsed[l] == 0 {
				fmt.Printf("Addding x value\n")
				d.xs[j].values = append(d.xs[j].values, s[l])
				j = j + 1
			}else{
				fmt.Printf("Addding y value\n")
				d.ys[k].values = append(d.ys[k].values, s[l])
				k = k +1
			}
		}
		j = 0
		k = 0
	}
	fmt.Printf("Data ready!\n")
	return d
}


func renderBar(d data, name string) error{

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: name})
	for i := range(d.xs){
		for j := range(d.ys){
			bar.AddXAxis(d.xs[i].values).AddYAxis(d.ys[j].header,d.ys[j].values)
		}
	}
	f, err := os.Create("tmp/plot.html")
	if err != nil{
		return err
	}
	bar.Render(f)
	return nil
}

func renderLine(d data, name string) error{

	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: name})
	for i := range(d.xs){
		line.AddXAxis(d.xs[i].values)
	}
	for j := range(d.ys){
		line.AddYAxis(d.ys[j].header, d.ys[j].values)
	}
	f, err := os.Create("tmp/plot.html")
	if err != nil{
		return err
	}
	line.Render(f)
	return nil
}

func renderScatter(d data, name string) error{

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(charts.TitleOpts{Title: name})
	for i := range(d.xs){
		scatter.AddXAxis(d.xs[i].values)
	}
	for j := range(d.ys){
		scatter.AddYAxis(d.ys[j].header, d.ys[j].values)
	}
	f, err := os.Create("tmp/plot.html")
	if err != nil{
		return err
	}
	scatter.Render(f)
	return nil
}

func plot(chartType int, d data, name string) error {
	fmt.Printf("Plotting...\n")
	switch chartType {
	case 0:
		err := renderBar(d, name)
		return err
	case 1:
		err := renderLine(d, name)
		return err
	case 2:
		err := renderScatter(d, name)
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

func renderTemplate(w http.ResponseWriter, tmpl string) {                                  
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)                                        
	if err != nil{                                                                              
		http.Error(w, err.Error(), http.StatusInternalServerError)                          
		return                                                                              
	}                                                                                           
}    

func plotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Plotting...\n")
	var msg Message
	var haveHeader bool
	json.NewDecoder(r.Body).Decode(&msg)
	fmt.Printf("%+v\n",msg)
	s := strings.Split(msg.Path,"\\")
	path := s[len(s)-1]
	plotType, _ := strconv.Atoi(msg.PlotType)
	title := msg.Title
	delimiter:= msg.Delimiter
	if msg.HaveTitles == "True"{
		haveHeader = true
		fmt.Printf("I have headers!\n")
	}else{
		haveHeader = false
		fmt.Printf("I do not have headers!\n")
	}
	format := msg.Format
	lines, _ := openCsv(path)
	d := getData(lines, delimiter, format, haveHeader)
	plot(plotType, d, title)
	http.Redirect(w, r, "/show/", http.StatusSeeOther)                                       
}             

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view")
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmp/plot.html")
	if err != nil{
		fmt.Printf("Error loading plot\n")
	}else{
		fmt.Printf("Showing plot...\n")
	}
	t.Execute(w, nil)
}

var port int
var templates = template.Must(template.ParseFiles("templates/view.html"))
func init() {
	port, _ = readConfig()

}

func main() {
	fmt.Printf("Listening on port %d...\n", port)
	http.HandleFunc("/view/", viewHandler)                
	http.HandleFunc("/plot/", plotHandler)
	http.HandleFunc("/show/", showHandler)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}
