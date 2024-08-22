package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"task4/internal/config"
	"task4/internal/logger"
)

func RunClient(cfg *config.Config) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		logger.Instance().Error("client no connect to server",
		 zap.String("client", cfg.ClientName))
		return
	}
	defer conn.Close()

	logger.Instance().Info("client connect to server",
		 zap.String("host", cfg.Host),
		 zap.Int("port", cfg.Port),
		 zap.String("client", cfg.ClientName))

	done := make(chan os.Signal, 1)
	defer close(done)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go readConsole(conn, cfg)
	go readServer(conn, done)

	<-done
	fmt.Println("exit client")		
}

func readConsole(conn net.Conn, cfg *config.Config) {
	for {
		reader, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			logger.Instance().Error("error", zap.Error(err))
			return
		}

		reader = cfg.ClientName + " : " + reader
		_, err = conn.Write([]byte(reader))
		if err != nil {
			logger.Instance().Error("error", zap.Error(err))
			return
		}
	}
}

func readServer(conn net.Conn, done chan os.Signal) {
	for {
		reader, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Instance().Info("EOF", zap.String("SERVER_STOP", conn.LocalAddr().String()))
			done<-os.Interrupt
			return
		}
		fmt.Print(reader)
	}
}
