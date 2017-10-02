package main

import (
	crand "crypto/rand"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	ttemplate "text/template"
	"time"
	"unicode/utf8"

	"chatengine/sessions"
	"chatengine/util"
	"crypto/sha512"
	"encoding/base64"
	"github.com/OneOfOne/cmap/stringcmap"
	"github.com/aviddiviner/gin-limit"
	"github.com/cznic/strutil"
	"github.com/didip/tollbooth"
	tb_config "github.com/didip/tollbooth/config"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"github.com/rs/xid"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})

	tlsConfig = &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		// compatibility loss:
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	limiters = map[byte]*tb_config.Limiter{
		'R': newLimiter(2, time.Millisecond*1500, "POST"),
		'W': newLimiter(4, time.Millisecond*1500, "POST"),
		'_': newLimiter(6, time.Second, "POST"),
		'X': newLimiter(8, time.Second, "POST"),
		'Z': newLimiter(10, time.Second, "POST"),
		'U': newLimiter(1000000, time.Second, "POST"),
	}

	funcMap = map[string]interface{}{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"isEven": func(i int) bool {
			return i%2 == 0
		},
		"stringNeq": func(a, b string) bool {
			return a != b
		},
		"relativeRoot": relativeRoot,
		"finalRelRoot": func(rr *string, ctx *gin.Context) string {
			if rr != nil && *rr != "" {
				return *rr
			}

			return relativeRoot(ctx)
		},
		"calculatePath": func(ctx *gin.Context, mode, path string) string {
			switch mode {
			case "absolute":
				return config.RootURL + path
			case "relative":
				return ctx.Request.URL.Path[:strings.LastIndexByte(ctx.Request.URL.Path, '/')] + path
			case "full":
				return path
			case "relroot":
				return relativeRoot(ctx) + path
			default:
				panic("invalid path mode supplied for path resolving")
			}
		},
		"head": func(title, desc, path string) template.HTML {
			return template.HTML(fmt.Sprintf(`
<meta charset="utf-8">
<meta name="description" content="%s">
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes, maximum-scale=2">
<meta name="HandheldFriendly" content="true">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
<meta name="format-detection" content="telephone=no">
<meta name="apple-mobile-web-app-title" content="ChatEngine">
<meta name="keywords" content="web,chat,nlp,language,natural,conversation,chatbot,cleverbot,cool">
<meta name="author" content="Armored Dragon">
<title>ChatEngine - %s</title>
<meta property="og:locale" content="en_US">
<meta property="og:site_name" content="ChatEngine">
<meta property="og:title" content="ChatEngine: %s">
<meta property="og:url" content="https://chatengine.xyz/%s">
<meta property="og:type" content="website">
<meta property="og:description" content="%s">
<meta property="og:image" content="https://chatengine.xyz/static/img/icon.png">
<link rel="apple-touch-startup-page" href="/static/img/bg.png">
<link rel="shortcut icon" type="image/png" href="/static/img/icon.png">
<meta name="twitter:title" content="ChatEngine: %s">
<meta name="twitter:url" content="https://chatengine.xyz/%s">
<meta name="twitter:description" content="%s">
<meta name="twitter:image" content="https://chatengine.xyz/static/img/icon.png">
<meta name="twitter:card" content="summary">`, desc, title, title, path, desc, title, path, desc))
		},
		"knockout": func(text string) template.HTML {
			return template.HTML(fmt.Sprintf(`<svg>
	<defs>
		<g id="text">
			<text text-anchor="middle" x="0" y="0" dy="1">%s</text>
		</g>
		<mask id="mask" x="0" y="0" width="100" height="50">
			<rect x="0" y="0" width="50" height="40" fill="#fff"></rect>
			<use xlink:href="#text"></use>
		</mask>
	</defs>
	<use xlink:href="#text" mask="url(#mask)"></use>
</svg>`,text))
		},
	}

	wcTokenMethods   = []string{"POST", "PATCH", "PUT"}
	wcGlobalKey      string
	wcRvqPath        string
	wcTokenSep       string
	wcOldTokenSep    = ""
	wcTokenJsErr     = []byte(`createMessage('them', 'An error occurred communicating with the server.\nTry reloading to fix the problem.');`)
	wcTokenSepChars  = []byte(" !\"#$%&()*+,-.0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~")
	wcTokenMap       = stringcmap.NewSize(64) // TODO: genx dedicated cmap KT=string,VT=*WebchatToken
	wcTokenLimiter   = newLimiters(12, time.Minute, wcTokenMethods)
	wcRand8Garbage   string
	wcFeatureHeaders = []string{
		"X-Mppr3cfda9b",
		"X-2jkg29ejfio",
		"X-D1iu2ocmi1291",
		"X-Djiagf9uc2389cu2",
		"X-Oi010auy3256ctb478vy53n847cmh2",
		"X-Fidsiifsjviinaniqo10309578j8ok",
		"X-Vvvdkwjeqixeiqomueio12x2r",
		"X-Pp001987n78mh912ham9",
		"X-92018p91ka90125yn5b",
		"X-123912c48923i4a1290y4a1289c4cy3hm9",
	}
	wcTempKey = genWcGeneralKey() // TODO: rotate this every few minutes or so

	tmplEngine *ttemplate.Template
	cgrMax     = big.NewInt(math.MaxInt64)
)

