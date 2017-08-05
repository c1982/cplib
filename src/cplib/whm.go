package cplib

import (
	"bytes"
	"bytes"
	"log"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"log"
)

//WHM API 1 (https://documentation.cpanel.net/display/SDK/Guide+to+WHM+API+1)
//
type WHM struct {
	Hostname   string
	Username   string
	Password   string
	Port       int
	SkipVerify bool
	Version string
	Debug bool
}

//CreateAcct (https://documentation.cpanel.net/display/SDK/WHM+API+1+Functions+-+createacct)
//This function creates a cPanel account. The function also sets up the new account's domain information.
func (w *WHM) CreateAcct(username,domain string, args Args){

}

func (w *WHM) execute(function string, arguments Args) error {
	var req *http.Request
	var reqUrl string
	
	vals := arguments.Values(w.Version)
	vals.Set("api.version", w.Version)	

	basic := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", w.Username, w.Password)))
	reqUrl = fmt.Sprintf("https://%s:%d/json-api/%s?%s", w.Hostname, w.Port, function, vals.Encode())

	req, err := http.NewRequest("GET", reqUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic))


	client := w.newHTTPClient()
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//TODO: Response body'i al

	if w.Debug {
		log.Println(reqUrl)
		log.Println(resp.Status)
		log.Println(function)
		log.Println(arguments)
		//log.Println(string(bytes))
	}

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