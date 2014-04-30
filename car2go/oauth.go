// Similar to the twitter example, but using an HTTP server instead of
// the command line.  Illustrates the need to preserve requestTokens across
// requests.  This does it very dumbly (using a global variable).
package car2go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
)

// var tokens map[string]*oauth.RequestToken
// var c *oauth.Consumer

// var store = sessions.NewCookieStore([]byte("***REMOVED***"), []byte("***REMOVED***"))

var client *Client

type Client struct {
	sessions *sessions.CookieStore
	consumer *oauth.Consumer
	tokens   map[string]*oauth.RequestToken
}

func NewOAuthClient(hashKey []byte, blockKey []byte, consumerKey string, consumerSecret string) *Client {
	client = &Client{
		sessions.NewCookieStore(hashKey, blockKey),
		oauth.NewConsumer(
			consumerKey,
			consumerSecret,
			oauth.ServiceProvider{
				RequestTokenUrl:   "https://www.car2go.com/api/reqtoken",
				AuthorizeTokenUrl: "https://www.car2go.com/api/authorize",
				AccessTokenUrl:    "https://www.car2go.com/api/accesstoken",
			},
		),
		make(map[string]*oauth.RequestToken),
	}
	return client
}

func (c *Client) SetDebug(debug bool) {
	c.consumer.Debug(debug)
}

func init() {

	// tokens = make(map[string]*oauth.RequestToken)

	// var consumerKey *string = flag.String(
	// 	"consumerkey",
	// 	"getacar",
	// 	"Consumer Key from Car2Go. See: https://dev.twitter.com/apps/new")

	// var consumerSecret *string = flag.String(
	// 	"consumersecret",
	// 	"***REMOVED***",
	// 	"Consumer Secret from Car2Go. See: https://dev.twitter.com/apps/new")

	// var port *int = flag.Int(
	// 	"port",
	// 	8889,
	// 	"Port to listen on.")

	// flag.Parse()

	// c = oauth.NewConsumer(
	// 	*consumerKey,
	// 	*consumerSecret,
	// 	oauth.ServiceProvider{
	// 		RequestTokenUrl:   "https://www.car2go.com/api/reqtoken",
	// 		AuthorizeTokenUrl: "https://www.car2go.com/api/authorize",
	// 		AccessTokenUrl:    "https://www.car2go.com/api/accesstoken",
	// 	},
	// )
	// c.Debug(true)
	// token := &oauth.AccessToken{"fXbl3o0LvQItAHdyUcF6OyPbi", "***REMOVED***", map[string]string{}}
	// response, err := c.Get(
	// 	"https://www.car2go.com/api/v2.1/accounts",
	// 	map[string]string{"loc": "milano", "format": "json"},
	// 	token)
	// log.Printf("%#v", response.StatusCode)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer response.Body.Close()

	// bits, err := ioutil.ReadAll(response.Body)
	// log.Printf("%#v", string(bits))
	// return

	http.HandleFunc("/oauth/car2go/login", RedirectUserToCar2Go)
	http.HandleFunc("/oauth/car2go", GetCar2GoToken)
	http.HandleFunc("/oauth/car2go/user", Car2GoUserDetail)
	// u := fmt.Sprintf(":%d", *port)
	// fmt.Printf("Listening on '%s'\n", u)
	// http.ListenAndServe(u, nil)
}

func RedirectUserToCar2Go(w http.ResponseWriter, r *http.Request) {
	//car2go currently use oob as callback url
	tokenUrl := fmt.Sprintf("http://%s/maketoken", r.Host)
	token, requestUrl, err := client.consumer.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		log.Fatal(err)
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	client.tokens[token.Token] = token
	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

func Car2GoUserDetail(w http.ResponseWriter, r *http.Request) {
	session, _ := client.sessions.Get(r, "car2go-session")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprint(w, session.Values["AccountId"], session.Values["Description"])
}

func GetCar2GoToken(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := client.consumer.AuthorizeToken(client.tokens[tokenKey], verificationCode)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	session, _ := client.sessions.Get(r, "car2go-session")
	session.Values["Token"] = accessToken.Token
	session.Values["Secret"] = accessToken.Secret

	accountId, description := GetAccountData(accessToken)

	session.Values["AccountId"] = accountId
	session.Values["Description"] = description

	session.Save(r, w)

	fmt.Fprintf(w, "%°v", accessToken)
	// response, err := c.Get(
	// 	"https://www.car2go.com/api/v2.1/accounts",
	// 	map[string]string{"loc": "milano", "format": "json"},
	// 	accessToken)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer response.Body.Close()

	// bits, err := ioutil.ReadAll(response.Body)
	// fmt.Fprintf(w, string(bits))
	// log.Printf("%#v", accessToken)

	// {"account":[{"accountId":750314,"description":"Andrea Leverano"}],"returnValue":{"code":0,"description":"Operation successful."}}
}

func GetAccountData(accessToken *oauth.AccessToken) (accountId int, description string) {

	response, err := client.consumer.Get(
		"https://www.car2go.com/api/v2.1/accounts",
		map[string]string{"loc": "milano", "format": "json"},
		accessToken)

	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	bits, _ := ioutil.ReadAll(response.Body)

	type ReturnValue struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	}
	type Account struct {
		AccountId   int    `json:"accountId"`
		Description string `json:"description"`
	}

	type AccountData struct {
		Accounts    []Account `json:"account"`
		ReturnValue `json:"returnValue"`
	}

	var ad AccountData
	json.Unmarshal(bits, &ad)
	accountId = ad.Accounts[0].AccountId
	description = ad.Accounts[0].Description
	return
}
