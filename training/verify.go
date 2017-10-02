package main

import (
	"github.com/abiosoft/ishell"
	"gopkg.in/cheggaaa/pb.v2"
)

func verify(c *ishell.Context, verifyMsgs []*ChatMessage) bool {
	bar := pb.Full.Start(len(verifyMsgs))
	for idx, msg := range verifyMsgs {
		if msg == nil {
			var prev *ChatMessage
			prevIdx := idx - 1
			if prevIdx < 0 {
				prev = &ChatMessage{
					Text: "none",
				}
			} else {
				prev = verifyMsgs[prevIdx]
			}
			if prev == nil {
				prev = &ChatMessage{
					Text: "[nil]",
				}
			}

			var next *ChatMessage
			nextIdx := idx + 1
			if nextIdx >= len(verifyMsgs) {
				next = &ChatMessage{
					Text: "none",
				}
			} else {
				next = verifyMsgs[nextIdx]
			}
			if next == nil {
				next = &ChatMessage{
					Text: "[nil]",
				}
			}

			bar.Finish()
			c.Printf("Error: a message is nil! Index: %d\n[Previous: %s, Next: %s]\n", idx, prev.Text, next.Text)
			return false
		}

		bar.Increment()
	}
	bar.Finish()

	return true
}
