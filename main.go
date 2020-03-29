package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

const URL = "https://jsonplaceholder.typicode.com/comments"
const templateName = "template/comments.html"
const outputPdfPath = "pdf/comments.pdf"
const tempFolderPath = "temp/"

type Comments []struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func parseTemplate(templateName string, data Comments) (buf *bytes.Buffer, err error){
	t, err := template.ParseFiles(templateName)
	if err != nil {
		log.Fatal(err)
	}
	buf = new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Fatal(err)
	}
	return buf, err
}

func generatePDF(templateContent string) {
	t := time.Now().Unix()
	fileName := tempFolderPath + strconv.FormatInt(int64(t), 10) + ".html"
	err1 := ioutil.WriteFile(fileName, []byte(templateContent), 0644)
	if err1 != nil {
		panic(err1)
	}
	f, err := os.Open(fileName)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA3)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(outputPdfPath)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(c Comments) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	dprint := func(str string) string {
		dashes := strings.Repeat(str, 4)
		return dashes

	}
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s", "POSTID", "ID", "NAME", "EMAIL")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s", dprint("-"), dprint("-"), dprint("-"), dprint("-"))
	for _, elem := range c {
		fmt.Fprintf(w, "\n %d\t%d\t%s\t%s", elem.PostID, elem.ID, elem.Name, elem.Email)
	}
}

func main() {

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	comments := Comments{}
	jsonErr := json.Unmarshal(body, &comments)
	if jsonErr != nil {
		log.Fatal(jsonErr)

	}
	prettyPrint(comments)
	buf, err := parseTemplate(templateName, comments)
	if err != nil {
		log.Fatal(err)
	}else{
		generatePDF(buf.String())
	}
}
