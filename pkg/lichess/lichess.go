package lichess

import(
	"encoding/json"
	"io/ioutil"
	"net/url"
	"net/http"
)

type PositionDataResults struct {
	Moves []struct {
		Uci string `json:"uci"`
		San string `json:"san"`
		White int `json:"white"`
		Draws int `json:"draws"`
		Black int `json:"black"`
	} `json:"moves"`
}

func PositionData(fenString string) (PositionDataResults, error) {
	escapedQuery := url.QueryEscape(fenString)
	lichessUrl := "https://explorer.lichess.ovh/lichess?topGames=0&recentGames=0&fen=" + escapedQuery
	resp, err := http.Get(lichessUrl)
	if err != nil {
		return PositionDataResults{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := PositionDataResults{}
	if err := json.Unmarshal(body, &data); err != nil {
		return PositionDataResults{}, err
	}
	return data, nil
}
