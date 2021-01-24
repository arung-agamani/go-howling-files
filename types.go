package main

// Message representing json plain string status response
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UserAuthMessage representing json UserAuthObject status response
type UserAuthMessage struct {
	Status string         `json:"status"`
	Data   UserAuthObject `json:"data"`
}

// DirListMessage representing json directory lists status response
type DirListMessage struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// DirTree data structure for directory tree
type DirTree struct {
	Name       string    `json:"name"`
	IsFile     bool      `json:"isFile"`
	CurrentDir string    `json:"currentDir"`
	Child      []DirTree `json:"child"`
}

// DirTreeMessage representing json response for directory tree
type DirTreeMessage struct {
	Status string  `json:"status"`
	Data   DirTree `json:"data"`
}
