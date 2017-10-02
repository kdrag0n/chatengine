package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"chatengine/sessions"
	"chatengine/util"
	"github.com/didip/tollbooth"
	tb_config "github.com/didip/tollbooth/config"
	"github.com/gin-gonic/gin"
)

const (
	recaptchaVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
)

var (
	emailRegexp     = regexp.MustCompile(`(?i)^[^@]+@[^@]+\.[a-z\-]+$`)
	registerLimiter = newLimiter(2, time.Second, "POST")
	loginLimiter    = newLimiter(100, time.Hour, "POST")
	httpClient      = &http.Client{
		Timeout: time.Second * 8,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
)

// User represents an user in the system.
type User struct {
	ID          int64     `gorm:"not null"`
	LastUpdated time.Time `gorm:"not null"`

	Email          string `gorm:"not null;type:varchar(100)"`
	Verified       bool   `gorm:"not null"`
	RegistrationIP string `gorm:"not null;type:varchar(45)"`
	Password       []byte `gorm:"type:binary(60)"`
}

func ratelimit(limiter *tb_config.Limiter, r *http.Request) bool {
	return tollbooth.LimitByRequest(limiter, r) != nil
}

func setupRegistration(bot *ChatBot, router *gin.Engine) {
	bot.db.AutoMigrate(&User{})

	router.GET("/signin", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "login")
	})
	router.GET("/login", func(ctx *gin.Context) {
		session := sessions.Get(ctx)
		if v, ok := session.Get("logged_in").(bool); ok && v {
			ctx.Redirect(http.StatusTemporaryRedirect, "manage/dashboard")
			return
		}

		ctx.HTML(http.StatusOK, "login.tmpl", nil)
	})
	router.POST("/login", func(ctx *gin.Context) {
		session := sessions.Get(ctx)
		if v, ok := session.Get("logged_in").(bool); ok && v {
			ctx.Redirect(http.StatusTemporaryRedirect, "manage/dashboard")
			return
		}

		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		rememberMe := ctx.PostForm("remember_me") == "on"
		errors := make([]string, 0, 2)
		passLen := len(password)

		errorEmail := false
		errorPassword := false

		if email == "" || !emailRegexp.MatchString(email) || len(email) > 100 {
			errorEmail = true
			errors = append(errors, "Please enter a valid e-mail address")
		}
		if passLen < 6 {
			errorPassword = true
			errors = append(errors, "You must enter a password with at least 6 characters")
		} else if passLen > 36 {
			errorPassword = true
			errors = append(errors, "Your password may not be longer than 36 characters")
		}

		if len(errors) > 0 {
			ctx.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"errorEmail":    errorEmail,
				"errorPassword": errorPassword,
				"errors":        errors,
			})
			return
		}

		if ratelimit(loginLimiter, ctx.Request) {
			ctx.HTML(http.StatusTooManyRequests, "login.tmpl", gin.H{
				"errors": []string{
					"Please slow down your requests",
				},
			})
			return
		}

		user := User{}
		bot.db.Where("email = ?", email).First(&user)
		passCorrect := false
		if len(user.Password) > 0 {
			passCorrect = bcrypt.CompareHashAndPassword(user.Password, util.StringToBytes(password)) == nil
		}

		if user.Email == "" || !passCorrect {
			ctx.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"errorEmail":    true,
				"errorPassword": true,
				"errors": []string{
					"Invalid e-mail or password",
				},
			})
			return
		}

		session.Set("user_id", user.ID)
		session.Set("logged_in_at", time.Now())
		session.Set("login_ip", ctx.ClientIP())
		session.Set("logged_in", true)
		session.Set("remember_me", rememberMe)
		session.Save()

		ctx.HTML(http.StatusOK, "redirect.tmpl", gin.H{
			"mode":        "relative",
			"destination": "manage/dashboard",
			"ctx":         ctx,
		})
	})

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.tmpl", gin.H{
			"recaptchaKey": config.RecaptchaSiteKey,
		})
	})
	router.POST("/register", func(ctx *gin.Context) {
		session := sessions.Get(ctx)
		if v, ok := session.Get("logged_in").(bool); ok && v {
			ctx.Redirect(http.StatusTemporaryRedirect, "manage/dashboard")
			return
		}

		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		confirm := ctx.PostForm("confirm")
		errors := make([]string, 0, 3)
		passLen := len(password)
		user := User{}

		errorEmail := false
		errorPassword := false
		errorConfirm := false

		if email == "" || !emailRegexp.MatchString(email) || len(email) > 100 {
			errorEmail = true
			errors = append(errors, "Please enter a valid e-mail address")
		}
		if passLen < 6 {
			errorPassword = true
			errors = append(errors, "You must enter a password with at least 6 characters")
		} else if passLen > 36 {
			errorPassword = true
			errors = append(errors, "Your password may not be longer than 36 characters")
		}
		if password != confirm {
			errorConfirm = true
			errors = append(errors, "Your passwords must match")
		}

		if len(errors) < 1 {
			bot.db.Where("email = ?", email).First(&user)
			if user.Email != "" {
				errorEmail = true
				errors = append(errors, "That account already exists")
			}
		}

		if len(errors) > 0 {
			ctx.HTML(http.StatusUnauthorized, "register.tmpl", gin.H{
				"errorEmail":    errorEmail,
				"errorPassword": errorPassword,
				"errorConfirm":  errorConfirm,
				"errors":        errors,
				"recaptchaKey":  config.RecaptchaSiteKey,
			})
			return
		}

		if ratelimit(registerLimiter, ctx.Request) {
			ctx.HTML(http.StatusTooManyRequests, "register.tmpl", gin.H{
				"errors": []string{
					"Please slow down your requests",
				},
				"recaptchaKey": config.RecaptchaSiteKey,
			})
			return
		}

		ip := ctx.ClientIP()

		recaptchaResponse := ctx.PostForm("g-recaptcha-response")
		isCaptchaDone := recaptchaResponse != ""
		if isCaptchaDone {
			recaptchaPostData := "secret=" + config.RecaptchaSecretKey + "&response=" + recaptchaResponse + "&remoteip=" + ip
			resp, err := httpClient.Post(recaptchaVerifyURL, "application/x-www-form-urlencoded", strings.NewReader(recaptchaPostData))

			if err != nil {
				logger.Errorf("Error requesting reCAPTCHA siteverify: %s", err)
				ctx.HTML(http.StatusServiceUnavailable, "register.tmpl", gin.H{
					"errors": []string{
						"We were unable to verify your request",
					},
					"recaptchaKey": config.RecaptchaSiteKey,
				})
				return
			}
			defer resp.Body.Close()

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("Error reading reCAPTCHA siteverify response: %s", err)
				ctx.HTML(http.StatusServiceUnavailable, "register.tmpl", gin.H{
					"errors": []string{
						"We were unable to verify your request",
					},
					"recaptchaKey": config.RecaptchaSiteKey,
				})
				return
			}

			response := make(map[string]interface{}, 5)
			err = json.Unmarshal(respBytes, &response)
			if err != nil {
				logger.Errorf("Error decoding reCAPTCHA siteverify response: %s", err)
				ctx.HTML(http.StatusInternalServerError, "register.tmpl", gin.H{
					"errors": []string{
						"We were unable to verify your request",
					},
					"recaptchaKey": config.RecaptchaSiteKey,
				})
				return
			}

			isCaptchaDone, _ = response["success"].(bool)
		}

		if !isCaptchaDone {
			ctx.HTML(http.StatusForbidden, "register.tmpl", gin.H{
				"errors": []string{
					"You must verify you are a human",
				},
				"recaptchaKey": config.RecaptchaSiteKey,
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(util.StringToBytes(password), 10)
		if err != nil {
			logger.Errorf("Error hashing password (register): %s", err)

			ctx.HTML(http.StatusInternalServerError, "register.tmpl", gin.H{
				"errors": []string{
					"An error occurred creating your account",
				},
				"recaptchaKey": config.RecaptchaSiteKey,
			})
			return
		}

		user = User{
			ID:             snowNode.Generate().Int64(),
			LastUpdated:    time.Now(),
			Email:          email,
			Verified:       false,
			RegistrationIP: ip,
			Password:       hashedPassword,
		}
		bot.db.Create(&user)

		session.Set("user_id", user.ID)
		session.Set("logged_in_at", time.Now())
		session.Set("login_ip", ip)
		session.Set("logged_in", true)
		session.Set("remember_me", true)
		session.Save()

		ctx.HTML(http.StatusOK, "redirect.tmpl", gin.H{
			"mode":        "relative",
			"destination": "manage/dashboard",
			"ctx":         ctx,
		})
	})

	router.GET("/signout", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "logout")
	})
	router.GET("/logout", func(ctx *gin.Context) {
		session := sessions.Get(ctx)
		if v, ok := session.Get("logged_in").(bool); ok && v {
			session.Set("logged_in", false)
		}
		session.Save()

		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
