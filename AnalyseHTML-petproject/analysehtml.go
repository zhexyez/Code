package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var InScript bool = false

type attr struct {
	name string
}

type elem struct {
	name string
	attrs []attr
}

func parseattr(el string) (out []attr, name string) {
	if el[0] == '/' || el[0] == '!' {
		if el == "/script" {
			InScript = false
		}
	  return nil, ""
	}
	namebuf := ""
	inattr := false
	firstattr := false
	tempbuf := ""
	tempslice := []string{}
	for iter, ch := range el {
		if ch == '"' && !inattr {
			inattr = true
		} else if ch =='"' && inattr {
			inattr = false
		}
		if ch == ' ' && !inattr {
			if firstattr {
				tempslice = append(tempslice, tempbuf)
				tempbuf = ""
				continue
			}
			namebuf = tempbuf
			firstattr = true
			if tempbuf == "script" {
				InScript = true
			}
			tempbuf = ""
			continue
		}
		if firstattr && iter == len(el) {
			tempslice = append(tempslice, tempbuf)
			break
		}
		tempbuf += string(ch)
	}
	if len(tempslice) == 0 {
		return nil, tempbuf
	}
	for _, val := range tempslice {
		out = append(out, attr{name:val})
	}
	return out, namebuf
}

func parse(html string) {
	elems := []elem{}
	elemopened := false
	tempbuf := ""
	for iter, ch := range html {
		if ch == '<' {
			if InScript {
				_,_ = parseattr(string(html[iter+1:iter+8]))
			}
			elemopened = true
			continue
		}
		if ch == '>' {
			if InScript { continue }
			if tempbuf == "" { continue }
			pattrs, elemname := parseattr(tempbuf)
			if pattrs != nil && elemname != "" {
				elems = append (elems, elem{name: elemname, attrs: pattrs})
			}
			elemopened = false
			tempbuf = ""
			continue
		}
		if elemopened && !InScript {
			tempbuf += string(ch)
			continue
		}
	}
	for _, e := range elems {
		fmt.Println(e.name)
		for _, i := range e.attrs {
			fmt.Println("|>", i.name)
		}
		fmt.Print("\n")
	}
}

func getThem(URL string) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	parse(string(body))
}

func main() {
	if len(os.Args) != 2 {
		log.Println("This program accepts only 1 argument of type URL")
		return
	}
	getThem(os.Args[1])
}