package fb_account_kit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	Version = "v1.2"
	APIUrl  = "https://graph.accountkit.com"
)

type Client struct {
	AccountKitUserID        int64
	AccessToken             string
	AppID                   string
	AppSecret               string
	TokenRefreshIntervalSec float64
}

func (client Client) GetMe() (result map[string]interface{}, err error) {
	result, err = getRequest("me", map[string]string{
		"access_token":    client.AccessToken,
		"appsecret_proof": client.appSecretProof(),
	})

	return
}

func (client *Client) fill(params map[string]interface{}) error {
	id, OK := params["id"].(string)
	if !OK {
		return errors.New("Not founded field 'id'")
	}
	client.AccountKitUserID, _ = strconv.ParseInt(id, 10, 64)

	accessToken, OK := params["access_token"].(string)
	if !OK {
		return errors.New("Not founded field 'access_token'")
	}
	client.AccessToken = accessToken

	refreshInterval, OK := params["token_refresh_interval_sec"].(float64)
	if !OK {
		return errors.New("Not founded field 'token_refresh_interval_sec'")
	}
	client.TokenRefreshIntervalSec = refreshInterval

	return nil
}

func (client Client) appSecretProof() string {
	mac := hmac.New(sha256.New, []byte(client.AppSecret))
	mac.Write([]byte(client.AccessToken))

	return hex.EncodeToString(mac.Sum(nil))
}

func (client Client) appSecretToken() string {
	return fmt.Sprintf("AA|%s|%s", client.AppID, client.AppSecret)
}

func CreateClient(authCode, appID, appSecret string) (client Client, err error) {
	client.AppSecret = appSecret
	client.AppID = appID

	var parsedBody map[string]interface{}
	parsedBody, err = getRequest("access_token", map[string]string{
		"grant_type":   "authorization_code",
		"code":         authCode,
		"access_token": client.appSecretToken(),
	})
	err = client.fill(parsedBody)
	return
}

func getRequest(endPoint string, params map[string]string) (parsedBody map[string]interface{}, err error) {
	var resp *http.Response
	resp, err = http.Get(createURL(endPoint, params))

	if err != nil {
		err = errors.New(fmt.Sprintf("Error open facebook API: ", err))
		return
	}

	if resp.StatusCode == http.StatusOK {
		parsedBody, err = parseRequestBody(resp.Body)
		return
	} else {
		err = errors.New("Facebook API error")
		return
	}
}

func createURL(endPoint string, params map[string]string) string {
	url := fmt.Sprintf("%s/%s/%s",
		APIUrl,
		Version,
		endPoint)
	if len(params) > 0 {
		paramsStr := "?"
		for key, val := range params {
			paramsStr += fmt.Sprintf("%s=%s&", key, val)
		}

		url += paramsStr[:len(paramsStr)-1]
	}

	return url
}

func parseRequestBody(io_body io.ReadCloser) (jsonBody map[string]interface{}, err error) {
	body, err := ioutil.ReadAll(io_body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error parse request body: %s", err))
		return
	}

	json.Unmarshal(body, &jsonBody)
	return
}
