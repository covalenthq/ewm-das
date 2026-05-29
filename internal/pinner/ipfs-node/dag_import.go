package ipfsnode

// dagImportRoot is one NDJSON line returned by POST /api/v0/dag/import.
type dagImportRoot struct {
	Root struct {
		Cid struct {
			Slash string `json:"/"`
		} `json:"Cid"`
		PinErrorMsg string `json:"PinErrorMsg"`
	} `json:"Root"`
}
