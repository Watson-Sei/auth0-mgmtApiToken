package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type HttpResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func main() {

	GetMgmtApiToken()

	c := cron.New()

	c.AddFunc("0 */12 * * *", GetMgmtApiToken)

	c.Start()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		if os.Getenv("AUTH0_MANAGEMENT_API_TOKEN") != "" {
			c.JSON(http.StatusOK, gin.H{
				"auth0_management_api_token": os.Getenv("AUTH0_MANAGEMENT_API_TOKEN"),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to Gin World",
			})
		}
	})
	r.Run()
}

func GetMgmtApiToken() {
	endpoint := os.Getenv("AUTH0_OAUTH_URL")

	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")
	payload.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	payload.Add("client_secret", os.Getenv("AUTH0_CLIENT_SECRET"))
	payload.Add("audience", os.Getenv("AUTH0_AUDIENCE"))

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(payload.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resBody HttpResponse
	if err := json.Unmarshal(body, &resBody); err != nil {
		log.Fatal(err)
	}

	fmt.Println(resBody.AccessToken)
	os.Setenv("AUTH0_MANAGEMENT_API_TOKEN", resBody.AccessToken)
}
