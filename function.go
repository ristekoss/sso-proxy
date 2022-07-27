package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gocolly/colly/v2"
	ssojwt "github.com/ristekoss/golang-sso-ui-jwt"
)

func init() {
	functions.HTTP("Proxy", Proxy)
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Request
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formData := map[string]string{
		"username": data.Username,
		"password": data.Password,
	}

	c := colly.NewCollector(colly.AllowURLRevisit())
	c.OnHTML("input[type=hidden]", func(e *colly.HTMLElement) {
		formData[e.Attr("name")] = e.Attr("value")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

  // TODO need to replace hardcoded url
	c.Visit("https://sso.ui.ac.id/cas2/login?service=http%3A%2F%2Flocalhost%3A8081%2F")
	c.Wait()
	fmt.Println(formData)

	c = c.Clone()
	c.OnResponse(func(r *colly.Response) {
		for _, val := range r.Headers.Values("Set-Cookie") {
			header := http.Header{}
			header.Add("Set-Cookie", val)
			req := http.Response{Header: header}
			c.SetCookies(r.Request.URL.String(), req.Cookies())
		}

		fmt.Println(r.Headers.Values("Set-Cookie"))
	})

	c.Post("https://sso.ui.ac.id/cas2/login?service=http%3A%2F%2Flocalhost%3A8081%2F", formData)
	c.Wait()

	var ticket string
	c = c.Clone()
	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		r, _ := http.NewRequest(http.MethodGet, req.URL.String(), nil)
		ticket = r.URL.Query()["ticket"][0]
		fmt.Println("ticket", ticket)
		return nil
	})
	c.Visit("https://sso.ui.ac.id/cas2/login?service=http%3A%2F%2Flocalhost%3A8081%2F")
	c.Wait()

	config := ssojwt.MakeSSOConfig(0, 0, "", "", "http%3A%2F%2Flocalhost%3A8081%2F", "")
	bodyBytes, _ := ssojwt.ValidatTicket(config, ticket)
	model, _ := ssojwt.Unmarshal(bodyBytes)
	res, _ := json.Marshal(model)

	fmt.Fprintln(w, string(res))
}
