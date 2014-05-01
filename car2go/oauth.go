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
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
)

// var tokens map[string]*oauth.RequestToken
// var c *oauth.Consumer

// var store = sessions.NewCookieStore([]byte("***REMOVED***"), []byte("***REMOVED***"))

var client *Client
var testMode string

const (
	sessionName = "car2go-session"
)

type Client struct {
	sessions *sessions.CookieStore
	consumer *oauth.Consumer
	tokens   map[string]*oauth.RequestToken
}

type Account struct {
	AccountID   interface{}
	Description interface{}
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

func (c *Client) SetTestMode(mode string) {
	testMode = mode
}

func init() {
	testMode = "1"

	http.HandleFunc("/oauth/car2go/login", RedirectUserToCar2Go)
	http.HandleFunc("/oauth/car2go", GetCar2GoToken)
	http.HandleFunc("/oauth/car2go/user", Car2GoUserAccount)
	http.HandleFunc("/oauth/car2go/book", Car2GoCreateBooking)
	http.HandleFunc("/oauth/car2go/bookings", Car2GoGetBookings)
	http.HandleFunc("/oauth/car2go/bookings/delete", Car2GoDelateBooking)
}

func RedirectUserToCar2Go(w http.ResponseWriter, r *http.Request) {
	//car2go currently use oob as callback url
	tokenUrl := fmt.Sprintf("http://%s/maketoken", r.Host)
	token, requestUrl, err := client.consumer.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	client.tokens[token.Token] = token
	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

func Car2GoUserAccount(w http.ResponseWriter, r *http.Request) {
	session, _ := client.sessions.Get(r, sessionName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	b, _ := json.Marshal(Account{session.Values["AccountId"], session.Values["Description"]})
	fmt.Fprint(w, string(b))
}

func Car2GoGetBookings(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	loc := values.Get("loc")
	accessToken, _ := getAccessTokenAndAccountId(r)
	response, err := client.consumer.Get(
		"https://www.car2go.com/api/v2.1/booking",
		map[string]string{"loc": loc, "format": "json", "test": testMode},
		accessToken)

	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	bits, _ := ioutil.ReadAll(response.Body)
	fmt.Fprint(w, string(bits))
}

func Car2GoDelateBooking(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	bookingId := values.Get("bookingId")
	accessToken, _ := getAccessTokenAndAccountId(r)
	response, err := client.consumer.Delete(
		"https://www.car2go.com/api/v2.1/booking/"+bookingId,
		map[string]string{"format": "json", "test": testMode},
		accessToken)

	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	bits, _ := ioutil.ReadAll(response.Body)
	fmt.Fprint(w, string(bits))
}

func Car2GoCreateBooking(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	vin := values.Get("vin")
	loc := values.Get("loc")
	accessToken, accountId := getAccessTokenAndAccountId(r)
	response, err := client.consumer.Post(
		"https://www.car2go.com/api/v2.1/bookings",
		map[string]string{"loc": loc, "format": "json", "vin": vin, "account": accountId, "test": testMode},
		accessToken)

	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	bits, _ := ioutil.ReadAll(response.Body)
	fmt.Fprint(w, string(bits))
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

	session, _ := client.sessions.Get(r, sessionName)
	session.Values["Token"] = accessToken.Token
	session.Values["Secret"] = accessToken.Secret

	accountId, description := GetAccountData(accessToken)

	session.Values["AccountId"] = accountId
	session.Values["Description"] = description

	session.Save(r, w)

	log.Printf("%#v", accessToken)
	log.Println(accountId, description)
	http.Redirect(w, r, "/?login=success", http.StatusTemporaryRedirect)
}

func GetAccountData(accessToken *oauth.AccessToken) (accountId int, description string) {

	response, err := client.consumer.Get(
		"https://www.car2go.com/api/v2.1/accounts",
		map[string]string{"loc": "milano", "format": "json", "test": testMode},
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

func getAccessTokenAndAccountId(r *http.Request) (accessToken *oauth.AccessToken, accountId string) {
	session, _ := client.sessions.Get(r, sessionName)
	accountId = strconv.Itoa(session.Values["AccountId"].(int))
	accessToken = &oauth.AccessToken{session.Values["Token"].(string), session.Values["Secret"].(string), map[string]string{}}
	return
}
