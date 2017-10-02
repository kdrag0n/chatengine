package main

import (
	"container/list"
	"time"
)

// DataRoot represents the root of all the chat data.
type DataRoot struct {
	Messages *list.List
}

// RawDataRoot represents the raw root of read chat data.
type RawDataRoot struct {
	Messages []*ChatMessage `bson:"messages"`
}

// Config represents the API key data.
type Config struct {
	ListenOn            string   `json:"listen"`
	API                 []string `json:"api"`
	Admin               []string `json:"admin"`
	BlockedIPs          []string `json:"blocked_ips"`
	SelflearnBlockedIPs []string `json:"selflearn_blocked_ips"`
	StackimpactKey      string   `json:"stackimpact"`
	WebchatURL          string   `json:"webchat_url"`
	Environment         string   `json:"env"`
}

// ChatMessage represents a single response.
type ChatMessage struct {
	Text         string         `bson:"text"`
	InResponseTo []*ResponseFor `bson:"in_response_to"`
	CreatedAt    time.Time      `bson:"created_at"`
	ExtraData    []ExtraData    `bson:"extra_data"`
	Occurrences  int            `bson:"occurrence"`
}

// ResponseFor represents a message which would get the parent response.
type ResponseFor struct {
	Text        string `bson:"text"`
	Occurrences int    `bson:"occurrence"`
}

// ExtraData represents extra data for a message.
type ExtraData struct {
	Name  string      `bson:"name"`
	Value interface{} `bson:"value"`
}
