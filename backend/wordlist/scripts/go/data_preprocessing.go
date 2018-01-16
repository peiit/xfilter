package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Text struct {
	Text string `json:"text"`
}

func WriteCsvFromDest(src string, dest string, tag string) {
	fi, err := os.Open(src)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	des, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer des.Close()
	w := csv.NewWriter(des)

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		text := Text{}
		json.Unmarshal(a, &text)
		var s []string
		if tag != "" {
			s = []string{tag, text.Text}
		} else {
			s = []string{text.Text}
		}
		w.Write(s)
		if err := w.Error(); err != nil {
			fmt.Printf("error writing csv:", err)
		}
	}

	w.Flush()
}

func WriteCsvFromSplit(src1 string, src2 string, dest string) {
	temp := [][]string{{}}

	f1, err := os.Open(src1)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f1.Close()
	f2, err := os.Open(src2)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f2.Close()

	br1 := bufio.NewReader(f1)
	br2 := bufio.NewReader(f2)

	des, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer des.Close()
	w := csv.NewWriter(des)
	for {
		a, _, c := br1.ReadLine()
		if c == io.EOF {
			break
		}
		s := []string{string(a)}
		temp = append(temp, s)
	}

	for {
		a, _, c := br2.ReadLine()
		if c == io.EOF {
			break
		}
		s := []string{string(a)}
		temp = append(temp, s)
	}

	for _, value := range temp {
		w.Write(value)
		if err := w.Error(); err != nil {
			fmt.Printf("error writing csv:", err)
		}
	}
	w.Flush()
}

func RawFileSplit() {
	// WriteCsvFromDest("../data/raw/sensitive.csv", "../data/split/sensitive_split.csv", "")
	WriteCsvFromDest("../../data/raw/160405/sensitive_approved.csv", "../../data/split/sensitive_approved_split.csv", "1")
	WriteCsvFromDest("../../data/raw/160405/sensitive_deleted.csv", "../../data/split/sensitive_deleted_split.csv", "0")
}

func WriteCsvForPreprocessing() {
	RawFileSplit()
	WriteCsvFromSplit("../../data/split/sensitive_approved_split.csv", "../../data/split/sensitive_deleted_split.csv", "../..//data/pre/sensitive.csv")
}

func main() {
	WriteCsvForPreprocessing()
}
