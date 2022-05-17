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
	White int `json:"white"`
	Draws int `json:"draws"`
	Black int `json:"black"`
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
	castleCorrect(&data)
	return data, nil
}

func castleCorrect(data *PositionDataResults){
	n := len(data.Moves)
	for i:=0 ; i<n ; i++ {
		if data.Moves[i].San == "O-O" || data.Moves[i].San == "O-O-O" || data.Moves[i].San == "O-O-O+" || data.Moves[i].San == "O-O+" {
			uci := data.Moves[i].Uci
			switch uci{
			case "e1h1":
				data.Moves[i].Uci = "e1g1"
			case "e8h8":
				data.Moves[i].Uci = "e8g8"
			case "e1a1":
				data.Moves[i].Uci = "e1c1"
			case "e8a8":
				data.Moves[i].Uci = "e8c8"
			}
		}
	}
}
