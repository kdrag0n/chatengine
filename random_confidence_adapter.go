package main

import (
	"math/rand"
)

// RandomConfidenceAdapter returns a random message with a random low-high confidence.
func RandomConfidenceAdapter(ctx *QueryCtx) PossibleResponse {
	confidence := float64(rand.Intn(650)) / 1000.0
	messages := *ctx.bot.data.Messages

	randomGetter := func() string {
		targetIdx := rand.Intn(messages.Len())
		half := messages.Len() / 2

		if targetIdx < half {
			currentIdx := 0

			for e := messages.Front(); e != nil; e = e.Next() {
				if currentIdx == targetIdx {
					return e.Value.(*ChatMessage).Text
				}

				currentIdx++
			}
		} else {
			currentIdx := messages.Len() - 1

			for e := messages.Back(); e != nil; e = e.Prev() {
				if currentIdx == targetIdx {
					return e.Value.(*ChatMessage).Text
				}

				currentIdx--
			}
		}

		return "A very strange internal error occurred."
	}

	if len(ctx.session.History) >= 6 {
		// TODO: make anti loop here
	}

	return PossibleResponse{
		useGetter:      true,
		responseGetter: randomGetter,
		confidence:     confidence,
	}
}
