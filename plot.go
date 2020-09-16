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

type data struct {
	x []string
	y []string
}

type Message struct{
	Delimiter string
	Path string  
	PlotType string
	Title string 
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
	line.AddXAxis(d.x)
	line.AddYAxis("Y",d.y)
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
	scatter.AddXAxis(d.x)
	scatter.AddYAxis("Y",d.y)
	f, err := os.Create("tmp/plot.html")
	if err != nil{
		return err
	}
	scatter.Render(f)
	return nil
}

func plot(chartType int, d data, name string) error {	
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
	json.NewDecoder(r.Body).Decode(&msg)
	fmt.Printf("%+v\n",msg)
	s := strings.Split(msg.Path,"\\")
	path := s[len(s)-1]
	plotType, _ := strconv.Atoi(msg.PlotType)
	title := msg.Title
	delimiter:= msg.Delimiter
	lines, _ := openCsv(path)
	d := getData(lines, delimiter)
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
