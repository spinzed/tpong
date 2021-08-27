package main

type Server struct {
	Clients       []Client
	EventListener chan StateEvent
}

func NewServer() *Server {
	s := &Server{
		Clients:       []Client{},
		EventListener: make(chan StateEvent),
	}

	go s.Listen()
	return s
}

func (s *Server) Listen() {
	for e := range s.EventListener {
		s.BroadcastEvent(e)
	}
}

func (s *Server) RegisterClient(cl Client) {
	s.Clients = append(s.Clients, cl)
}

func (s *Server) UnregisterClient(name string) {
	cp := make([]Client, 0)
	for _, c := range s.Clients {
		if c.Name() != name {
			cp = append(cp, c)
		}
	}
	s.Clients = cp
}

func (s *Server) BroadcastEvent(ev StateEvent) {
	for _, cl := range s.Clients {
		cl.RegisterEvent(ev)
	}
}
