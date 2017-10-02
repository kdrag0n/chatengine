package main

import (
	"bufio"
	"container/list"
	"io"
	"os"
	"sync"
	"time"

	"github.com/abiosoft/ishell"
	"gopkg.in/cheggaaa/pb.v2"
)

func chatlogTrainer(c *ishell.Context) {
	if !haveData {
		c.Println("There must be a chat data container!")
		return
	}

	c.Print("Path to chatlog file: ")
	path := c.ReadLine()

	file, err := os.Open(path)
	if err != nil {
		c.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.Printf("Error stat()ing file: %s\n", err)
		return
	}
	byteCount := stat.Size()

	c.Println("Reading lines and starting tasks...")
	r := bufio.NewReader(file)
	messages := data.Messages
	var _wg sync.WaitGroup
	wg := &_wg
	var prevLine string
	bar := pb.Full.Start64(byteCount)

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			bar.Finish()
			c.Printf("Error reading line from file: %s\n", err)
			return
		}
		lineLen := len(line)
		if lineLen == 1 { // empty
			continue
		}

		wg.Add(1)
		go chatlogTrainLine(messages, prevLine, line[:lineLen-1], wg)
		prevLine = line

		bar.Add(lineLen)
	}
	bar.Finish()

	c.Println("Waiting for tasks to finish...")
	wg.Wait()

	c.Println("Finished!")
}

func chatlogTrainLine(messages *list.List, prev string, line string, wg *sync.WaitGroup) {
	defer wg.Done()

	var foundMessage *ChatMessage
	for e := messages.Front(); e != nil; e = e.Next() {
		message := e.Value.(*ChatMessage)

		if message.Text == line {
			foundMessage = message
			break
		}
	}

	if foundMessage != nil {
		foundMessage.Occurrences++

		if prev != "" {
			irtLine := prev

			incremented := false
			for _, irt := range foundMessage.InResponseTo {
				if irt.Text == irtLine {
					irt.Occurrences++
					incremented = true
				}
			}

			if !incremented {
				irt := &ResponseFor{
					Text:        irtLine,
					Occurrences: 1,
				}

				foundMessage.InResponseTo = append(foundMessage.InResponseTo, irt)
			}
		}
	} else {
		irtLen := 0
		if prev != "" {
			irtLen = 1
		}

		message := &ChatMessage{
			Text:         line,
			InResponseTo: make([]*ResponseFor, irtLen),
			CreatedAt:    time.Now(),
			ExtraData:    make([]ExtraData, 0, 0),
			Occurrences:  1,
		}

		if irtLen == 1 {
			message.InResponseTo[0] = &ResponseFor{
				Text:        prev,
				Occurrences: 1,
			}
		}

		messages.PushBack(message)
	}
}
