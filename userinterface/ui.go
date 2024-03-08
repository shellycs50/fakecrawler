package userinterface

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func getAnswer(prompt string) (string, error) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	sentence, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	} else {
		return sentence[:len(sentence)-2], nil
	}
}

func urlIsLive(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

func isCompatibleNumber(num string) (int, error) {
	result, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}
	if result < 1 || result > 100 {
		return 0, errors.New("number not in valid range")
	}
	return result, nil
}

func isValidFileName(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return match
}

func GetCrawlArgs() (string, int, string) {
	var ui_errors []error
	user_url, err := getAnswer("\nEnter Url: ")
	ui_errors = append(ui_errors, err)
	user_depth, err := getAnswer("\nEnter Traversal Depth: (2-5 is sensible for most PCs) ")
	ui_errors = append(ui_errors, err)
	user_filename, err := getAnswer("\nName of Output .txt file: ")
	ui_errors = append(ui_errors, err)
	for _, err := range ui_errors {
		if err != nil {
			log.Fatalln("Error: ", err)
		}
	}
	isLive, err := urlIsLive(user_url)
	if err != nil {
		log.Fatalln(err)
	}
	if !isLive {
		log.Fatalln("Error: url is not live")
	}

	intdepth, err := isCompatibleNumber(user_depth)
	if err != nil {
		log.Fatalln(err)
	}

	if !isValidFileName(user_filename) {
		log.Fatalln("Error: filename must be alphanumeric")
	}
	return user_url, intdepth, user_filename

}
