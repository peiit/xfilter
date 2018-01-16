package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/huichen/sego"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

var Segmenter sego.Segmenter
var StopWords map[string]bool
var CorpusSum CorpusCount
var Terms map[string]Term
var TermsList []string
var CorpusX map[string]float64
var CorpusY map[string]float64
var TermsSum TermsCount

const (
	Weight = 2.0 // α
	// for KNN
	Threshold = 0.7 // ε
	K         = 15  // k
)

type Term struct {
	Label int32
	A     float64 // contains `x` c1 sum num
	B     float64 // contains `x` c0 sum num
	C     float64 // c1 not contains `x` sum num
	D     float64 // c0 not contains `x` sum num
	Nx    float64 // `x` sum num in c1
	Ny    float64 // `y` sum num in c0
	TFIDF float64 // A / X * log(N / A) + (B / Y * log(N / B) * α
	CHI   float64 // CHI
}

type Pair struct {
	Index int32
	Label int32
	CHI   float64
	TFIDF float64
}

type CorpusCount struct {
	X float64 // c1
	Y float64 // c0
}

type TermsCount struct {
	X float64 // c1 sum term
	Y float64 // c0 sum term
}

type PairList map[string]Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].CHI < p[j].CHI }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	initVariable()
	loadStopWords()
	loadDict()
	loadData()
	pl := caculateEigenvalue()
	// exportCsv(pl, "../data/pre/knn_training.csv", "KNN")
	exportCsv(pl, "../data/pre/svm_feature.csv", "feature")
	exportSVMCsv(pl, "../data/pre/svm_training.csv")
	fmt.Println("done!")
}

func initVariable() {
	StopWords = make(map[string]bool)
	Terms = make(map[string]Term)
	CorpusX = make(map[string]float64)
	CorpusY = make(map[string]float64)
}

func isStopWord(s string) bool {
	if StopWords[s] {
		return true
	} else {
		return false
	}
}

func updateCorpusCount(c string, perLineTexts map[string]bool) {
	// update `text` in corpus `X` or `Y`count
	var corpusCount map[string]float64
	switch c {
	case "1":
		corpusCount = CorpusX
		break
	case "0":
		corpusCount = CorpusY
		break
	}
	for text, _ := range perLineTexts {
		count, ok := corpusCount[text]
		if ok {
			count += 1
		} else {
			count = 1
		}
		corpusCount[text] = count
	}
}

func setOrUpdateTerms(c string, text string) {
	// check if Nx ++ OR Ny ++
	term, _ := Terms[text]
	switch c {
	case "1":
		term.Nx += 1
		term.Label = 1
		TermsSum.X += 1
		break
	case "0":
		term.Ny += 1
		term.Label = 0
		TermsSum.Y += 1
		break
	}
	Terms[text] = term
}

func chineseSegment(c string, s string) {
	segments := Segmenter.Segment([]byte(s))
	var perLineTexts map[string]bool
	perLineTexts = make(map[string]bool)
	if len(segments) > 0 {
		for _, seg := range segments {
			token := seg.Token()
			var text = token.Text()
			if !isStopWord(text) && len(strings.TrimSpace(text)) != 0 {
				TermsList = append(TermsList, text) // put in termslist
				if !perLineTexts[text] {
					perLineTexts[text] = true
				}
				setOrUpdateTerms(c, text)
			}
		}
		// update corpus count
		updateCorpusCount(c, perLineTexts)
	}
}

func caculateEigenvalue() PairList {
	// caculate CHI and TFIDF
	for key, term := range Terms {
		term.A = CorpusX[key]
		term.B = CorpusY[key]
		term.C = CorpusSum.X - term.A
		term.D = CorpusSum.Y - term.B
		cn := CorpusSum.X + CorpusSum.Y
		a := term.A
		b := term.B
		c := term.C
		d := term.D
		nx := term.Nx
		ny := term.Ny

		if a == 0 && b != 0 {
			term.TFIDF = nx/TermsSum.X*1.0 + ny/TermsSum.Y*(1.0+math.Log10(cn/b))*Weight
		} else if a != 0 && b == 0 {
			term.TFIDF = nx/TermsSum.X*(1.0+math.Log10(cn/a)) + ny/TermsSum.Y*1.0*Weight
		} else if a == 0 && b == 0 {
			term.TFIDF = nx/TermsSum.X*1.0 + ny/TermsSum.Y*1.0*Weight
		} else {
			term.TFIDF = nx/TermsSum.X*(1.0+math.Log10(cn/a)) + ny/TermsSum.Y*(1.0+math.Log10(cn/b))*Weight
		}
		term.CHI = float64(cn) * math.Pow(float64(a*d-b*c), 2) / float64((a+c)*(b+d)*(a+b)*(c+d))
		Terms[key] = term
	}
	// sort
	pl := make(PairList, len(Terms))
	i := 0
	for key, term := range Terms {
		i++
		pl[key] = Pair{int32(i), term.Label, term.CHI, term.TFIDF}
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func exportSVMCsv(pl PairList, dest string) {
	des, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer des.Close()
	w := csv.NewWriter(des)
	for _, t := range TermsList {
		var s []string
		s = string{fmt.Sprintf("%d", p.Label), fmt.Sprintf("%d:%.6f", p.Index, p.TFIDF)}
		w.Write(s)
		if err := w.Error(); err != nil {
			fmt.Printf("error writing csv:", err)
		}
	}
	w.Flush()
}

func exportCsv(pl PairList, dest string, class string) {
	des, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer des.Close()
	w := csv.NewWriter(des)
	for key, p := range pl {
		var s []string
		switch class {
		case "KNN":
			if p.CHI > Threshold {
				s = []string{fmt.Sprintf("%d", p.Label), fmt.Sprintf("%.6f", p.CHI), fmt.Sprintf("%.6f", p.TFIDF)}
			}
			break
		case "SVM":
			s = []string{fmt.Sprintf("%d", p.Label), fmt.Sprintf("%d:%.6f", index+1, p.TFIDF)}
			break
		case "feature":
			s = []string{fmt.Sprintf("%d", index+1), p.Key}
			break
		}
		w.Write(s)
		if err := w.Error(); err != nil {
			fmt.Printf("error writing csv:", err)
		}
	}
	w.Flush()
}

func loadDict() {
	dictSrc := "../data/dictionary.txt,../data/sensitive.txt"
	Segmenter.LoadDictionary(dictSrc)
}

func loadStopWords() {
	fi, err := os.Open("../data/stop_words.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		StopWords[string(a)] = true
	}
}

func loadData() {
	fi, err := os.Open("../data/pre/sensitive.csv")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	reader := csv.NewReader(fi)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if record[0] == "0" {
			CorpusSum.X += 1
		} else {
			CorpusSum.Y += 1
		}
		chineseSegment(record[0], record[1])
	}
}
