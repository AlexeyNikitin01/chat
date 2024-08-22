package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"task4/internal/config"

	"task4/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RunServer(cfg *config.Config) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return
	}
	logger.Instance().Info("server started", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))
	defer l.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var connMap = &sync.Map{}

	clientChan := make(chan net.Conn)
	defer close(clientChan)

	go func() {
		for {
			clientConn, err := l.Accept()
			if err != nil {
				return
			}
			clientChan <- clientConn
		}
	}()

	for {
		select {
		case <-done:
			logger.Instance().Info("server stop",
				zap.String("host", cfg.Host),
				zap.Int("port", cfg.Port))
			connMap.Range(func(key, value any) bool {
				if conn, ok := value.(net.Conn); ok {
					conn.Close()
				}
				return true
			})
			return
		case clientConn := <-clientChan:
			id := uuid.New().String()
			connMap.Store(id, clientConn)
			logger.Instance().Info("new client", zap.String("id", id))
			go userConnect(id, clientConn, connMap)
		}
	}
}

func userConnect(id string, c net.Conn, connMap *sync.Map) {
	defer func() {
		c.Close()
		connMap.Delete(id)
		logger.Instance().Info("client exit", zap.String("id", id))
	}()

	for {
		userInput, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			logger.Instance().Error("no read string from user", zap.String("id", id))
			return
		}

		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if key != id {
					conn.Write([]byte(userInput))
				}
			}
			return true
		})
	}
}
