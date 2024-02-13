package callbackserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type CallbackServer struct {
	addr     string
	path     string
	listener net.Listener
	server   *http.Server
	redirect string

	respCh chan (Response)
}

type Response struct {
	Error string `json:"error"`
	Code  string `json:"code"`
	State string `json:"state"`
}

func New(redirect string) *CallbackServer {
	return &CallbackServer{
		addr: "127.0.0.1:4446",
		path: "/callback",
		// TODO: don't hard code this so that we can make the library generic
		redirect: redirect,
		respCh:   make(chan (Response)),
	}
}

func (s *CallbackServer) Listen() error {
	if s.listener != nil {
		return nil
	}

	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = l
	return nil
}

func URI() string {
	// TODO: don't hardcode
	return "http://127.0.0.1:4446/callback"
}

func (s *CallbackServer) Port() int {
	if s.listener == nil {
		return 0
	}
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *CallbackServer) Start() error {
	if s.listener == nil {
		err := s.Listen()
		if err != nil {
			return err
		}
	}

	s.server = &http.Server{
		Addr: s.addr,
	}

	http.HandleFunc(s.path, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		resp := Response{
			Error: q.Get("error"),
			Code:  q.Get("code"),
			State: q.Get("state"),
		}

		if s.redirect != "" {
			http.Redirect(w, r, s.redirect, http.StatusSeeOther)
		} else {
			_, _ = w.Write([]byte("OK"))
			w.WriteHeader(http.StatusOK)
		}
		s.respCh <- resp
	})

	if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *CallbackServer) WaitForResponse() Response {
	resp := <-s.respCh
	_ = s.Stop()
	return resp
}

func (s *CallbackServer) Stop() error {
	// We try to gracefully shutdown the server for 100ms. After that time we'll shutdown
	// forcefully. The main reason for this is that we want to give the browser a chance
	// to receive and handle the redirect response before we shut down the server.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	return s.server.Shutdown(ctx)
}
