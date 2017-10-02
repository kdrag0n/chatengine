package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"chatengine/util"
	"github.com/OneOfOne/cmap/stringcmap"
	"github.com/bwmarrin/snowflake"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
	"runtime/debug"
	"flag"
)

var (
	validAPIKeys        map[string]struct{}
	validAdminKeys      map[string]struct{}
	blockedIPs          map[string]struct{}
	selflearnBlockedIPs map[string]struct{}
	logicAdapters       = []logicAdapter{
		DistanceAdapter,
		RandomConfidenceAdapter,
	}
	chatSessions  = stringcmap.NewSize(64) // TODO: genx dedicated cmap KT=string,VT=*ChatSession
	initFuncs     = [...]func(){}
	userGreetings = [...]string{
		"hi",
		"hey",
		"hello",
		"heya",
		"heyo",
		"greetings",
		"hi there",
		"yo",
		"g'day",
		"good day",
		"morning",
		"good morning",
		"good night",
		"gnight",
		"g'night",
		"mornin",
		"hiya",
		"hallo",
	}
	userGreetingRegexp = func() *regexp.Regexp {
		regex := bytes.NewBufferString(`(?i)^\b?(?:`)

		for idx, greeting := range userGreetings {
			regex.WriteString(strings.Replace(greeting, " ", `\s+`, -1))

			if idx != len(userGreetings)-1 {
				regex.WriteByte('|')
			}
		}

		regex.WriteString(`)\b`)

		return regexp.MustCompile(regex.String())
	}()
	greetings = [...]string{
		"Hey!",
		"Oh hey there!",
		"Heya!",
		"Hmm. Fluffy clouds. AAH hi, I didn't see you there.",
		"Yeah? Oh, hi there.",
		"Hi, I guess.",
		"Don't disturb me...",
		"I'm kind of busy right now, but hi.",
		"Hi. Bye.",
		"Greetings, partner.",
		"Isn't it a beautiful day?",
		"Mmmm. Delicious. Hey, I'm busy eating right now.",
		"Hello!",
		"Hello, friend.",
		"Hello, random person.",
		"Hello, human.",
	}
	logger = func() *zap.SugaredLogger {
		logger, _ := zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build()

		return logger.Sugar()
	}()
	webchatURL        string
	config            Config
	selflearn         = true
	autowriteChatData = true
	snowNode          = func() (node *snowflake.Node) {
		node, err := snowflake.NewNode(4)
		if err != nil {
			panic(err)
		}

		return
	}()
	startTime = time.Now()
	isDebug = false
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "activate debug mode (log requests, reload templates, etc)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	go func() {
		for range time.Tick(time.Hour) {
			rand.Seed(time.Now().UnixNano())
		}
	}()
}

type logicAdapter func(*QueryCtx) PossibleResponse

// DataRoot represents the root of all the chat data.
type DataRoot struct {
	Messages *list.List
	Version  uint8
	Flags    uint32
}

// RawDataRoot represents the raw root of read chat data.
type RawDataRoot struct {
	Messages []*ChatMessage `bson:"messages"`
	Version  uint8          `bson:"version"`
	Flags    uint32         `bson:"flags"`
}

// Config represents the server configuration.
type Config struct {
	ListenOn            string   `json:"listen"`
	API                 []string `json:"api"`
	Admin               []string `json:"admin"`
	BlockedIPs          []string `json:"blocked_ips"`
	SelflearnBlockedIPs []string `json:"selflearn_blocked_ips"`
	RootURL             string   `json:"root_url"`
	Environment         string   `json:"env"`
	RecaptchaSiteKey    string   `json:"recaptcha_site_key"`
	RecaptchaSecretKey  string   `json:"recaptcha_secret_key"`
	DSN                 string   `json:"dsn"`
	RedisProtocol       string   `json:"redis_protocol"`
	RedisServer         string   `json:"redis_server"`
	RedisPassword       string   `json:"redis_password"`
	RedisDB             int      `json:"redis_database"`
}

