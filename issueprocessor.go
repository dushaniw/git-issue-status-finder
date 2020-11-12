package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fmt.Println("starting.......")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "enter_your_git_access_token_here"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Open the file
	csvInputFile, err := os.Open("input.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvInputFile)

	csvOutputFile, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	csvWriter := csv.NewWriter(csvOutputFile)
	i := 1
	// Iterate through the records of input csv file
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		re := regexp.MustCompile("https://github.com/wso2/product-apim/issues/(\\d)+")
		issueLink := re.FindString(record[0])
		if issueLink != "" {
			re := regexp.MustCompile("https://github.com/wso2/product-apim/issues/")
			firstHalf := re.FindString(issueLink)
			if firstHalf != "" {
				issueNo, _ := strconv.Atoi(issueLink[len(firstHalf):])
				issue, _, err := client.Issues.Get(ctx, "wso2", "product-apim", issueNo)
				if err != nil {
					log.Fatal(err)
				}
				issueState := issue.GetState()
				csvRecord := []string{issueLink, issueState}
				err = csvWriter.Write(csvRecord)
				if err != nil {
					log.Fatal("Cannot write to csv file")
				}
				fmt.Printf("%d GitIssue: %s\t status: %s\n", i, issueLink, issueState)
				i = i + 1
			}
		}
	}
	csvWriter.Flush()
	csvOutputFile.Close()
}
