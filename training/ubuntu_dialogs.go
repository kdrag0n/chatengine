package main

import (
	"archive/tar"
	"compress/gzip"
	"container/list"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/abiosoft/ishell"
	"gopkg.in/cheggaaa/pb.v2"
)

func ubuntuDialogsTrainer(c *ishell.Context) {
	if !haveData {
		c.Println("There must be a chat data container!")
		return
	}

	c.Print("Path to Ubuntu Dialogs tar.gz file: ")
	dialogsPath := c.ReadLine()

	dialogsFile, err := os.Open(dialogsPath)
	if err != nil {
		c.Printf("Error opening file: %s\n", err)
		return
	}
	defer dialogsFile.Close()

	gz, err := gzip.NewReader(dialogsFile)
	if err != nil {
		c.Printf("Error decompressing file: %s\n", err)
		return
	}

	tar := tar.NewReader(gz)

	c.Println("Reading data and launching tasks...")
	messages := data.Messages
	var _wg sync.WaitGroup
	wg := &_wg
	bar := pb.Full.Start(1852868)

	for {
		file, err := tar.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			bar.Finish()
			c.Printf("Error reading file from tar: %s\n", err)
			return
		}

		if !strings.HasSuffix(file.Name, ".tsv") {
			continue
		}

		tsv := csv.NewReader(tar)
		tsv.Comma = '\t'
		tsv.LazyQuotes = true
		tsv.FieldsPerRecord = 4
		rows, err := tsv.ReadAll()
		if err != nil {
			bar.Finish()
			c.Printf("Error reading records from file %s: %s\n", file.Name, err)
			return
		}

		wg.Add(1)
		go ubuntuTrainRows(messages, rows, wg)

		bar.Increment()
	}
	bar.Finish()

	c.Println("Waiting for tasks to finish...")
	wg.Wait()

	c.Println("Finished!")
}

func ubuntuTrainRows(messages *list.List, rows [][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	history := list.New()

	for _, row := range rows {
		line := row[3]
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

			if history.Len() > 0 {
				irtLine := history.Back().Value.(string)

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
			if history.Len() > 0 {
				irtLen = 1
			}

			speaker := ExtraData{
				Name:  "speaker",
				Value: row[1],
			}

			var timeData *ExtraData
			rowTime, err := time.Parse("2006-01-02T15:04:05.000Z", row[0])
			if err == nil {
				timeData = &ExtraData{
					Name:  "time",
					Value: rowTime,
				}
			}

			var addressingSpeaker *ExtraData
			if len(strings.TrimSpace(row[2])) > 0 {
				addressingSpeaker = &ExtraData{
					Name:  "addressing_speaker",
					Value: row[2],
				}
			}

			var extraDataSlice []ExtraData
			if timeData != nil && addressingSpeaker != nil {
				extraDataSlice = []ExtraData{speaker, *timeData, *addressingSpeaker}
			} else if timeData != nil {
				extraDataSlice = []ExtraData{speaker, *timeData}
			} else if addressingSpeaker != nil {
				extraDataSlice = []ExtraData{speaker, *addressingSpeaker}
			} else {
				extraDataSlice = []ExtraData{speaker}
			}

			message := &ChatMessage{
				Text:         line,
				InResponseTo: make([]*ResponseFor, irtLen),
				CreatedAt:    time.Now(),
				ExtraData:    extraDataSlice,
				Occurrences:  1,
			}

			if irtLen == 1 {
				message.InResponseTo[0] = &ResponseFor{
					Text:        history.Back().Value.(string),
					Occurrences: 1,
				}
			}

			messages.PushBack(message)
		}

		history.PushBack(line)
	}
}
