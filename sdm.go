package main

import (
	"io/ioutil"
	"strings"
	"net/http"
	"os"
	"fmt"
)

func main() {
	// Load tags (articles) list from file
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	input :=strings.Split(string(b), "\r\n")

	//Some statistic
	fmt.Println("Statistics:")
	fmt.Print("Input tags count: ")
	fmt.Println(len(input))
	fmt.Println("Parsing starts:")

	// Parse wikipedia by tags list
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, i := range input{
		fmt.Print("Parse " + i)
		result:= getWikiHeadlines(i)
		//Append output file with results
		if _, err = f.WriteString(strings.Join(result,"\r\n")); err != nil {
			panic(err)
		}
		fmt.Println(" - OK")
	}

	fmt.Println("Parsing finished!")
	fmt.Scanln()
}
// Func gets headlines from raw article text by finding specific meta-tag
func getWikiHeadlines(word string) []string{
	var result = make([]string,0)
	resp, _ := http.Get("http://en.wikipedia.org/wiki/"+word) //Get article by tag
	body, _ := ioutil.ReadAll(resp.Body)
	pageString :=string(body) //Source of wiki article web page
	result = append(result,"\r\n--------------------")
	result = append(result,strings.ToTitle(word)+":")
	//Find headlines (meta tags) and parse it if exist
	i:= strings.Index(pageString,"mw-headline")
	for  i!= -1 {
		result = append(result,parser(pageString,i))
		pageString = pageString[i+2:]
		i= strings.Index(pageString,"mw-headline")
	}
	return result
}

func parser(s string, i int)string {
	var result string
	i = i+17; //Offset from meta teg beginning
	for string(s[i]) != "\""{
		result +=string(s[i])
		i++;
	}
	return strings.ToLower(result)
}