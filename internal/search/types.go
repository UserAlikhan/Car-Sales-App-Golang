package search

type CarPostSearchResponse struct {
	Shards   Shards      `json:"_shards"`
	Hits     HitsWrapper `json:"hits"`
	TimedOut bool        `json:"timed_out"`
	Took     int         `json:"took"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// Outer hits object
type HitsWrapper struct {
	Total    Total      `json:"total"`
	MaxScore float64    `json:"max_score"`
	Hits     []HitEntry `json:"hits"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitEntry struct {
	ID     string     `json:"_id"`
	Index  string     `json:"_index"`
	Score  float64    `json:"_score"`
	Source CarPostDoc `json:"_source"`
}
