package vo

type PingResult struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Version     struct {
		Number                           string `json:"number"`
		BuildFlavor                      string `json:"build_flavor"`                        // e.g. "oss" or "default"
		BuildType                        string `json:"build_type"`                          // e.g. "docker"
		BuildHash                        string `json:"build_hash"`                          // e.g. "b7e28a7"
		BuildDate                        string `json:"build_date"`                          // e.g. "2019-04-05T22:55:32.697037Z"
		BuildSnapshot                    bool   `json:"build_snapshot"`                      // e.g. false
		LuceneVersion                    string `json:"lucene_version"`                      // e.g. "8.0.0"
		MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`  // e.g. "6.7.0"
		MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"` // e.g. "6.0.0-beta1"
	} `json:"version"`
	TagLine string `json:"tagline"`
}
