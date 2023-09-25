package callbackserver

import (
	"net"
	"net/http"
)

type CallbackServer struct {
	addr     string
	path     string
	listener net.Listener
	server   *http.Server

	respCh chan (Response)
}

type Response struct {
	Error string `json:"error"`
	Code  string `json:"code"`
	State string `json:"state"`
}

func New() *CallbackServer {
	return &CallbackServer{
		addr:   "127.0.0.1:4446",
		path:   "/callback",
		respCh: make(chan (Response)),
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

func (s *CallbackServer) URI() string {
	// TODO: don't hardcode
	return "http://localhost:4446/callback"
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
		s.respCh <- resp

		// TODO: define success and error redirects
		http.Redirect(w, r, "https://www.jetpack.io", http.StatusSeeOther)
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
	return s.server.Close()
}
