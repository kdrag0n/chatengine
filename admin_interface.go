package main

import (
	"bytes"
	"container/list"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"runtime"

	"chatengine/util"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"gopkg.in/mgo.v2/bson"
)

const (
	generalTimeFmt = "Mon, Jan 2 3:04 PM MST"
	preciseTimeFmt = "Mon, Jan 2 3:04:05.000 PM -0700 MST"

	modeExactMatch = "exact_match"
	modeSubstring  = "substring"
	modeRegexp     = "regexp"
)

var (
	needKey = template.HTML(`API keys are only accepted in the <code>adminKey</code> cookie, or the <code>key</code> query string parameter.`)
)

func requireAdminKey() func(*gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/admin/login" {
			ctx.Next()
			return
		}

		key, _ := ctx.Cookie("adminKey")
		if key == "" {
			ctx.HTML(http.StatusUnauthorized, "admin.tmpl", gin.H{
				"title":   "Unauthorized",
				"desc":    "A valid admin key is required.",
				"content": needKey,
			})

			ctx.Abort()
			return
		}

		if _, ok := validAdminKeys[key]; !ok {
			ctx.HTML(http.StatusUnauthorized, "admin.tmpl", gin.H{
				"title":   "Unauthorized",
				"desc":    "A valid admin key is required.",
				"content": needKey,
			})

			ctx.Abort()
		}

		ctx.Set("key", key)
		ctx.Next()
	}
}

func endpointHelp(endpoint, help string) template.HTML {
	return template.HTML(`<a href="` + endpoint + `">/` + endpoint + `</a> - ` + help)
}

func _dcAnyMatchIRT(msg *ChatMessage) bool {
	for _, irt := range msg.InResponseTo {
		if filterTest(filterPrep(irt.Text)) {
			return true
		}
	}

	return false
}

func genMsgCardHTML(cardName string, msg *ChatMessage, index *int) string {
	indexKey := ""
	if index != nil {
		indexKey = fmt.Sprintf(` <form action="get" method="post" enctype="application/x-www-form-urlencoded" style="display:inline;">
<button type="submit" name="index" value="%d" style="padding:0.1em;"><code>%d</code></button>
</form> <form action="delete" method="post" enctype="application/x-www-form-urlencoded" style="display:inline;">
<button type="submit" name="index" value="%d" style="padding:0.1em;">Delete</button>
</form>`, *index, *index, *index)
	}

	var irtList, exdList string
	for _, irt := range msg.InResponseTo {
		_key := "s"
		if irt.Occurrences == 1 {
			_key = ""
		}

		irtList += fmt.Sprintf("<li><code>%s</code> (%d occurrence%s)</li>", template.HTMLEscapeString(irt.Text), irt.Occurrences, _key)
	}
	for _, exd := range msg.ExtraData {
		exdList += fmt.Sprintf("<li><code>%s</code>: <code>%s</code></li>", template.HTMLEscapeString(exd.Name), exd.Value)
	}

	return fmt.Sprintf(`<div style="padding: 0.5em 0.75em;border-radius: 1.5em;background-color: rgb(143, 227, 229);display: inline-block;">
    <h3>%s%s</h3>
    Text: <code>%s</code><br>
    Created at: %s<br>
    Occurrences: %d<br>
    <br>
    In response to:
    <ul>
    %s
    </ul>
    Extra data:
    <ul>
    %s
    </ul>
</div>`, cardName, indexKey, template.HTMLEscapeString(msg.Text), msg.CreatedAt.Format("Mon, Jan 2 3:04:05 PM -0700 MST 2006"),
		msg.Occurrences, irtList, exdList)
}

