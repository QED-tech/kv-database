package network

import (
	"database/internal/database"
	"database/internal/database/config"
	"database/internal/shared/logger"
	"database/internal/tools"
	"net"
)

type TCPServer struct {
	database *database.Database
	logger   logger.Logger
	config   *config.Config
}

func NewTCPServer(
	database *database.Database,
	logger logger.Logger,
	conf *config.Config,
) *TCPServer {
	return &TCPServer{database: database, logger: logger, config: conf}
}

func (s *TCPServer) Listen() error {
	listener, err := net.Listen("tcp", s.config.Network.Address)
	if err != nil {
		s.logger.Errorf("[server] failed to listen tpc port: %v", err)

		return err
	}

	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			s.logger.Errorf("[server] failed to close listener: %v", err)
		}
	}(listener)

	s.logger.Infof("[server] listening tcp connection on %s", s.config.Network.Address)
	sem := tools.NewSemaphore(s.config.Network.MaxConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Errorf("[server] failed to accept connection: %v", err)
			continue
		}

		go func(c net.Conn) {
			sem.Acquire()
			defer sem.Release()

			s.handleConnection(c)
		}(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			s.logger.Errorf("[server] failed to close connection: %v", err)
		}
	}(conn)

	buf := make([]byte, 2048)

	reads, err := conn.Read(buf)
	if err != nil {
		s.logger.Errorf("[server] failed to read tcp message: %v", err)
		return
	}
	out := s.database.Handle(string(buf[:reads]))

	_, err = conn.Write([]byte(out))
	if err != nil {
		s.logger.Errorf("[server] failed to write tcp message: %v", err)
	}
}
