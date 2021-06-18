package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const authKey = "ec093dd5-bbe3-4d8e-bdac-314b40afb796"
const dailyUrl = "https://de-challenge.ltvco.com/v1/songs/daily"
const monthly = "https://de-challenge.ltvco.com/v1/songs/monthly"

type Repository struct {
	Client *http.Client
}

func New() *Repository {
	return &Repository{Client: &http.Client{}}
}

type SongResponse struct {
	Song_id     string `json:"song_id"`
	Released_at string `json:"released_at"`
	Duration    string `json:"duration"`
	Artist      string `json:"artist"`
	Name        string `json:"name"`
	Stats       Stats  `json:"stats"`
}

type Stats struct {
	Last_played_at float64 `json:"last_played_at"`
	Times_played   int32   `json:"times_played"`
	Global_rank    int32   `json:"global_rank"`
}

//YYYY-MM-DD
func (r *Repository) Daily(releasedAt string) ([]SongResponse, error) {
	args := fmt.Sprintf("{api_key: %s, released_at: %s}", authKey, releasedAt)
	resp, err := r.httpCall("GET", dailyUrl, strings.NewReader(args), releasedAt)

	if err != nil {
		return nil, err
	}

	var list []SongResponse
	err = json.Unmarshal([]byte(resp), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

//YYYY-MM
func (r *Repository) Monthly(released_at string) ([]SongResponse, error) {
	//args := fmt.Sprintf("{api_key: %s, released_at: %s}", authKey, released_at)

	resp, err := r.httpCall("GET", monthly, strings.NewReader(""), released_at)
	if err != nil {
		return nil, err
	}

	var list []SongResponse
	err = json.Unmarshal([]byte(resp), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) httpCall(action, url string, msg *strings.Reader, releasedAt string) (string, error) {
	req, err := http.NewRequest(action, url, msg)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("api_key", authKey)
	q.Add("released_at", releasedAt)
	req.URL.RawQuery = q.Encode()

	resp, err := r.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)

		return "", err
	}

	return string(bytes), nil
}
