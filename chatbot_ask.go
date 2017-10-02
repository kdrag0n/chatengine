package main

import (
	"math/rand"
	"time"
	"container/list"

	"chatengine/util"
	"github.com/gin-gonic/gin"
)

var (
	lastAskTimes = list.New()
	logAskTime = make(chan time.Duration, 10)
	startingResponse = PossibleResponse{
		response:   "¯\\_(ツ)_/¯",
		confidence: -16384.0,
	}
)

// Ask returns a response provided a question to answer.
func (bot *ChatBot) Ask(query string, qLength int, session *ChatSession, request *gin.Context, clientIP string) (string, float64, error) {
	start := time.Now()
	bestResponse := startingResponse
	var bestMatch *ChatMessage
	var err error

	isCJK := util.ContainsCJK(query)
	fmtQuery := util.Format(query, isCJK)
	qRunes := []rune(fmtQuery)
	// println(fmtQuery) // STOPSHIP // TODO: no. // NOTE: this has to go

	if !isCJK && userGreetingRegexp.MatchString(fmtQuery) {
		bestResponse = PossibleResponse{
			response:   greetings[rand.Intn(len(greetings))],
			confidence: 1.0,
		}
	} else {
		queryCtx := &QueryCtx{
			Query:       fmtQuery,
			QueryRunes:  qRunes,
			origQuery:   query,
			isCJK:       isCJK,
			queryLength: qLength,
			bot:         bot,
			session:     session,
			request:     request,
		}

		options := make([]PossibleResponse, len(logicAdapters))

		for i, adapter := range logicAdapters {
			options[i] = adapter(queryCtx)
		}

		for _, resp := range options {
			if resp.confidence > bestResponse.confidence {
				if resp.userBestMatch != nil {
					bestMatch = resp.userBestMatch
				}
				bestResponse = resp
			}
		}
	}

	bestResponseText := bestResponse.response
	if bestResponse.useGetter {
		bestResponseText = bestResponse.responseGetter()
	}

	session.LastModified = time.Now()
	job := &selflearnJob{
		clientIP:         clientIP,
		qLength:          qLength,
		bot:              bot,
		session:          session,
		fmtQuery:         fmtQuery,
		qRunes:           qRunes,
		bestResponse:     bestResponse,
		bestResponseText: bestResponseText,
		bestMatch:        bestMatch,
	}
	slJobQueue <- job
	logAskTime <- time.Since(start)

	return bestResponseText, bestResponse.confidence, err
}

func askTimeLogger() {
	for {
		duration := <-logAskTime

		if lastAskTimes.Len() == 100 {
			lastAskTimes.Remove(lastAskTimes.Front())
		}

		lastAskTimes.PushBack(duration)
	}
}