// HTTPResponse represents a response to the client.
type HTTPResponse struct {
	context    *gin.Context
	status     int
	Success    bool    `json:"success"`
	Session    string  `json:"session,omitempty"`
	Error      string  `json:"error,omitempty"`
	Response   string  `json:"response"`
	Confidence float64 `json:"confidence"`
	Filled     bool    `json:"-"`
}

func (resp *HTTPResponse) err(status int, err string) {
	resp.status = status
	resp.Success = false
	resp.Error = err
	resp.Response = "An error occurred: " + err
	resp.Filled = true
}

func (resp *HTTPResponse) respond(message string, confidence float64) {
	resp.status = http.StatusOK
	resp.Success = true
	resp.Error = ""
	resp.Response = message
	resp.Confidence = confidence
	resp.Filled = true
}

func newResponse(ctx *gin.Context) HTTPResponse {
	return HTTPResponse{
		context:  ctx,
		status:   http.StatusInternalServerError,
		Success:  true,
		Error:    "",
		Response: "",
		Session:  "",
		Filled:   false,
	}
}

func relativeRoot(ctx *gin.Context) string {
	relRoot := strings.Repeat("../", strings.Count(ctx.Request.URL.Path, "/")-1)
	if relRoot == "" {
		relRoot = "./"
	}

	return relRoot
}

type webchatToken struct {
	token     string
	ip        string
	random    string
	host      string
	href      string
	userAgent string
	path      string
	features  string
	createdAt time.Time
}

// RequestData represents the data a client sends with a request.
type RequestData struct {
	Session           string `json:"session"`
	Query             string `json:"query"`
	ModernizrFeatures string `json:"f"`
}

