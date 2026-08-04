package broker

func (s *Server) routes() {
	s.mux.HandleFunc("/replicate", s.ReplicateHandler)
	s.mux.HandleFunc("/append", s.AppendHandler)
	s.mux.HandleFunc("/read", s.ReadHandler)
	s.mux.HandleFunc("/promote", s.PromoteHandler)
	s.mux.HandleFunc("/metadata", s.MetadataHandler)
}
