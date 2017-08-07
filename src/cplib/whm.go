package cplib

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//WHM API 1 (https://documentation.cpanel.net/display/SDK/Guide+to+WHM+API+1)
type WHM struct {
	Hostname   string
	Username   string
	Password   string
	Port       int
	SkipVerify bool
	Version    string
	Debug      bool
}

//CreateAcct (https://documentation.cpanel.net/display/SDK/WHM+API+1+Functions+-+createacct)
//This function creates a cPanel account. The function also sets up the new account's domain information.
func (w *WHM) CreateAcct(username, domain string, args Args) (out CreateAcctResponse, err error) {

	if args == nil {
		args = Args{}
	}

	args["username"] = username
	args["domain"] = domain

	result, err := w.execute("createacct", args)

	return result.(CreateAcctResponse), err
}

//RemoveAcct (https://documentation.cpanel.net/display/SDK/WHM+API+1+Functions+-+removeacct)
//This function deletes a cPanel or WHM account.
func (w *WHM) RemoveAcct(username string) (out MetadataResponse, err error) {

	args := Args{}
	args["username"] = username

	result, err := w.execute("removeacct", args)

	if err != nil {
		return out, err
	}

	return result.(MetadataResponse), err
}

func (w *WHM) execute(function string, arguments Args) (out interface{}, err error) {
	var req *http.Request
	var reqURL string

	vals := arguments.Values(w.Version)
	vals.Set("api.version", w.Version)

	basic := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", w.Username, w.Password)))
	reqURL = fmt.Sprintf("https://%s:%d/json-api/%s?%s", w.Hostname, w.Port, function, vals.Encode())

	req, err = http.NewRequest("GET", reqURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic))

	client := w.newHTTPClient()
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if w.Debug {
		log.Println(reqURL)
		log.Println(resp.Status)
		log.Println(function)
		log.Println(arguments)
		log.Println(string(bytes))
	}

	json.Unmarshal(bytes, out)

	if err != nil {
		return nil, err
	}

	return out, err
}

func (w *WHM) newHTTPClient() *http.Client {

	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: w.SkipVerify,
		},
	}

	return client
}