// HandleHTTP handles incoming HTTP requests at /ask for the chatbot.
func (bot *ChatBot) HandleHTTP(ctx *gin.Context) {
	response := newResponse(ctx)
	clientIP := ctx.ClientIP()
	defer func() {
		if !response.Filled {
			return
		}

		ctx.JSON(response.status, response)
	}()

	if _, ok := blockedIPs[clientIP]; ok {
		response.err(http.StatusForbidden, "You are not permitted to use this service.")
		return
	}

	providedKey := ctx.GetHeader("Authorization")
	var wcToken *webchatToken
	isWebchat := false

	if len(providedKey) < 24 {
		response.err(http.StatusUnauthorized, "You must supply a valid API key!")
		return
	}

	if _, ok := validAPIKeys[providedKey]; !ok {
		if k, ok := wcTokenMap.GetOK(providedKey); ok {
			wcToken = k.(*webchatToken)
			isWebchat = true
		} else {
			response.err(http.StatusUnauthorized, "You must supply a valid API key!")
			return
		}
	}

	referrer := ctx.GetHeader("Referer")
	userAgent := ctx.GetHeader("User-Agent")
	headerIdx := int(math.Floor(float64(util.CurrentTimeMillis())/300000)) % 10
	println(ctx.GetHeader(wcFeatureHeaders[headerIdx]) != wcToken.features)
	// STOPSHIP // TODO: remove that print

	if isWebchat && (referrer != webchatURL || referrer != wcToken.href || ctx.Request.Host != wcToken.host ||
		clientIP != wcToken.ip || userAgent != wcToken.userAgent ||
		ctx.GetHeader(wcFeatureHeaders[headerIdx]) != wcToken.features ||
		!strings.HasPrefix(userAgent, "Mozilla/5.0")) {
		response.err(http.StatusUnauthorized, "You must supply a valid API key!\u200b")
		return
	}

	limiter, ok := limiters[providedKey[0]]
	if !ok {
		if isWebchat {
			limiter = limiters['W']
		} else {
			limiter = limiters['_']
		}
	}

	httpError := tollbooth.LimitByRequest(limiter, ctx.Request)
	if httpError != nil {
		response.err(httpError.StatusCode, "You have reached the ratelimit for your key!")
		return
	}

	var data RequestData
	err := ctx.BindJSON(&data)
	if err != nil {
		response.err(http.StatusBadRequest, "Failed to decode JSON! "+err.Error())
		return
	}

	var sessionID string
	if data.Session == "" {
		sessionID = xid.New().String()
		response.Session = sessionID
	} else if len(data.Session) > 48 {
		response.err(http.StatusBadRequest, "That session ID is too long!")
		return
	} else {
		sessionID = data.Session
	}

	session, gotSessionOk := chatSessions.Get(sessionID).(*ChatSession)
	if gotSessionOk {
		if clientIP != session.ClientIP {
			response.err(http.StatusUnauthorized, "That session is already being used by another client!")
			return
		}
	} else {
		have := 0
		chatSessions.ForEach(func(_ string, val interface{}) bool {
			if val.(*ChatSession).ClientIP == clientIP {
				have++
			}

			return true
		})

		if have >= 100 {
			response.err(http.StatusForbidden, "You may not create more than 100 sessions!")
			return
		}

		session = &ChatSession{
			ID:           sessionID,
			History:      make([]HistoryEntry, 0, 8),
			BeginTime:    time.Now(),
			LastModified: time.Now(),
			ClientIP:     clientIP,
		}
		chatSessions.Set(sessionID, session)
	}

	qLength := utf8.RuneCountInString(data.Query)

	if qLength > 150 {
		data.Query = data.Query[:150]
		qLength = 150
	} else if qLength < 1 {
		response.err(http.StatusBadRequest, "You must specify some text!")
		return
	}

	resp, conf, err := bot.Ask(data.Query, qLength, session, ctx, clientIP)
	if err != nil {
		panic(err)
	}

	response.respond(resp, conf)
}

func newLimiter(limit int64, thresTime time.Duration, method string) *tb_config.Limiter {
	limiter := tollbooth.NewLimiter(limit, thresTime)
	limiter.IPLookups = []string{"CF-Connecting-IP", "X-Forwarded-For", "RemoteAddr", "X-Real-IP"}
	limiter.Methods = []string{method}

	return limiter
}

func newLimiters(limit int64, thresTime time.Duration, methods []string) *tb_config.Limiter {
	limiter := tollbooth.NewLimiter(limit, thresTime)
	limiter.IPLookups = []string{"CF-Connecting-IP", "X-Forwarded-For", "RemoteAddr", "X-Real-IP"}
	limiter.Methods = methods

	return limiter
}

func optimizedLogger(out io.Writer) gin.HandlerFunc {
	isTerm := true

	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) {
		isTerm = false
	}

	logQueue := make(chan string, 2)
	go func() {
		for {
			fmt.Println(<-logQueue)
		}
	}()

	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		ctx.Next()

		latency := time.Since(start)

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		statusCode := ctx.Writer.Status()
		var statusColor, methodColor string
		if isTerm {
			statusColor = colorForStatus(statusCode)
			methodColor = colorForMethod(method)
		}

		logQueue <- fmt.Sprintf("%v |%s %3d %s| %12v | %14s |%s  %s %-7s %s",
			start.Format("15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			path,
		)
	}
}

func panicHandler(ctx *gin.Context, err interface{}) {
	path := ctx.Request.URL.Path

	if path == "/ask" && ctx.Request.Method == "POST" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error":      "An error occurred while processing your request.",
			"response":   "An error occurred!",
			"confidence": 0,
		})
	} else {
		ctx.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"admin": strings.HasPrefix(path, "/admin/"),
			"err":   err,
		})
	}

	ctx.Abort()
}

