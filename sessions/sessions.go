package sessions

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const (
	ContextKey  = "session"
	Name        = "main"
	errorFormat = "[sessions] ERROR! %s\n"
)

var Store store

type store interface {
	sessions.Store
	Options(Options)
}

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the session.
	Clear()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error
	// Session returns the internal gorilla session.
	Session() *sessions.Session
}

type session struct {
	name    string
	request *http.Request
	store   store
	session *sessions.Session
	dirty   bool
	writer  http.ResponseWriter
}

func (s *session) Get(key interface{}) interface{} {
	return s.session.Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.session.Values[key] = val
	s.dirty = true
}

func (s *session) Delete(key interface{}) {
	delete(s.session.Values, key)
	s.dirty = true
}

func (s *session) Clear() {
	for key := range s.session.Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.session.AddFlash(value, vars...)
	s.dirty = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.dirty = true
	return s.session.Flashes(vars...)
}

func (s *session) Options(options Options) {
	s.session.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (s *session) Save() error {
	if s.dirty {
		e := s.store.Save(s.request, s.writer, s.session)
		if e == nil {
			s.dirty = false
		}
		return e
	}
	return nil
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *session) Dirty() bool {
	return s.dirty
}

func Get(ctx *gin.Context) Session {
	if s, ok := ctx.Get(ContextKey); ok {
		return s.(Session)
	} else {
		sess, err := Store.Get(ctx.Request, Name)
		if err != nil {
			panic(err)
		}

		s = &session{
			name:    Name,
			request: ctx.Request,
			store:   Store,
			session: sess,
			dirty:   false,
			writer:  ctx.Writer,
		}
		ctx.Set(ContextKey, s)
		return s.(Session)
	}
}