func setupAdminInterface(bot *ChatBot, router *gin.Engine) {
	group := router.Group("/admin")
	tmplFuncMap := template.FuncMap{
		"timeNow":           time.Now,
		"listNew":           list.New,
		"bsonMarshal":       bson.Marshal,
		"bsonUnmarshal":     bson.Unmarshal,
		"ioutilReadFile":    ioutil.ReadFile,
		"ioutilWriteFile":   ioutil.WriteFile,
		"randIntn":          rand.Intn,
		"randFloat64":       rand.Float64,
		"osOpen":            os.Open,
		"osMkdir":           os.Mkdir,
		"templateNew":       template.New,
		"botLoadChatData":   bot.LoadChatData,
		"botWriteChatData":  bot.WriteChatData,
		"requireAdminKey":   requireAdminKey,
		"endpointHelp":      endpointHelp,
		"genMsgCardHTML":    genMsgCardHTML,
		"_dcAnyMatchIRT":    _dcAnyMatchIRT,
		"format":            util.Format,
		"filterTest":        filterTest,
		"filterPrep":        filterPrep,
		"filterPrepNoLower": filterPrepNoLower,
		"getElemAt":         util.GetElemAt,
		"getFilterWordList": func() string {
			return filterWordList
		},
		"getFilterSet": func() map[string]struct{} {
			return filterSet
		},
	}
	tmplVars := map[string]interface{}{
		"data":                bot.data,
		"messages":            bot.data.Messages,
		"adminGroup":          group,
		"generalTimeFmt":      generalTimeFmt,
		"bot":                 bot,
		"osStdout":            os.Stdout,
		"osStdin":             os.Stdin,
		"osStderr":            os.Stderr,
		"needKey":             needKey,
		"startTime":           startTime,
		"selflearn":           selflearn,
		"autowriteChatData":   autowriteChatData,
		"config":              config,
		"logger":              logger,
		"greetings":           greetings,
		"userGreetings":       userGreetings,
		"userGreetingsRegexp": userGreetingRegexp,
		"chatSessions":        chatSessions,
	}

	help := func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"title":    "ChatEngine Admin Interface",
			"desc":     template.HTML("Accepts basic GET requests.<br>The built-in Go profiler, pprof, is also supported. Use it like <code>go tool pprof chatengine http://chatengine_server/admin/pprof/profile</code>.<br>Pprof's usage: <code>go tool pprof [program binary file] [pprof server url]</code>"),
			"listhead": "Available endpoints:",
			"list": []template.HTML{
				template.HTML(`<a href="#">/</a> or <a href="help">/help</a> - get endpoint help (this page)`),
				endpointHelp("login", "log into the admin interface"),
				endpointHelp("logout", "logout of the admin interface"),
				endpointHelp("dashboard", "experimental new full-on admin dashboard"),
				endpointHelp("keys/get", "get loaded key data"),
				endpointHelp("keys/list", "list valid API keys"),
				endpointHelp("keys/reload", "reload key data from disk"),
				endpointHelp("chat/reload", "reload chat data from disk"),
				endpointHelp("chat/write", "write chat data to disk"),
				endpointHelp("chat/selflearn", "toggle chat self-learning"),
				endpointHelp("chat/autowrite", "toggle auto-writing of chat data to disk"),
				endpointHelp("chat/dataclean", "clean chat data, removing responses caught by the filter"),
				endpointHelp("chat/reformat", `(re)format all chat data, ensuring "proper" punctuation`),
				endpointHelp("chat/clear", "clear the chat data, removing <strong>ALL</strong> messages"),
				endpointHelp("chat/count", "get the number of chat messages stored"),
				endpointHelp("chat/random", "get a random chat message"),
				endpointHelp("chat/search", "search by regex/substring/match, in all messages"),
				endpointHelp("chat/delete", "delete chat message(s) by regex/substring/match/index/index range"),
				endpointHelp("chat/get", "get a chat message by index, including previous and next"),
				endpointHelp("sessions/clean", "remove sessions not used for 15 minutes"),
				endpointHelp("sessions/prune", "clear <strong>ALL</strong> active sessions"),
				endpointHelp("sessions/list", "list currently active session IDs"),
				endpointHelp("sessions/view", "view chat history of a session"),
				endpointHelp("pprof", "pprof, the go profiler"),
				endpointHelp("pprof/*", "endpoints for use by the pprof tool"),
				endpointHelp("errors/panic", "raise a panic for error recovery testing"),
				endpointHelp("errors/panic?err=runtime error: something", "panic with the specified error"),
				endpointHelp("manage/template", "evaluate a template string"),
				endpointHelp("manage/uptime", "get the current uptime"),
				endpointHelp("manage/free_mem", "free unused memory pages to the OS"),
				endpointHelp("manage/stack", "get the full current stack"),
				endpointHelp("manage/gc_stats", "get the garbage collector's stats"),
				endpointHelp("manage/gc_percent", "set the minimum ratio between new and old data (percent) for a GC"),
				endpointHelp("manage/max_threads", "set the maximum number of OS threads the server may create"),
				endpointHelp("manage/heap_dump", "write a full heap dump to a file, or send the data"),
				template.HTML(`<a href="manage/get_heap_dump__nconfirm">/manage/get_heap_dump</a> - download a full heap dump`),
				endpointHelp("manage/sysinfo", "get system information"),
			},
		})
	}

	group.Use(requireAdminKey())
	{
		group.GET("/", help)
		group.GET("/help", help)
		group.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_login.tmpl", gin.H{
				"now": time.Now(),
			})
		})
		group.POST("/login", func(ctx *gin.Context) {
			key := ctx.PostForm("key")
			if len(key) != 96 {
				ctx.HTML(http.StatusBadRequest, "admin_login.tmpl", gin.H{
					"now":   time.Now(),
					"error": "You must enter a valid key!",
				})
				return
			}

			if _, ok := validAdminKeys[key]; !ok {
				ctx.HTML(http.StatusUnauthorized, "admin_login.tmpl", gin.H{
					"now":   time.Now(),
					"error": "That key isn't a valid admin key!",
				})
				return
			}

			ctx.SetCookie("adminKey", key, 31556926, "/", "", false, true)
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Admin Authentication",
				"desc":  "You have been logged into the admin interface.",
			})
		})
		group.GET("/logout", func(ctx *gin.Context) {
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:     "adminKey",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   0,
				HttpOnly: true,
			})

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Admin Authentication",
				"desc":  "You have been logged out of the admin interface.",
			})
		})

		group.GET("/dashboard", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_dashboard.tmpl", nil)
		})

		group.GET("/keys/get", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, config)
		})
		group.GET("/keys/list", func(ctx *gin.Context) {
			keys := make([]string, len(validAPIKeys))
			i := 0
			for key := range validAPIKeys {
				keys[i] = key
				i++
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "API Key List",
				"desc":  "Here is a list of all currently valid API keys.",
				"list":  keys,
			})
		})
		group.GET("/keys/reload", func(ctx *gin.Context) {
			_, err := loadConfig()

			if err != nil {
				ctx.HTML(http.StatusOK, "error.tmpl", gin.H{
					"admin": true,
					"desc":  "Failed to reload key data.",
					"err":   err,
				})
			} else {
				ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
					"title": "Key Data",
					"desc":  "The key data was successfully reloaded.",
				})
			}
		})

		group.GET("/chat/reload", func(ctx *gin.Context) {
			err := bot.LoadChatData()

			if err != nil {
				ctx.HTML(http.StatusOK, "error.tmpl", gin.H{
					"admin": true,
					"desc":  "Failed to reload chat data.",
					"err":   err,
				})
			} else {
				ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
					"title": "Chat Data",
					"desc":  "The chat data was successfully reloaded.",
				})
			}
		})
		group.GET("/chat/write", func(ctx *gin.Context) {
			err := bot.WriteChatData()

			if err != nil {
				ctx.HTML(http.StatusOK, "error.tmpl", gin.H{
					"admin": true,
					"desc":  "Failed to write chat data.",
					"err":   err,
				})
			} else {
				ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
					"title": "Chat Data",
					"desc":  "The chat data was successfully written.",
				})
			}
		})
		group.GET("/chat/selflearn", func(ctx *gin.Context) {
			selflearn = !selflearn

			key := "ff"
			if selflearn {
				key = "n"
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Self-Learning",
				"desc":  "Chat self-learning is now o" + key + ".",
			})
		})
		group.GET("/chat/autowrite", func(ctx *gin.Context) {
			autowriteChatData = !autowriteChatData

			key := "ff"
			if autowriteChatData {
				key = "n"
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Auto-Writing",
				"desc":  "Chat auto-writing is now o" + key + ".",
			})
		})
		group.GET("/chat/dataclean", func(ctx *gin.Context) {
			removed := 0
			toRemove := make([]*list.Element, 0, 1000)
			for e := bot.data.Messages.Front(); e != nil; e = e.Next() {
				msg := e.Value.(*ChatMessage)

				if msg.Text == "" || filterTest(filterPrep(msg.Text)) || _dcAnyMatchIRT(msg) {
					toRemove = append(toRemove, e)
					removed++
				}
			}

			for _, e := range toRemove {
				bot.data.Messages.Remove(e)
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Data",
				"desc":  template.HTML(fmt.Sprintf("The chat data has been scanned for messages violating the filter.<br>%d messages were removed.", removed)),
			})
		})
		group.GET("/chat/reformat", func(ctx *gin.Context) {
			processed := bot.data.Messages.Len()
			modified := 0
			for e := bot.data.Messages.Front(); e != nil; e = e.Next() {
				msg := e.Value.(*ChatMessage)
				oldText := msg.Text
				msg.Text = util.Format(msg.Text, util.ContainsCJK(msg.Text))
				_modified := false

				for _, irt := range msg.InResponseTo {
					oldItext := irt.Text
					irt.Text = util.Format(irt.Text, util.ContainsCJK(irt.Text))

					if oldItext != irt.Text {
						_modified = true
					}
				}

				if _modified || oldText != msg.Text {
					modified++
				}
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Data",
				"desc":  template.HTML(fmt.Sprintf("All stored responses and their <code>InResponseTo</code>s have been verified for format.<br>%d messages were processed. %d messages were modified.", processed, modified)),
			})
		})
		group.GET("/chat/clear", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Data",
				"desc":  template.HTML(`<h2>Are you sure you want to clear <strong>all</strong> the chat data?</h2><br><a href="clear__confirm">Yes, I'm sure.</a><br><a style="font-size: 20pt;" href="../">No, I'm not sure.</a>`),
			})
		})
		group.GET("/chat/clear__confirm", func(ctx *gin.Context) {
			bot.data.Messages = list.New()

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Data",
				"desc":  template.HTML("The chat data has been <strong>cleared</strong>."),
			})
		})
		group.GET("/chat/count", func(ctx *gin.Context) {
			count := bot.data.Messages.Len()

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Data",
				"desc":  template.HTML(fmt.Sprintf("There are <strong>%d</strong> chat messages stored.", count)),
			})
		})
		group.GET("/chat/random", func(ctx *gin.Context) {
			text := "What? No response?"
			defer func() {
				ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
					"title": "Random Chat Message",
					"desc":  template.HTML(text),
				})
			}()

			targetIdx := rand.Intn(bot.data.Messages.Len())
			elem := util.GetElemAt(targetIdx, bot.data.Messages)
			if elem == nil {
				text = fmt.Sprintf("Found a... non existent message?! Index: %d", targetIdx)
				return
			}

			msg := elem.Value.(*ChatMessage)
			if msg == nil {
				text = fmt.Sprintf("Found a nil message! Index: %d", targetIdx)
			} else {
				text = genMsgCardHTML("Message", msg, &targetIdx)
			}
		})
		group.GET("/chat/search", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Chat Message Search",
				"desc": template.HTML(`Here you can search for chat messages.<br>
<br>
<form action="search" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Search by exact text match</h3>
	<input type="text" placeholder="Text to search for..." name="message" required>
	<button type="submit">Search</button>
</form>
<br>
<form action="search" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Search by substring</h3>
	<input type="text" placeholder="Substring to search for..." name="substr" required>
	<button type="submit">Search</button>
</form>
<br>
<form action="search" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Search by regular expression</h3>
	<input type="text" placeholder="Regexp to match with..." name="regex" required>
	<button type="submit">Search</button>
</form>`),
			})
		})
		group.POST("/chat/search", func(ctx *gin.Context) {
			var regex *regexp.Regexp
			var mode, text string
			if text = ctx.PostForm("message"); text != "" {
				mode = modeExactMatch
			} else if text = ctx.PostForm("substr"); text != "" {
				mode = modeSubstring
			} else if text = ctx.PostForm("regex"); text != "" {
				mode = modeRegexp
				regex = regexp.MustCompile(text)
			} else {
				ctx.HTML(http.StatusBadRequest, "admin.tmpl", gin.H{
					"title": "Chat Message Search",
					"desc":  "You must provide a search method and text/input for that search method!",
				})

				return
			}

			results := make([]template.HTML, 0, 10)
			idx := 0
			for e := bot.data.Messages.Front(); e != nil; e = e.Next() {
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

					results = append(results, template.HTML(fmt.Sprintf("<strong>Warning</strong>: a message is <code>nil</code>! Index: %d<br>[Previous: %s, Next: %s]", idx, template.HTMLEscapeString(prev.Text), template.HTMLEscapeString(next.Text))))
					continue
				}

				var cond bool
				switch mode {
				case modeExactMatch:
					cond = msg.Text == text
				case modeSubstring:
					cond = strings.Contains(msg.Text, text)
				case modeRegexp:
					cond = regex.MatchString(msg.Text)
				}

				if cond {
					results = append(results, template.HTML(fmt.Sprintf(`<form action="get" method="post" enctype="application/x-www-form-urlencoded" style="display:inline;">
    <button type="submit" name="index" value="%d" style="padding:0;"><code>%d</code></button>
</form> %s`, idx, idx, template.HTMLEscapeString(msg.Text))))
				}

				idx++
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Message Search Results",
				"desc": template.HTML(fmt.Sprintf(`Search mode: <code>%s</code><br>
Searched for: <code>%s</code><br>
Results found: <strong>%d</strong><br>
<br>
`, mode, text, len(results))),
				"listhead": "Here are the results:",
				"list":     results,
			})
		})
		group.GET("/chat/delete", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Delete Chat Message(s)",
				"desc": template.HTML(`
<form action="delete" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Delete by index</h3>
	<input type="number" placeholder="Enter message index..." name="index" required>
	<button type="submit">Delete</button>
</form>
<br>
<form action="delete" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Delete by exact text match</h3>
	<input type="text" placeholder="Text to search for..." name="message" required>
	<button type="submit">Delete</button>
</form>
<br>
<form action="delete" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Delete by substring</h3>
	<input type="text" placeholder="Substring to search for..." name="substr" required>
	<button type="submit">Delete</button>
</form>
<br>
<form action="delete" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Delete by regular expression</h3>
	<input type="text" placeholder="Regexp to match with..." name="regex" required>
	<button type="submit">Delete</button>
</form>`),
			})
		})
		group.POST("/chat/delete", func(ctx *gin.Context) {
			var regex *regexp.Regexp
			var index int
			var mode, text string

			if text = ctx.PostForm("message"); text != "" {
				mode = modeExactMatch
			} else if text = ctx.PostForm("substr"); text != "" {
				mode = modeSubstring
			} else if text = ctx.PostForm("regex"); text != "" {
				mode = modeRegexp
				regex = regexp.MustCompile(text)
			} else if text = ctx.PostForm("index"); text != "" {
				mode = "index"
				parsedUI64, err := strconv.ParseUint(text, 10, 32)
				if err != nil {
					panic(err)
				}
				index = int(parsedUI64)
			} else {
				ctx.HTML(http.StatusBadRequest, "admin.tmpl", gin.H{
					"title": "Chat Message Search",
					"desc":  "You must provide a search method and text/input for that search method!",
				})

				return
			}

			results := make([]string, 0, 10)
			toDelete := make([]*list.Element, 0, 25)
			idx := 0
			for e := bot.data.Messages.Front(); e != nil; e = e.Next() {
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

					results = append(results, fmt.Sprintf("<strong>Warning</strong>: a message is <code>nil</code>! Index: %d<br>[Previous: %s, Next: %s]", idx, template.HTMLEscapeString(prev.Text), template.HTMLEscapeString(next.Text)))
					continue
				}

				var cond bool
				switch mode {
				case modeExactMatch:
					cond = msg.Text == text
				case modeSubstring:
					cond = strings.Contains(msg.Text, text)
				case modeRegexp:
					cond = regex.MatchString(msg.Text)
				case "index":
					cond = idx == index
				}

				if cond {
					toDelete = append(toDelete, e)
					results = append(results, genMsgCardHTML("Deleted message", msg, &idx))
				}

				idx++
			}

			for _, e := range toDelete {
				bot.data.Messages.Remove(e)
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Delete Chat Message",
				"desc": template.HTML(fmt.Sprintf(`Delete mode: <code>%s</code><br>
Searched for: <code>%s</code><br>
Messages deleted: <strong>%d</strong><br>
<br>%s`, mode, text, len(results), strings.Join(results, "<br>"))),
			})
		})
		group.GET("/chat/get", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Get Chat Message",
				"desc": template.HTML(`
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Message Index:</h3>
	<input type="number" placeholder="Enter message index..." name="index" required>
	<button type="submit">Get</button>
</form>`),
			})
		})
		group.POST("/chat/get", func(ctx *gin.Context) {
			parsedUI64, err := strconv.ParseUint(ctx.PostForm("index"), 10, 32)
			if err != nil {
				panic(err)
			}

			index := int(parsedUI64)
			iPlus1 := index + 1
			iMinus1 := index - 1
			e := util.GetElemAt(index, bot.data.Messages)
			if e == nil {
				ctx.HTML(http.StatusBadRequest, "admin.tmpl", gin.H{
					"title": "Get Chat Message",
					"desc": template.HTML(fmt.Sprintf(`There's no message at index <code>%d</code>!
<br>
Maybe try getting:
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
    <button type="submit" name="index" value="%d">
        The previous message (index <code>%d</code>)
    </button>
</form>
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
    <button type="submit" name="index" value="0">
        The first message (index <code>0</code>)
    </button>
</form>
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
    <button type="submit" name="index" value="%d">
        The last message (index <code>%d</code>)
    </button>
</form>`, index, iMinus1, iMinus1, bot.data.Messages.Len()-1, bot.data.Messages.Len()-1)),
				})

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

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Get Chat Message",
				"desc": template.HTML(fmt.Sprintf(`
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
    <button type="submit" name="index" value="0">
        The first message (index <code>0</code>)
    </button>
</form>
<form action="get" method="post" enctype="application/x-www-form-urlencoded">
    <button type="submit" name="index" value="%d">
        The last message (index <code>%d</code>)
    </button>
</form>
<br><br>
%s<br><br>
%s<br><br>
%s
`, bot.data.Messages.Len()-1, bot.data.Messages.Len()-1, genMsgCardHTML("Previous", prev, &iMinus1), genMsgCardHTML("Target", msg, &index), genMsgCardHTML("Next", next, &iPlus1))),
			})
		})

		group.GET("/sessions/clean", func(ctx *gin.Context) {
			sessionClean()

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Session Management",
				"desc":  template.HTML("The sessions have been cleaned.<br>Sessions inactive for 15 minutes have been removed."),
			})
		})
		group.GET("/sessions/prune", func(ctx *gin.Context) {
			deleted := make([]string, chatSessions.Len())
			i := 0

			chatSessions.ForEach(func(id string, _ interface{}) bool {
				chatSessions.Delete(id)
				deleted[i] = id
				i++

				return true
			})

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Session Management",
				"desc":  template.HTML("The sessions have been pruned.<br>All sessions have been removed and discarded."),
			})
		})
		group.GET("/sessions/list", func(ctx *gin.Context) {
			l := make([]template.HTML, chatSessions.Len())
			i := 0

			chatSessions.ForEach(func(id string, rSession interface{}) bool {
				session := rSession.(*ChatSession)

				l[i] = template.HTML(fmt.Sprintf(`<a href="view?id=%s">%s</a> - <code>%s</code> - Created at %s, last modified %s`,
					url.QueryEscape(id), template.HTMLEscapeString(id),
					session.ClientIP,
					session.BeginTime.Format(generalTimeFmt),
					session.LastModified.Format(generalTimeFmt)))
				i++

				return true
			})

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Session List",
				"desc": template.HTML(fmt.Sprintf("Here is a list of all the currently active sessions.<br>There are <strong>%d</strong> sessions in total.",
					len(l))),
				"list": l,
			})
		})
		group.GET("/sessions/view", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "View Session",
				"desc": template.HTML(`
<form action="view" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Session ID:</h3>
	<input type="text" placeholder="Enter session ID..." name="id" required>
	<button type="submit">View</button>
</form>`),
			})
		})
		group.POST("/sessions/view", func(ctx *gin.Context) {
			session, ok := chatSessions.Get(ctx.PostForm("id")).(*ChatSession)
			if !ok {
				ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
					"title": "Session Viewer",
					"desc":  "That session doesn't exist!",
				})
				return
			}

			ctx.HTML(http.StatusOK, "messages.tmpl", gin.H{
				"title":    "Session Viewer",
				"desc":     "Here is a chat transcript of session " + session.ID + ".",
				"messages": session.History,
			})
		})

		group.GET("/pprof", func(ctx *gin.Context) {
			pprof.Index(ctx.Writer, ctx.Request)
		})
		group.GET("/pprof/block", pprofHandler(pprof.Index))
		group.GET("/pprof/heap", pprofHandler(pprof.Index))
		group.GET("/pprof/profile", pprofHandler(pprof.Profile))
		group.POST("/pprof/symbol", pprofHandler(pprof.Symbol))
		group.GET("/pprof/symbol", pprofHandler(pprof.Symbol))
		group.GET("/pprof/trace", pprofHandler(pprof.Trace))

		group.GET("/errors/panic", func(ctx *gin.Context) {
			err := ctx.DefaultQuery("err", "a panic for testing error/panic recovery systems")
			panic(err)
		})

		group.GET("/manage/template", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Evaluate Template",
				"desc": template.HTML(`
<textarea style="font-family:'Fira Code','Source Code Pro','Noto Mono','Ubuntu Mono',monospace;" form="tfid" placeholder="Enter template..." name="tmpl" required>
</textarea>
<form action="template" method="post" enctype="application/x-www-form-urlencoded" id="tfid">
	<button type="submit">Evaluate</button>
</form>`),
			})
		})
		group.POST("/manage/template", func(ctx *gin.Context) {
			tmpStr := ctx.PostForm("tmpl")
			tmpl, err := template.New("chatengine_admin_eval").Funcs(tmplFuncMap).Parse(tmpStr)
			if err != nil {
				panic(err)
			}

			tmplVars["ctx"] = ctx
			resultBuf := bytes.NewBufferString("")
			before := time.Now()
			err = tmpl.Execute(resultBuf, tmplVars)

			results := resultBuf.String()
			if err != nil {
				results += fmt.Sprintf(`

<<
ERROR evaluating template!
%s
>>`, err)
			}

			varGen := bytes.NewBufferString("")
			for k, v := range tmplVars {
				varGen.WriteString("<li><code>")
				varGen.WriteString(k)
				varGen.WriteString("</code>: <pre class=p><code>")
				varGen.WriteString(fmt.Sprintf("%s", v))
				varGen.WriteString("</code></pre>")
			}

			funcGen := bytes.NewBufferString("")
			for k, v := range tmplFuncMap {
				funcGen.WriteString("<li><code>")
				funcGen.WriteString(k)
				funcGen.WriteString("</code>: <pre class=p><code>")
				fstr := fmt.Sprintf("%s", v)
				funcGen.WriteString(fstr[4:strings.LastIndexByte(fstr, '=')])
				funcGen.WriteString("</code></pre>")
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Template Evaluation",
				"desc": template.HTML(fmt.Sprintf(`
<style>.p {
	display: inline-block;
	margin: 0 0;
}</style>
Took <strong>%s</strong> to execute.<br><br>
Input:
<pre><code>%s</code></pre>
Output:
<pre><code>%s</code></pre>
<br><br>
<p>
<details>
    <summary>Variables</summary>
    <p></p>
    <ul>
    %s
    </ul>
    <p></p>
</details>
</p>
<p>
<details>
    <summary>Functions</summary>
    <p></p>
    <ul>
    %s
    </ul>
    <p></p>
</details>
</p>`, time.Since(before), tmpStr, results, varGen.String(), funcGen.String())),
			})
		})
		group.GET("/manage/uptime", func(ctx *gin.Context) {
			systime, err := host.Uptime()
			if err != nil {
				panic(err)
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Server Uptime",
				"desc": fmt.Sprintf("This server has been up for %s. The system has been up for %s.",
					time.Since(startTime), time.Duration(systime*1000000000)),
			})
		})
		group.GET("/manage/free_mem", func(ctx *gin.Context) {
			debug.FreeOSMemory()

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Free OS Memory",
				"desc":  "The garbage collector has been called to free as much memory to the OS as possible.",
			})
		})
		group.GET("/manage/stack", func(ctx *gin.Context) {
			buf := make([]byte, 1<<16)
			stack := buf[:runtime.Stack(buf, true)]

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Full Program Stacks",
				"desc":  template.HTML(fmt.Sprintf(`<pre><code>%s</code></pre>`, util.BytesToString(stack))),
			})
		})
		group.GET("/manage/gc_stats", func(ctx *gin.Context) {
			gcStats := debug.GCStats{}
			debug.ReadGCStats(&gcStats)

			pauseDuraHTML := bytes.NewBufferString(`<ul>`)
			for _, duration := range gcStats.Pause {
				pauseDuraHTML.WriteString(`<li>`)
				pauseDuraHTML.WriteString(duration.String())
				pauseDuraHTML.WriteString(`</li>`)
			}
			pauseDuraHTML.WriteString(`</ul>`)

			pauseEndHTML := bytes.NewBufferString(`<ul>`)
			for _, t := range gcStats.PauseEnd {
				pauseEndHTML.WriteString(`<li>`)
				pauseEndHTML.WriteString(t.Format(preciseTimeFmt))
				pauseEndHTML.WriteString(`</li>`)
			}
			pauseEndHTML.WriteString(`</ul>`)

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title":    "Garbage Collector Statistics",
				"desc":     "All sublists are sorted most recent events first.",
				"listhead": "Values:",
				"list": []template.HTML{
					template.HTML("Last Collection: " + gcStats.LastGC.Format(preciseTimeFmt)),
					template.HTML(fmt.Sprintf("Number of Collections: %d", gcStats.NumGC)),
					template.HTML("Last Pause Durations: " + pauseDuraHTML.String()),
					template.HTML("Last Pause End Times: " + pauseEndHTML.String()),
					template.HTML("Total Pause Duration: " + gcStats.PauseTotal.String() + " (for all collections)"),
				},
			})
		})
		group.GET("/manage/gc_percent", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Garbage Collector Threshold",
				"desc": template.HTML(`
Set the minimum percentage of new data that must have been allocated compared to old data to trigger a garbage collection.<br><br>
As quoted from the Go documentation: SetGCPercent sets the garbage collection target percentage: a collection is triggered when the ratio of freshly allocated data to live data remaining after the previous collection reaches this percentage. SetGCPercent returns the previous setting. The initial setting is the value of the <code>GOGC</code> environment variable at startup, or 100 if the variable is not set. A negative percentage disables garbage collection.<br><br>
<form action="gc_percent" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Percentage:</h3>
	<input type="number" placeholder="Enter threshold percentage..." name="percent" required>
	<button type="submit">Set Threshold</button>
</form>`),
			})
		})
		group.POST("/manage/gc_percent", func(ctx *gin.Context) {
			arg := ctx.PostForm("percent")
			percent, err := strconv.Atoi(arg)
			if err != nil {
				ctx.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
					"admin": true,
					"desc": template.HTML(`Invalid percentage supplied!<br>
<a href="gc_percent">Back</a>`),
					"err": err,
				})
				return
			}

			old := debug.SetGCPercent(percent)
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Garbage Collector Threshold",
				"desc": template.HTML(fmt.Sprintf(`The garbage collector threshold is now set to <strong>%d%%</strong>.<br>
The old threshold was %d%%.`, percent, old)),
			})
		})
		group.GET("/manage/max_threads", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Maximum Thread Count",
				"desc": template.HTML(`
Set the maximum amount of OS threads the server may create.<br><br>
As quoted from the Go documentation: SetMaxThreads sets the maximum number of operating system threads that the Go program can use. If it attempts to use more than this many, the program crashes. SetMaxThreads returns the previous setting. The initial setting is 10,000 threads.<br><br>
The limit controls the number of operating system threads, not the number of goroutines. A Go program creates a new thread only when a goroutine is ready to run but all the existing threads are blocked in system calls, cgo calls, or are locked to other goroutines due to use of <code>runtime.LockOSThread</code>.<br><br>
SetMaxThreads is useful mainly for limiting the damage done by programs that create an unbounded number of threads. The idea is to take down the program before it takes down the operating system.<br><br>
<form action="max_threads" method="post" enctype="application/x-www-form-urlencoded">
	<h3>Thread Count:</h3>
	<input type="number" placeholder="Enter thread count..." name="max" required>
	<button type="submit">Set Maximum</button>
</form>`),
			})
		})
		group.POST("/manage/max_threads", func(ctx *gin.Context) {
			arg := ctx.PostForm("max")
			max, err := strconv.Atoi(arg)
			if err != nil {
				ctx.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
					"admin": true,
					"desc": template.HTML(`Invalid count supplied!<br>
<a href="max_threads">Back</a>`),
					"err": err,
				})
				return
			}

			old := debug.SetMaxThreads(max)
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Maximum Thread Count",
				"desc": template.HTML(fmt.Sprintf(`The maximum thread count is now <strong>%d</strong>.<br>
The old maximum was %d.`, max, old)),
			})
		})
		group.GET("/manage/heap_dump", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Heap Dump",
				"desc": template.HTML(`
Write a full heap dump to a file.<br><br>
<form action="max_threads" method="post" enctype="application/x-www-form-urlencoded">
	<h3>File Name/Path:</h3>
	<input type="text" placeholder="Enter filename/path..." name="max" required>
	<button type="submit"><strong>Write Heap Dump</strong></button>
</form>
<h2>Or...</h2>
Send a full heap dump.<br>
<a href="get_heap_dump">Download Heap Dump</a>`),
			})
		})
		group.POST("/manage/heap_dump", func(ctx *gin.Context) {
			path := ctx.PostForm("path")
			f, err := os.Create(path)
			if err != nil {
				ctx.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
					"admin": true,
					"desc": template.HTML(`Invalid path supplied!<br>
<a href="heap_dump">Back</a>`),
					"err": err,
				})
				return
			}

			debug.WriteHeapDump(f.Fd())
			f.Close()

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Heap Dump",
				"desc":  template.HTML(fmt.Sprintf("The heap has been dumped to <code>%s</code>.", path)),
			})
		})
		group.GET("/manage/get_heap_dump", func(ctx *gin.Context) {
			if runtime.GOOS == "linux" {
				w := ctx.Writer

				ptrVal := reflect.ValueOf(w)
				val := reflect.Indirect(ptrVal)

				// w is a "http.response" struct from which we get the 'conn' field
				valconn := val.FieldByName("conn")
				val1 := reflect.Indirect(valconn)

				// which is a http.conn from which we get the 'rwc' field
				ptrRwc := val1.FieldByName("rwc").Elem()
				rwc := reflect.Indirect(ptrRwc)

				// which is net.TCPConn from which we get the embedded conn
				val1conn := rwc.FieldByName("conn")
				val2 := reflect.Indirect(val1conn)

				// which is a net.conn from which we get the 'fd' field
				fdmember := val2.FieldByName("fd")
				val3 := reflect.Indirect(fdmember)

				// which is a netFD from which we get the 'sysfd' field
				netFdPtr := val3.FieldByName("sysfd")

				fd := int(netFdPtr.Int())
				ctx.Writer.Header()["Content-Type"] = []string{"application/octet-stream"}
				debug.WriteHeapDump(uintptr(fd))
			} else {
				ctx.String(http.StatusServiceUnavailable, "This only works on Linux")
			}
		})
		group.GET("/manage/get_heap_dump__nconfirm", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "Get Heap Dump?",
				"desc": template.HTML(`<h2>Are you sure you want to download a <strong>full</strong> heap dump?</h2><br>
<a href="get_heap_dump">Yes, I'm sure.</a><br><a style="font-size: 20pt;" href="../">No, I'm not sure.</a>`),
			})
		})
		group.GET("/manage/sysinfo", func(ctx *gin.Context) {
			finalHTML := bytes.NewBuffer(make([]byte, 0))
			finalHTML.Grow(8192)

			sections := make([]infoSection, 8)

			{
				gcStats := debug.GCStats{}
				debug.ReadGCStats(&gcStats)

				memStats := runtime.MemStats{}
				runtime.ReadMemStats(&memStats)

				pauseDuraHTML := bytes.NewBufferString(`<ul>`)
				for _, duration := range gcStats.Pause {
					pauseDuraHTML.WriteString(`<li>`)
					pauseDuraHTML.WriteString(duration.String())
					pauseDuraHTML.WriteString(`</li>`)
				}
				pauseDuraHTML.WriteString(`</ul>`)

				pauseEndHTML := bytes.NewBufferString(`<ul>`)
				for _, t := range gcStats.PauseEnd {
					pauseEndHTML.WriteString(`<li>`)
					pauseEndHTML.WriteString(t.Format(preciseTimeFmt))
					pauseEndHTML.WriteString(`</li>`)
				}
				pauseEndHTML.WriteString(`</ul>`)

				sections[0] = infoSection{
					title: "Go Runtime",
					content: genMap{
						"Version": runtime.Version(),
						"GOOS": runtime.GOOS,
						"GOARCH": runtime.GOARCH,
						"GOROOT": codeSect(runtime.GOROOT()),
						"Goroutines": runtime.NumGoroutine(),
						"Usable CPUs": runtime.NumCPU(),
						"Cgo Calls": runtime.NumCgoCall(),
						"Garbage Collector": genMap{
							"Last Collection": gcStats.LastGC.Format(preciseTimeFmt),
							"Number of Collections": gcStats.NumGC,
							"Last Pause Durations": pauseDuraHTML.String(),
							"Last Pause End Times": pauseEndHTML.String(),
							"Total Pause Duration": gcStats.PauseTotal.String() + " (for all collections)",
							"Time Used": fmt.Sprintf("%.1f%% of program's total time", memStats.GCCPUFraction * 100),
							"Metadata Size": byteNum(memStats.GCSys),
							"Forced Collections": memStats.NumForcedGC,
							"Next Target Heap": byteNum(memStats.NextGC),
						},
						"Current Stack": codeSect(string(debug.Stack())),
						"Actions": []string{
							`<a href="free_mem">Free OS Memory</a>`,
							`<a href="heap_dump">Set Heap Dump</a>`,
							`<a href="max_threads">Set Maximum Thread Count</a>`,
							`<a href="gc_percent">Set Garbage Collection Threshold</a>`,
						},
						"Heap": genMap{
							"Used for Objects": byteNum(memStats.HeapAlloc),
							"Idle": byteNum(memStats.HeapIdle),
							"In Use": byteNum(memStats.HeapInuse),
							"Objects": memStats.HeapObjects,
							"Allocated": byteNum(memStats.HeapSys),
							"Released": byteNum(memStats.HeapReleased),
							"Object Frees": memStats.Frees,
							"Object Allocations": memStats.Mallocs,
							"Total Allocated": byteNum(memStats.TotalAlloc),
						},
						"Internal Structures": genMap{
							"<code>mcache</code> Allocated": byteNum(memStats.MCacheSys),
							"<code>mcache</code> In Use": byteNum(memStats.MCacheInuse),
							"<code>mspan</code> Allocated": byteNum(memStats.MSpanSys),
							"<code>mspan</code> In Use": byteNum(memStats.MSpanInuse),
						},
						"Stack": genMap{
							"In Use": byteNum(memStats.StackInuse),
							"Allocated": byteNum(memStats.StackSys),
						},
						"Pointer Lookups": memStats.Lookups,
						"Total Allocated Memory": byteNum(memStats.Sys),
						"Miscellaneous Allocated": byteNum(memStats.OtherSys),
					},
				}
			}

			{
				hif, _ := host.Info()
				usrs, _ := host.Users()
				sens, _ := host.SensorsTemperatures()

				usrMap := make(map[string]string, len(usrs))
				senMap := make(map[string]string, len(sens))

				for _, usr := range usrs {
					usrMap[usr.User] = fmt.Sprintf("On %s@%s \u2022 Started %d", usr.Terminal, usr.Host, usr.Started)
				}
				for _, sen := range sens {
					senMap[sen.SensorKey] = fmt.Sprintf("%.8f", sen.Temperature)
				}

				sections[1] = infoSection{
					title: "Host",
					content: genMap{
						"Uptime": time.Duration(hif.Uptime * 1e9),
						"Boot Time": time.Unix(int64(hif.BootTime), 0),
						"Kernel Version": hif.KernelVersion,
						"Platform": hif.Platform,
						"Family": hif.PlatformFamily,
						"Version": hif.PlatformVersion,
						"Virtualization Role": hif.VirtualizationRole,
						"Virtualization System": hif.VirtualizationSystem,
						"Hostname": hif.Hostname,
						"Host ID": hif.HostID,
						"OS": hif.OS,
						"Processes": hif.Procs,
						"Users": usrMap,
						"Sensors": sens,
					},
				}
			}

			{
				ct, _ := cpu.Counts(false)
				lct, _ := cpu.Counts(true)
				infos, _ := cpu.Info()
				times, _ := cpu.Times(false)
				first := infos[0]
				t := times[0]

				flags := bytes.NewBufferString(`<ul>`)
				for _, fl := range first.Flags {
					flags.WriteString(`<li><code>`)
					flags.WriteString(fl)
					flags.WriteString(`</code></li>`)
				}
				flags.WriteString(`</ul>`)

				sections[2] = infoSection{
					title: "CPU",
					content: genMap{
						"Count": ct,
						"Logical Count": lct,
						"CPU 0 Information": genMap{
							"CPU": first.CPU,
							"Cores": first.Cores,
							"Core ID": first.CoreID,
							"Family": first.Family,
							"Model": first.Model,
							"Model Name": first.ModelName,
							"Stepping": first.Stepping,
							"Vendor": first.VendorID,
							"Clock Speed": fmt.Sprintf("%.4f MHz (%.4f GHz)", first.Mhz, first.Mhz / 1000.0),
							"Microcode": first.Microcode,
							"Cache Size": first.CacheSize,
							"Physical ID": first.PhysicalID,
							"Flags": flags.String(),
						},
						"Total Time": fmt.Sprintf("%.4f seconds", t.Total()),
						"Times": genMap{
							"CPU": t.CPU,
							"User": t.User,
							"System": t.System,
							"Idle": t.Idle,
							"IRQ": t.Irq,
							"SoftIRQ": t.Softirq,
							"Nice": t.Nice,
							"iowait": t.Iowait,
							"Steal": t.Steal,
							"Stolen": t.Stolen,
							"Guest": t.Guest,
							"Guest Nice": t.GuestNice,
						},
					},
				}
			}

			{
				parts, _ := disk.Partitions(true)
				partitions := make([]genMap, len(parts))

				for i, part := range parts {
					partitions[i] = genMap{
						"Device": codeSect(part.Device),
						"Filesystem": part.Fstype,
						"Mount Point": part.Mountpoint,
						"Options": codeSect(part.Opts),
					}
				}

				sections[3] = infoSection{
					title: "Disk",
					content: genMap{
						"Partitions": partitions,
					},
				}
			}

			{
				avg, _ := load.Avg()
				misc, _ := load.Misc()

				sections[4] = infoSection{
					title: "Load",
					content: genMap{
						"Load Average (last minute)": avg.Load1,
						"Load Average (last 5 minutes)": avg.Load5,
						"Load Average (last 15 minutes)": avg.Load15,
						"Processes Running": misc.ProcsRunning,
						"Processes Blocked": misc.ProcsBlocked,
						"Context": misc.Ctxt,
					},
				}
			}

			{
				virt, _ := mem.VirtualMemory()
				swap, _ := mem.SwapMemory()

				sections[5] = infoSection{
					title: "Memory",
					content: genMap{
						"Total": byteNum(virt.Total),
						"Free": byteNum(virt.Free),
						"Used": fmt.Sprintf("%.1f%% (%s)", virt.UsedPercent, byteNum(virt.Used)),
						"Active": byteNum(virt.Active),
						"Inactive": byteNum(virt.Inactive),
						"Available": byteNum(virt.Available),
						"Wired": byteNum(virt.Wired),
						"Cached": byteNum(virt.Cached),
						"Shared": byteNum(virt.Shared),
						"Dirty": byteNum(virt.Dirty),
						"Buffers": byteNum(virt.Buffers),
						"Page Tables": byteNum(virt.PageTables),
						"Slab": byteNum(virt.Slab),
						"Swap Cached": byteNum(virt.SwapCached),
						"Writeback": byteNum(virt.Writeback),
						"Swap": genMap{
							"Total": byteNum(swap.Total),
							"Free": byteNum(swap.Free),
							"Swap In": byteNum(swap.Sin),
							"Swap Out": byteNum(swap.Sout),
							"Used": fmt.Sprintf("%.1f%% (%s)", swap.UsedPercent, byteNum(swap.Used)),
						},
					},
				}
			}

			{
				ifaces, _ := net.Interfaces()
				ifMap := make([]genMap, len(ifaces))

				for idx, iface := range ifaces {
					addrs := make([]codeSect, len(iface.Addrs))
					for idx, addr := range iface.Addrs {
						addrs[idx] = codeSect(addr.Addr)
					}

					ifMap[idx] = genMap{
						"Name": codeSect(iface.Name),
						"Hardware Address": codeSect(iface.HardwareAddr),
						"MTU": iface.MTU,
						"Addresses": addrs,
						"Flags": iface.Flags,
					}
				}

				sections[6] = infoSection{
					title: "Network",
					content: genMap{
						"Interfaces": ifMap,
					},
				}
			}

			{
				pids, _ := process.Pids()
				procs := make([]codeSect, len(pids))

				for idx, pid := range pids {
					procs[idx] = codeSect(fmt.Sprintf("%d", pid))
				}

				sections[7] = infoSection{
					title: "Process",
					content: genMap{
						"PIDs": procs,
					},
				}
			}

			for _, section := range sections {
				finalHTML.WriteString(`<details><summary>`)
				finalHTML.WriteString(section.title)
				finalHTML.WriteString(`</summary><p></p><ul>`)
				for k, v := range section.content {
					finalHTML.WriteString(`<li>`)
					finalHTML.WriteString(k)
					finalHTML.WriteString(": ")

					renderType(finalHTML, v)

					finalHTML.WriteString(`</li>`)
				}
				finalHTML.WriteString(`</ul><p></p></details>`)
			}

			ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
				"title": "System Information",
				"desc":  template.HTML(strings.Replace(finalHTML.String(), "0.000", "0", -1)),
			})
		})
	}
}

type genMap = map[string]interface{}
type codeSect string

type infoSection struct {
	title string
	content genMap
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