func e404(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error.tmpl", gin.H{
		"title": "Not Found",
		"desc": template.HTML(`The path you requested doesn't exist on this server!<br>
I wonder where it went... or did you make a typo?`),
		"ctx": ctx,
	})
}

func e405(ctx *gin.Context) {
	ctx.HTML(http.StatusMethodNotAllowed, "error.tmpl", gin.H{
		"desc": template.HTML(`An intruder has been detected in the system. This incident has been reported... to nobody.<br>
<a href="javascript:void(0);" onclick="window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';">Here's a cat video for now.</a>`),
		"err": template.HTML("Error <strong>405 Method Not Allowed</strong>: Intruder detected. Get out of here. You probably shouldn't be here anyway."),
		"ctx": ctx,
	})
}

func templateInit(router *gin.Engine) {
	if tmpl, err := template.New("chatengine").Funcs(funcMap).ParseGlob("templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}

	if tmpl, err := ttemplate.New("chatengine").Funcs(funcMap).ParseGlob("templates/*"); err == nil {
		tmplEngine = tmpl
	} else {
		panic(err)
	}
}

func createHTTPHandler(bot *ChatBot) *gin.Engine {
	wcGlobalKey = genWcGeneralKey()
	wcRvqPath = "/rvq/" + genWcGeneralKey()
	wcTokenSep = genTokenSep()
	wcRand8Garbage = genRand8Garbage()

	router := gin.New()
	if isDebug {
		optimizedLogger(os.Stdout)
	}
	router.Use(nice.Recovery(panicHandler))

	logger.Info("Initializing session store...")
	sessionStore, err := sessions.RedisStore(bot.rPool, []byte("eef906c81b41efa40130c37e645e7d40"))
	if err != nil {
		panic(err)
	}

	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   31556926,
		Secure:   false, // true = https only cookie, not sent over plain http
		HttpOnly: true,
	})

	router.Use(limit.MaxAllowed(28))
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Server", "ChatEngine/1.2.1")
		ctx.Next()
	})

	if isDebug {
		router.Use(func(ctx *gin.Context) {
			templateInit(router)
			ctx.Next()
		})
	}

	templateInit(router)

	router.GET("/chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "chat.tmpl", gin.H{
			"tokenJsGlobalKey": wcGlobalKey,
		})
	})
	router.Static("/static", "static")
	router.Static("/docs", "slate_docs/build")
	router.GET("/", func(ctx *gin.Context) {
		adminKey, _ := ctx.Cookie("adminKey")

		ctx.HTML(http.StatusOK, "landing.tmpl", gin.H{
			"adminKey": adminKey,
		})
	})

	router.GET("/scvfi72a9an1.js", func(ctx *gin.Context) {
		globalKey := ctx.Query("0abpo3inr4")

		ctx.Writer.Header()["Content-Type"] = []string{"application/javascript; charset=utf-8"}

		if globalKey != wcGlobalKey {
			ctx.Status(http.StatusBadRequest)
			ctx.Writer.Write(wcTokenJsErr)
			return
		}

		ctx.Status(http.StatusOK)
		tmplEngine.ExecuteTemplate(ctx.Writer, "webchat_token_obf.js", map[string]interface{}{
			"global_key":  wcGlobalKey,
			"payload_sep": wcTokenSep,
			"target_url":  wcRvqPath,
			"temp_key":    wcTempKey,
		})
	})

	router.Any(wcRvqPath, func(ctx *gin.Context) {
		success := false
		defer func() {
			if !success {
				ctx.String(http.StatusBadRequest, "oops")
			}
		}()

		idx := int(math.Floor(float64(util.CurrentTimeMillis())/180000)) % 3
		correctMethod := wcTokenMethods[idx]
		if ctx.Request.Method != correctMethod {
			return
		}

		key := ctx.Query("v9f0c")
		globalKey := ctx.Query("p87da6b5cz")
		if globalKey != wcGlobalKey {
			return
		}

		expectedKey := wcTempKey
		if key != expectedKey {
			return
		}

		rawPayloadBytes, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}

		rawPayload2, err := strutil.Base64Decode(rawPayloadBytes)
		if err != nil {
			return
		}

		rawPayload := string(rawPayload2)
		payload := strings.Split(rawPayload, wcTokenSep)
		if len(payload) != 6 {
			if wcOldTokenSep != "" && strings.Contains(rawPayload, wcOldTokenSep) {
				payload = strings.Split(rawPayload, wcOldTokenSep)
				if len(payload) != 6 {
					return
				}
			} else {
				return
			}
		}

		ip := ctx.ClientIP()
		random := payload[0]
		host := payload[1]
		href := payload[2]
		userAgent := payload[3]
		path := payload[4]
		features := strings.TrimSpace(payload[5])

		if err2 := tollbooth.LimitByKeys(wcTokenLimiter, []string{ip, href, userAgent, features}); err2 != nil {
			return
		}

		token := &webchatToken{
			token:     genWcTokenString(),
			ip:        ip,
			random:    random,
			host:      host,
			href:      href,
			userAgent: userAgent,
			path:      path,
			features:  features,
			createdAt: time.Now(),
		}
		wcTokenMap.Set(token.token, token)

		ctx.String(http.StatusOK, wcRand8Garbage+token.token)
		success = true
	})

	router.GET("/status", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "status.tmpl", nil)
	})

	router.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		for _, route := range router.Routes() {
			if route.Path == path {
				e405(ctx)
				return
			}
		}

		e404(ctx)
	})
	router.NoMethod(e405)

	// Set up all the components
	go wcTokenSepUpdater()
	go wcTokenJanitor()
	go wcRand8GarbageUpdater()
	setupRegistration(bot, router)
	setupManagement(bot, router)
	setupAPI(bot, router)
	setupAdminInterface(bot, router)
	setupTroubleshooting(bot, router)

	go http.ListenAndServe("127.0.0.1:2084", nil)

	return router
}

