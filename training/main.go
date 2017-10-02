package main

import (
	"container/list"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/abiosoft/ishell"
	"gopkg.in/cheggaaa/pb.v2"
	"gopkg.in/mgo.v2/bson"
)

type intTrainer func(*ishell.Context)

var (
	trainers = []intTrainer{
		chatlogTrainer,
		ubuntuDialogsTrainer,
	}
	haveData = false
	data     DataRoot
	rawData  RawDataRoot
)

func addCommands(shell *ishell.Shell, cmds ...*ishell.Cmd) {
	for _, cmd := range cmds {
		shell.AddCmd(cmd)
	}
}

func rawToData(rawData RawDataRoot) DataRoot {
	newMsgList := list.New()
	for _, msg := range rawData.Messages {
		newMsgList.PushBack(msg)
	}

	newData := DataRoot{
		Messages: newMsgList,
	}

	return newData
}

func main() {
	shell := ishell.New()
	startTime := time.Now()

	addCommands(shell, &ishell.Cmd{
		Name: "load",
		Help: "load data from the specified file",
		Func: func(c *ishell.Context) {
			c.Print("File path: ")
			path := c.ReadLine()

			c.Println("Loading data...")
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()
			bytes, err := ioutil.ReadFile(path)
			c.ProgressBar().Stop()
			if err != nil {
				c.Println(err)
				return
			}

			c.Println("Decoding data...")
			c.ProgressBar().Start()
			err = bson.Unmarshal(bytes, &rawData)
			c.ProgressBar().Stop()
			if err != nil {
				c.Println("Error decoding data:", err)
				return
			}
			data = rawToData(rawData)
			haveData = true

			c.Println("Data loaded.")
		},
	}, &ishell.Cmd{
		Name: "new",
		Help: "create a new chat data container",
		Func: func(c *ishell.Context) {
			c.Print("Are you sure you want to create a new, empty data container?\nType (y/n): ")
			choice := c.ReadLine()
			if choice != "y" {
				c.Println("Not proceeding.")
				return
			}

			data = DataRoot{
				Messages: list.New(),
			}
			haveData = true

			c.Println("New data container created.")
		},
	}, &ishell.Cmd{
		Name: "write",
		Help: "write current chat data to disk",
		Func: func(c *ishell.Context) {
			c.Print("Are you sure you want to write the current chat data to disk?\nType (y/n): ")
			choice := c.ReadLine()
			if choice != "y" {
				c.Println("Not proceeding.")
				return
			}

			c.Print("File path: ")
			path := c.ReadLine()

			c.Println("Converting data...")

			msgList := *data.Messages
			msgSlice := make([]*ChatMessage, msgList.Len())

			for i, e := 0, msgList.Front(); e != nil; i, e = i+1, e.Next() {
				msgSlice[i] = e.Value.(*ChatMessage)
			}

			c.Println("Encoding data...")
			bytes, err := bson.Marshal(RawDataRoot{
				Messages: msgSlice,
			})
			if err != nil {
				c.Println(err)
				return
			}

			c.Println("Verifying data...")
			var decodedRoot RawDataRoot
			err = bson.Unmarshal(bytes, &decodedRoot)
			if err != nil {
				c.Println("Error decoding data:", err)
				return
			}
			verifyMsgs := decodedRoot.Messages

			res := verify(c, verifyMsgs)
			if !res {
				c.Print("Write anyway? (y/n): ")
				if c.ReadLine() != "y" {
					return
				}
			}

			c.Println("Writing data to file...")
			c.ProgressBar().Indeterminate(true)
			c.ProgressBar().Start()

			err = ioutil.WriteFile(path, bytes, 644)
			if err != nil {
				c.ProgressBar().Stop()
				c.Println(err)
			} else {
				c.ProgressBar().Stop()
				c.Println("Data written.")
			}
		},
	}, &ishell.Cmd{
		Name: "verify",
		Help: "verify the current chat data",
		Func: func(c *ishell.Context) {
			if !haveData {
				c.Println("There must be a chat data container!")
				return
			}

			c.Println("Converting to slice...")
			msgList := *data.Messages
			msgSlice := make([]*ChatMessage, msgList.Len())

			for i, e := 0, msgList.Front(); e != nil; i, e = i+1, e.Next() {
				msgSlice[i] = e.Value.(*ChatMessage)
			}

			c.Println("Verifying...")
			verify(c, msgSlice)
		},
	}, &ishell.Cmd{
		Name: "count",
		Help: "get the number of messages in the data",
		Func: func(c *ishell.Context) {
			c.Printf("There are %d messages.\n", data.Messages.Len())
		},
	}, &ishell.Cmd{
		Name: "chatlog",
		Help: "train from a chatlog file",
		Func: chatlogTrainer,
	}, &ishell.Cmd{
		Name: "ubuntu",
		Help: "train from ubuntu dialogs tsv data",
		Func: ubuntuDialogsTrainer,
	}, &ishell.Cmd{
		Name: "uptime",
		Help: "get the uptime of this trainer shell",
		Func: func(c *ishell.Context) {
			c.Printf("This shell has been running for %s.\n", time.Since(startTime))
		},
	}, &ishell.Cmd{
		Name: "random",
		Help: "get a random message",
		Func: func(c *ishell.Context) {
			if !haveData {
				c.Println("There must be a chat data container!")
				return
			}
			messages := data.Messages
			if messages.Len() < 1 {
				c.Println("There are no messages!")
				return
			}

			targetIdx := rand.Intn(messages.Len())
			elem := getElemAt(targetIdx, messages)
			if elem == nil {
				c.Println("Found a... non existent message?! Index:", targetIdx)
				return
			}
			msg := elem.Value.(*ChatMessage)
			if msg == nil {
				c.Println("Found a nil message! Index:", targetIdx)
			} else {
				c.Println(msg)
			}
		},
	}, &ishell.Cmd{
		Name: "search",
		Help: "search for a message in the data",
		Func: func(c *ishell.Context) {
			c.Print("Search for: ")
			want := c.ReadLine()
			if want == "" {
				c.Println("You must specify a message to search for!")
				return
			}

			idx := 0
			bar := pb.Full.Start(data.Messages.Len())
			for e := data.Messages.Front(); e != nil; e = e.Next() {
				msg := e.Value.(*ChatMessage)
				if msg == nil {
					var prev *ChatMessage
					prevElem := e.Prev()
					if prevElem == nil {
						prev = &ChatMessage{
							Text: "none",
						}
					} else {
						prev = prevElem.Value.(*ChatMessage)
					}
					if prev == nil {
						prev = &ChatMessage{
							Text: "[nil]",
						}
					}

					var next *ChatMessage
					nextElem := e.Next()
					if nextElem == nil {
						next = &ChatMessage{
							Text: "none",
						}
					} else {
						next = nextElem.Value.(*ChatMessage)
					}
					if next == nil {
						next = &ChatMessage{
							Text: "[nil]",
						}
					}

					c.Printf("Warning: a message is nil! Index: %d\n[Previous: %s, Next: %s]\n", idx, prev.Text, next.Text)
					continue
				}

				if msg.Text == want {
					ctxPrintfAntivet(c, "%d: %s", idx, msg)
				}

				idx++
				bar.Increment()
			}

			bar.Finish()
			c.Println("Search finished.")
		},
	}, &ishell.Cmd{
		Name: "template",
		Help: "evaluate a template string",
		Func: func(c *ishell.Context) {
			c.Print("Template: ")
			tmp := c.ReadLine()
			if tmp == "" {
				c.Println("You must specify a template string!")
				return
			}

			tmpl, err := template.New("trainer_eval").Parse(tmp)
			if err != nil {
				c.Println("Error parsing template:", err)
				return
			}

			err = tmpl.Execute(os.Stdout, map[string]interface{}{
				"timeNow":              time.Now,
				"listNew":              list.New,
				"pbFull":               pb.Full,
				"bsonMarshal":          bson.Marshal,
				"bsonUnmarshal":        bson.Unmarshal,
				"ioutilReadFile":       ioutil.ReadFile,
				"ioutilWriteFile":      ioutil.WriteFile,
				"randIntn":             rand.Intn,
				"randFloat64":          rand.Float64,
				"osOpen":               os.Open,
				"osMkdir":              os.Mkdir,
				"osStdout":             os.Stdout,
				"osStdin":              os.Stdin,
				"osStderr":             os.Stderr,
				"templateNew":          template.New,
				"ishellNew":            ishell.New,
				"ctx":                  c,
				"c":                    c,
				"data":                 data,
				"rawData":              rawData,
				"ubuntuDialogsTrainer": ubuntuDialogsTrainer,
				"verify":               verify,
				"ubuntuTrainRows":      ubuntuTrainRows,
				"getElemAt":            getElemAt,
				"haveData":             haveData,
			})
			if err != nil {
				c.Println("Error executing template:", err)
			}
		},
	}, &ishell.Cmd{
		Name: "get",
		Help: "get a message by its index",
		Func: func(c *ishell.Context) {
			c.Print("Index: ")
			idxStr := c.ReadLine()
			if idxStr == "" {
				c.Println("You must specify an index!")
				return
			}
			parsed, err := strconv.ParseInt(idxStr, 10, 32)
			if err != nil {
				c.Println("Error parsing index:", err)
				return
			}

			elem := getElemAt(int(parsed), data.Messages)
			if elem == nil {
				c.Println("No such message index!")
				return
			}
			msg := elem.Value.(*ChatMessage)
			if msg == nil {
				c.Println("That message is nil!")
				return
			}

			c.Println(msg)
		},
	}, &ishell.Cmd{
		Name: "delete",
		Help: "delete a message, by index",
		Func: func(c *ishell.Context) {
			c.Print("Index: ")
			idxStr := c.ReadLine()
			if idxStr == "" {
				c.Println("You must specify an index!")
				return
			}
			parsed, err := strconv.ParseInt(idxStr, 10, 32)
			if err != nil {
				c.Println("Error parsing index:", err)
				return
			}

			elem := getElemAt(int(parsed), data.Messages)
			if elem == nil {
				c.Println("No such message index!")
				return
			}

			data.Messages.Remove(elem)
			c.Println("Message deleted.")
		},
	}, &ishell.Cmd{
		Name: "surroundings",
		Help: "Get a message by index, and its surroundings",
		Func: func(c *ishell.Context) {
			c.Print("Index: ")
			idxStr := c.ReadLine()
			if idxStr == "" {
				c.Println("You must specify an index!")
				return
			}
			parsed, err := strconv.ParseInt(idxStr, 10, 32)
			if err != nil {
				c.Println("Error parsing index:", err)
				return
			}

			e := getElemAt(int(parsed), data.Messages)
			if e == nil {
				c.Println("No such message index!")
				return
			}
			msg := e.Value.(*ChatMessage)
			if msg == nil {
				msg = &ChatMessage{
					Text: "[nil]",
				}
			}

			var prev *ChatMessage
			prevElem := e.Prev()
			if prevElem == nil {
				prev = &ChatMessage{
					Text: "none",
				}
			} else {
				prev = prevElem.Value.(*ChatMessage)
			}
			if prev == nil {
				prev = &ChatMessage{
					Text: "[nil]",
				}
			}

			var next *ChatMessage
			nextElem := e.Next()
			if nextElem == nil {
				next = &ChatMessage{
					Text: "none",
				}
			} else {
				next = nextElem.Value.(*ChatMessage)
			}
			if next == nil {
				next = &ChatMessage{
					Text: "[nil]",
				}
			}

			ctxPrintfAntivet(c, "Previous: %s\nTarget: %s\nNext: %s\n", prev, msg, next)
		},
	})

	shell.Println("Welcome to the ChatEngine interactive trainer.")
	if len(os.Args) > 1 && os.Args[1] == "exit" {
		shell.Process(os.Args[2:]...)
	} else {
		// start shell
		shell.Run()
	}
}
