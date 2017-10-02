package main

import (
	"chatengine/levenshtein"
	"strings"
	"time"
)

const (
	slWorkerCount = 4
)

type selflearnJob struct {
	clientIP         string
	qLength          int
	bot              *ChatBot
	session          *ChatSession
	fmtQuery         string
	qRunes           []rune
	bestResponse     PossibleResponse
	bestResponseText string
	bestMatch        *ChatMessage
}

var (
	slJobQueue   = make(chan *selflearnJob, 24)
	slWorkerPool = make(chan chan *selflearnJob, slWorkerCount)
)

type selflearnWorker struct {
	jobChan chan *selflearnJob
	quit    chan struct{}
}

func selflearnJobDispatcher() {
	for i := 0; i < slWorkerCount; i++ {
		worker := selflearnWorker{
			jobChan: make(chan *selflearnJob),
			quit:    make(chan struct{}),
		}
		worker.start()
	}

	for {
		select {
		case job := <-slJobQueue:
			worker := <-slWorkerPool
			worker <- job
		}
	}
}

func (w selflearnWorker) start() {
	go func() {
		for {
			slWorkerPool <- w.jobChan

			select {
			case job := <-w.jobChan:
				doJob(job)
			case <-w.quit:
				return
			}
		}
	}()
}

func doJob(_job *selflearnJob) {
	defer handlePanic()
	job := *_job

	ipIsSLBlocked := false
	if _, ok := selflearnBlockedIPs[job.clientIP]; ok {
		ipIsSLBlocked = true
	}
	shouldSelflearn := selflearn && job.qLength <= 75 && !ipIsSLBlocked

	var userIRT []*ResponseFor
	if sessionLen := len(job.session.History); sessionLen > 0 {
		oldEntry := job.session.History[sessionLen-1]
		irtMsg := oldEntry.BestMatch
		if irtMsg == nil {
			irtMsg = oldEntry.Message
		}

		userIRT = []*ResponseFor{{
			Text:        irtMsg.Text,
			Occurrences: 1,
		}}
	} else {
		userIRT = make([]*ResponseFor, 0, 0)
	}

	var userMessage *ChatMessage
	for e := job.bot.data.Messages.Front(); e != nil; e = e.Next() {
		val := e.Value.(*ChatMessage)

		if val.Text == job.fmtQuery {
			userMessage = val

			if shouldSelflearn {
				val.Occurrences++
			}
		}
	}

	if userMessage == nil {
		userMessage = &ChatMessage{
			Text:         job.fmtQuery,
			textRunes:    job.qRunes,
			InResponseTo: userIRT,
			CreatedAt:    time.Now(),
			ExtraData:    make([]ExtraData, 0, 0),
			Occurrences:  1,
		}

		if shouldSelflearn && !filterTest(filterPrepNoLower(strings.ToLower(job.fmtQuery))) {
			job.bot.data.Messages.PushBack(userMessage)

			qRunes := job.qRunes
			qLen := len(qRunes)

			queryBestCache.ForEach(func(query string, rawResults interface{}) bool {
				results := rawResults.(*distanceResults)
				iterQuery := []rune(query)
				iqLen := len(iterQuery)

				distance := levenshtein.Distance(iterQuery, qRunes, iqLen, qLen)
				var conf float64

				if distance == 0 {
					conf = 1.0
				} else {
					conf = 1.0 - (float64(distance) / float64(argMax(iqLen, qLen)))
				}

				if conf >= results.respConfidence {
					results.bestMessages = append(results.bestMessages, userMessage)
				}

				return true
			})
		}
	}

	var ourChatMessage *ChatMessage
	if job.bestResponse.chatMessage != nil {
		ourChatMessage = job.bestResponse.chatMessage
		ourChatMessage.Occurrences++
	} else {
		ourChatMessage = &ChatMessage{
			Text:      job.bestResponseText,
			textRunes: []rune(job.bestResponseText),
			InResponseTo: []*ResponseFor{&ResponseFor{
				Text:        userMessage.Text,
				Occurrences: 1,
			}},
			CreatedAt:   userMessage.CreatedAt,
			ExtraData:   make([]ExtraData, 0, 0),
			Occurrences: 1,
		}
	}

	job.session.History = append(job.session.History, HistoryEntry{
		Message:   userMessage,
		BestMatch: job.bestMatch,
		CreatedAt: userMessage.CreatedAt,
	}, HistoryEntry{
		Message:   ourChatMessage,
		BestMatch: ourChatMessage,
		CreatedAt: userMessage.CreatedAt,
	})
}
