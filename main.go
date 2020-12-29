package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Settings struct {
	Credentials struct {
		AccessKeyID     string `json:"AccessKeyID"`
		SecretAccessKey string `json:"SecretAccessKey"`
	} `json:"credentials"`
	FunctionName string `json:"functionName"`
	File         string `json:"file"`
	Region       string `json:"region"`
}

func main() {

	var settings Settings

	// Open settings
	settingsBytes, settingsErr := getFileBytes("./settings.json")
	if settingsErr != nil {
		fmt.Println("ERROR Opening settings.json: " + settingsErr.Error())
		return
	}

	// Unmarshal into struct
	_ = json.Unmarshal(settingsBytes, &settings)

	// Open zip file to upload
	zipBytes, zipErr := getFileBytes(settings.File)
	if zipErr != nil {
		fmt.Println("ERROR Opening Zip File " + settings.File + ": " + zipErr.Error())
		return
	}

	// Create AWS Session
	creds := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     settings.Credentials.AccessKeyID,
		SecretAccessKey: settings.Credentials.SecretAccessKey,
	})
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(settings.Region),
			Credentials: creds,
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("ERROR CREATING AWS SESSION: " + err.Error())
		return
	}

	// Connect session to Lambda service
	sessConn := session.Must(sess, err)
	lamService := lambda.New(sessConn)

	// Setup lambda function for updating
	codeUpdateInput := lambda.UpdateFunctionCodeInput{
		FunctionName: &settings.FunctionName,
		ZipFile:      zipBytes,
	}

	// Update lambda function
	_, err = lamService.UpdateFunctionCode(&codeUpdateInput)
	if err != nil {
		fmt.Println("ERROR UPDATING LAMBDA FUNCTION: " + err.Error())
		return
	}

	fmt.Println("Lambda Function " + settings.FunctionName + " Has Been Successfully Updated!")
}

// getFileBytes ...
func getFileBytes(file string) (fileBytes []byte, err error) {
	// Open File
	reader, openErr := os.Open(file)
	if openErr != nil {
		return fileBytes, openErr
	}
	defer reader.Close()

	// Read file
	fileBytes, err = ioutil.ReadAll(reader)

	return fileBytes, err
}
