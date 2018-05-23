package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type webhook struct {
	Action      string `json:"action"`
	PullRequest struct {
		CreatedAt    string  `json:"created_at"`
		ClosedAt     *string `json:"closed_at"`
		Additions    int     `json:"additions"`
		Deletions    int     `json:"deletions"`
		ChangedFiles int     `json:"changed_files"`
		Base         struct {
			Repo struct {
				Name string `json:"full_name"`
			} `json:"repo"`
		} `json:"base"`
	} `json:"pull_request"`
}

// VerifySignature checks the github signature on the HTTP request
func VerifySignature(headers map[string]string, body string) error {
	rawToken := os.Getenv("WEBHOOK_SECRET_TOKEN")
	signature, err := GetSignature(headers)
	if err != nil {
		return errors.New("Missing signature")
	}
	mac := hmac.New(sha1.New, []byte(rawToken))
	mac.Write([]byte(body))
	expectedMAC := mac.Sum(nil)
	expectedHubSignature := fmt.Sprintf("sha1=%s", hex.EncodeToString(expectedMAC))

	if signature != expectedHubSignature {
		return errors.New("Signature did not match")
	}
	return nil
}

// IsPullRequest filters out non-PR webhooks.
func IsPullRequest(headers map[string]string) bool {
	for k, v := range headers {
		if strings.ToLower(k) == "x-github-event" {
			if v == "pull_request" {
				return true
			}
			fmt.Printf("ignoring X-GitHub-Event %s\n", v)
		}
	}
	return false
}

// GetSignature retrieves the Github signature from the headers, ignoring case
func GetSignature(headers map[string]string) (string, error) {
	for k, v := range headers {
		if strings.ToLower(k) == "x-hub-signature" {
			return v, nil
		}
	}
	return "", errors.New("Didn't find request signature in headers")
}

// Respond wraps an API Gateway response to be less verbose.
func Respond(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: message, StatusCode: status}, nil
}

func decodePayload(headers map[string]string, body string) (events.APIGatewayProxyResponse, webhook, error) {
	hook := webhook{}
	err := VerifySignature(headers, body)
	result := events.APIGatewayProxyResponse{}
	if err != nil {
		result.StatusCode = 400
		result.Body = "Invalid Signature"
		return result, hook, errors.New("Invalid Signature")
	}

	result.StatusCode = 200
	result.Body = "Ignored notification. LOVE, TIME, IDEAS..."
	return result, hook, errors.New("Ignore notification")
}

// HandleRequest is the main entry point for the lambda processing.
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := Respond

	response, _, err := decodePayload(request.Headers, request.Body)
	if err != nil {
		// don't return an error from the lambda handler even if we are returning
		// an error code in the HTTP response.
		fmt.Println(err)
		return response, nil
	}

	return r(200, "The heart of metaphor is inference.")
}

func main() {
	lambda.Start(HandleRequest)
}