func genTokenSep() string {
	targetLen := 12 + rand.Intn(20)
	data := make([]byte, targetLen)
	wtscLen := len(wcTokenSepChars)

	for i := 0; i < targetLen; i++ {
		data[i] = wcTokenSepChars[rand.Intn(wtscLen)]
	}

	return util.BytesToString(data)
}

func cGenRand() int64 {
	res, err := crand.Int(crand.Reader, cgrMax)
	if err != nil {
		return rand.Int63n(math.MaxInt64)
	}

	return res.Int64()
}

func genRand8Garbage() string {
	targetLen := 8
	data := make([]byte, targetLen)
	wtscLen := len(wcTokenSepChars)

	for i := 0; i < targetLen; i++ {
		data[i] = wcTokenSepChars[rand.Intn(wtscLen)]
	}

	return util.BytesToString(data)
}

func genWcGeneralKey() string {
	return fmt.Sprintf("%d%d%d%s%d%d%d", cGenRand(), rand.Int31n(math.MaxInt32), cGenRand(), xid.New().String(),
		cGenRand(), snowNode.Generate(), rand.Int31n(math.MaxInt32))
}

func genWcTokenString() string {
	data := make([]byte, 32)
	_, err := crand.Read(data)
	if err != nil {
		logger.Error("Failed to get 32 bytes of random data for token, using pseudo-random generator.", err)
		for i := 0; i < 32; i++ {
			data[i] = byte(rand.Intn(math.MaxUint8))
		}
	}

	hasher := sha512.New384()
	hasher.Write(data)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func wcRand8GarbageUpdater() {
	for range time.Tick(time.Minute) {
		wcRand8Garbage = genRand8Garbage()
	}
}

func wcTokenSepUpdater() { // every 20 minutes
	for {
		time.Sleep(time.Minute * 10)
		wcOldTokenSep = wcTokenSep
		wcTokenSep = genTokenSep()
		time.Sleep(time.Minute * 10)
		wcOldTokenSep = ""
	}
}

func wcTokenJanitor() {
	for range time.Tick(time.Minute * 5) {
		now := time.Now()
		wcTokenMap.ForEach(func(key string, rToken interface{}) bool {
			token := rToken.(*webchatToken)
			if now.Sub(token.createdAt).Round(time.Minute) >= 5 {
				wcTokenMap.Delete(key)
			}

			return true
		})
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