// ChatMessage represents a single response.
type ChatMessage struct { // TODO: move a lot of this chat stuff to chatengine/chat or something, also to generate cmaps
	Text         string `bson:"text"`
	textRunes    []rune
	IsCJK        bool           `bson:"is_cjk"`
	InResponseTo []*ResponseFor `bson:"in_response_to"`
	CreatedAt    time.Time      `bson:"created_at"`
	ExtraData    []ExtraData    `bson:"extra_data"`
	Occurrences  int            `bson:"occurrence"`
}

// ResponseFor represents a message which would get the parent response.
type ResponseFor struct {
	Text        string `bson:"text"`
	textRunes   []rune
	Occurrences int `bson:"occurrence"`
}

// ExtraData represents extra data for a message.
type ExtraData struct {
	Name  string      `bson:"name"`
	Value interface{} `bson:"value"`
}

// ChatBot is the main chatbot object, where logic takes place.
type ChatBot struct {
	data  DataRoot
	db    *gorm.DB
	rPool *redis.Pool
}

// PossibleResponse represents a possible response to return, with a set confidence level.
type PossibleResponse struct {
	response       string
	useGetter      bool
	responseGetter func() string
	confidence     float64
	chatMessage    *ChatMessage
	userBestMatch  *ChatMessage
}

// ChatSession represents an active chat session.
type ChatSession struct {
	ID           string
	History      []HistoryEntry
	BeginTime    time.Time
	LastModified time.Time
	ClientIP     string
}

// HistoryEntry represents a single message entry in the history of a chat session.
type HistoryEntry struct {
	Message   *ChatMessage
	BestMatch *ChatMessage
	CreatedAt time.Time
}

// QueryCtx represents the full context of a query request, sent over HTTP.
type QueryCtx struct {
	Query       string
	QueryRunes  []rune
	origQuery   string
	isCJK       bool
	queryLength int
	bot         *ChatBot
	request     *gin.Context
	session     *ChatSession
}

func sessionClean() {
	chatSessions.ForEach(func(id string, sessionInterface interface{}) bool {
		session := sessionInterface.(*ChatSession)

		if time.Since(session.LastModified).Minutes() > 14.85 {
			chatSessions.Delete(id)
		}

		return true
	})
}

func sessionReaper() {
	for range time.Tick(time.Minute * 5) {
		sessionClean()
	}
}

func chatDataWriter(bot *ChatBot) {
	for range time.Tick(time.Minute * 5) {
		if autowriteChatData {
			err := bot.WriteChatData()

			if err != nil {
				logger.Errorf("Error writing chat data: %s", err)
			}
		}
	}
}

func (message *ChatMessage) String() string {
	return message.Text
}

// LoadChatData loads chat data from disk, discarding the in-memory copy.
func (bot *ChatBot) LoadChatData() error {
	isBson := false

	fileReader, err := os.Open("data.gob")
	if err != nil {
		if os.IsNotExist(err) {
			fileReader, err = os.Open("chat.bson")
			if err != nil {
				return err
			}
			defer fileReader.Close()

			isBson = true
		} else {
			return err
		}
	}

	reader := bufio.NewReader(fileReader)

	var rawData RawDataRoot
	if isBson {
		logger.Info("Detected old chat.bson datastore. Converting to data.gob...")
		var bytes []byte
		bytes, err = ioutil.ReadAll(reader)
		if err != nil {
			return err
		}

		err = bson.Unmarshal(bytes, &rawData)
		if err != nil {
			return err
		}
	} else {
		err = gob.NewDecoder(reader).Decode(&rawData)
		if err != nil {
			return err
		}
	}

	newMsgList := list.New()
	for _, msg := range rawData.Messages {
		newMsgList.PushBack(msg)
	}

	if rawData.Version < 1 {
		logger.Info("Detected older data version. Upgrading...")
		for e := newMsgList.Front(); e != nil; e = e.Next() {
			msg := e.Value.(*ChatMessage)
			msg.IsCJK = util.ContainsCJK(msg.Text)
		}

		rawData.Version = 1
		logger.Info("Data format upgraded.")
	}

	logger.Info("Performing rune transformation...")
	for e := newMsgList.Front(); e != nil; e = e.Next() {
		msg := e.Value.(*ChatMessage)
		msg.textRunes = []rune(msg.Text)

		for _, respFor := range msg.InResponseTo {
			respFor.textRunes = []rune(respFor.Text)
		}
	}
	logger.Info("Full rune transformation finished.")

	bot.data = DataRoot{
		Messages: newMsgList,
		Version:  rawData.Version,
		Flags:    rawData.Flags,
	}

	if isBson {
		err = bot.WriteChatData()
		if err != nil {
			return err
		}

		logger.Info("Converted to data.gob.")
	}

	return nil
}

