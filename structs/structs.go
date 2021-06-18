package structs

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

// [
// {
// "song_id": "a7d8feae-cac5-40c2-8272-53b4089636c7",
// "released_at": "2021-01-21",
// "duration": "3m22s",
// "artist": "Weezer",
// "name": "All My Favorite Songs",
// "stats": {
// "last_played_at": 337193736372486642,
// "times_played": 98621,
// "global_rank": 87
// }
// }
// ]
