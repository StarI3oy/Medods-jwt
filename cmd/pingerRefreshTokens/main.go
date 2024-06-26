package pingerRefreshTokens

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "localhost:8080/refresh?refresh_token=eyJHdWlkIjoiMTIzIiwiRXhwaXJlcyI6IjIwMjQtMDQtMDRUMTM6NDA6NTAuNTYyMjIzNiswNTowMCJ9"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
