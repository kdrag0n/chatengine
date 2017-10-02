package main

import (
	"math/rand"
	"time"

	"chatengine/levenshtein"
	"github.com/OneOfOne/cmap/stringcmap"
)

var (
	queryBestCache = stringcmap.NewSize(64) // TODO: genx dedicated cmap KT=string,VT=*DistanceResults
)

type distanceResults struct {
	bestMessages   []*ChatMessage
	respTo         *ResponseFor
	respConfidence float64
	creationTime   time.Time
}

func distanceCacheJanitor() {
	for range time.Tick(time.Minute * 2) {
		now := time.Now()

		queryBestCache.ForEach(func(key string, val interface{}) bool {
			results := val.(*distanceResults)

			if now.Sub(results.creationTime).Hours() > 2.9 {
				queryBestCache.Delete(key)
			}

			return true
		})

		if queryBestCache.Len() > 6144 {
			queryBestCache = stringcmap.NewSize(64)
		}
	}
}

// DistanceAdapter returns the best response based on distance comparisions.
func DistanceAdapter(ctx *QueryCtx) PossibleResponse {
	var results distanceResults
	cacheRes := queryBestCache.Get(ctx.Query)

	if cacheRes != nil {
		results = *(cacheRes.(*distanceResults))
	} else {
		var bestMessages []*ChatMessage
		respTo := &ResponseFor{}
		var respConfidence float64
		respConfidence = 0.0
		queryRunes := ctx.QueryRunes
		qLen := len(queryRunes)
		notCtxIsCJK := !ctx.isCJK

		for e := ctx.bot.data.Messages.Front(); e != nil; e = e.Next() {
			msg := e.Value.(*ChatMessage)
			if msg.IsCJK && notCtxIsCJK {
				continue
			}

			if len(msg.InResponseTo) < 1 {
				continue
			}

			for _, respFor := range msg.InResponseTo {
				tLen := len(respFor.textRunes)

				distance := levenshtein.Distance(queryRunes, respFor.textRunes, qLen, tLen)
				var confidence float64

				if distance == 0 {
					confidence = 1.0
				} else {
					confidence = 1.0 - (float64(distance) / float64(argMax(qLen, tLen)))
				}

				if confidence > respConfidence {
					bestMessages = make([]*ChatMessage, 1, 4)
					bestMessages[0] = msg

					respTo = respFor
					respConfidence = confidence
				} else if confidence == respConfidence {
					bestMessages = append(bestMessages, msg)
				}
			}
		}

		results = distanceResults{
			bestMessages:   bestMessages,
			respTo:         respTo,
			respConfidence: respConfidence,
			creationTime:   time.Now(),
		}
	}
	queryBestCache.Set(ctx.Query, &results)

	if len(results.bestMessages) < 1 {
		return PossibleResponse{
			response:   "",
			confidence: 0.0,
		}
	}

	bestMessage := results.bestMessages[rand.Intn(len(results.bestMessages))]
	confidence := results.respConfidence - (float64(rand.Intn(360)) / 1000.0)
	if confidence > 0.5 {
		results.respTo.Occurrences++

		if ctx.Query != results.respTo.Text {
			bestMessage.InResponseTo = append(bestMessage.InResponseTo, &ResponseFor{
				Text:        ctx.Query,
				Occurrences: 1,
			})
		}
	}

	return PossibleResponse{
		response:    bestMessage.Text,
		confidence:  confidence,
		chatMessage: bestMessage,
	}
}