// WriteChatData writes the current in-memory chat data to disk.
func (bot *ChatBot) WriteChatData() error {
	msgList := *bot.data.Messages
	msgSlice := make([]*ChatMessage, msgList.Len())

	for i, e := 0, msgList.Front(); e != nil; i, e = i+1, e.Next() {
		msgSlice[i] = e.Value.(*ChatMessage)
	}

	atomNum := rand.Intn(1024)
	atomFilename := fmt.Sprintf("data.gob.atom%d", atomNum)
	writer, err := os.Create(atomFilename)
	if err != nil {
		return err
	}
	defer writer.Close()

	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(RawDataRoot{
		Messages: msgSlice,
		Version:  bot.data.Version,
		Flags:    bot.data.Flags,
	})
	if err != nil {
		return err
	}
	writer.Close()

	err = os.Rename(atomFilename, "data.gob")
	return err
}

func loadConfig() (*Config, error) {
	logger.Info("Loading config...")
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Errorf("Error loading config: %s", err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		logger.Errorf("Error loading config: %s", err)
		return nil, err
	}

	webchatURL = config.RootURL + "chat"

	empStruct := struct{}{}
	validAPIKeys = make(map[string]struct{}, len(config.API))
	for _, key := range config.API {
		validAPIKeys[key] = empStruct
	}

	validAdminKeys = make(map[string]struct{}, len(config.Admin))
	for _, key := range config.Admin {
		validAdminKeys[key] = empStruct
	}

	blockedIPs = make(map[string]struct{}, len(config.BlockedIPs))
	for _, ip := range config.BlockedIPs {
		blockedIPs[ip] = empStruct
	}

	selflearnBlockedIPs = make(map[string]struct{}, len(config.SelflearnBlockedIPs))
	for _, ip := range config.SelflearnBlockedIPs {
		selflearnBlockedIPs[ip] = empStruct
	}

	logger.Info("Config loaded.")
	return &config, nil
}

func main() {
	debug.SetMaxStack(384000000)

	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	logger.Info("Connecting to MySQL database...")
	logger.Sync()
	db, err := gorm.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.DB().SetConnMaxLifetime(time.Minute * 15)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(32)

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	logger.Info("Connecting to Redis database...")
	logger.Sync()
	pool := &redis.Pool{
		MaxIdle:     48,
		MaxActive:   32,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err2 := redis.Dial(config.RedisProtocol, config.RedisServer,
				redis.DialDatabase(config.RedisDB), redis.DialPassword(config.RedisPassword))
			if err2 != nil {
				return nil, err2
			}

			return c, nil
		},
	}
	defer pool.Close()

	_, err = pool.Get().Do("PING")
	if err != nil {
		panic(err)
	}

	bot := &ChatBot{
		data:  DataRoot{},
		db:    db,
		rPool: pool,
	}

	logger.Info("Loading chat data...")
	logger.Sync()
	err = bot.LoadChatData()
	if err != nil {
		panic(err)
	}

	for _, ifunc := range initFuncs {
		ifunc()
	}

	logger.Info("Setting up HTTP server...")
	logger.Sync()
	gin.SetMode(gin.ReleaseMode)
	engine := createHTTPHandler(bot)

	logger.Info("Starting background tasks...")
	logger.Sync()
	go sessionReaper()
	go chatDataWriter(bot)
	go distanceCacheJanitor()
	go util.FormatCacheJanitor()
	go selflearnJobDispatcher()
	go askTimeLogger()
	defer bot.WriteChatData()
	go func() {
		server := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSConfig:    tlsConfig,
			Addr:         config.ListenOn,
			Handler:      engine,
		}

		eerr := server.ListenAndServe()
		if eerr != nil {
			panic(eerr)
		}
	}()
	logger.Infof("HTTP server listening on %s", config.ListenOn)
	<-make(chan struct{})

	logger.Info("Shutting down.")
}
