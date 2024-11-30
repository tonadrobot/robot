package bot

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
	"gopkg.in/macaron.v1"
)

func generateCode() string {
	code := ""
	num := strconv.Itoa(generateRandNum(99))

	words, err := urlToLines(WordsUrl)
	if err != nil {
		loge(err)
	}

	code = words[generateRandNum(2047)]

	return code + num
}

func generateRandNum(max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	min := 0
	rn := rand.Intn(max-min+1) + min
	return rn
}

func urlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return linesFromReader(resp.Body)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getTgId(ctx *macaron.Context) int64 {
	tgids := ctx.Params("telegramid")
	tgid, err := strconv.Atoi(tgids)
	if err != nil {
		loge(err)
	}
	return int64(tgid)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
