package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	env := os.Getenv("ENV")
	err := godotenv.Load(GetEnvFileOrPanic(env))
	if err != nil {
		log.Fatal("error loading .env file", err)
		return
	}

	logger := GetLogger()

	logger.Debug("running in debug mode")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", GetRedisHostOrPanic(), GetRedisPortOrPanic()),
		Password: GetRedisPassword(),
		DB:       GetRedisDBOrPanic(),
	})
	redisStore := NewRedisStore(redisClient)
	rateLimiter := NewRateLimiter(redisStore)

	chatCrawlerDetector := NewChatCrawlerDetector()
	tunnelManager := NewTunnelManager()
	httpServer := NewHTTPServer(redisStore, tunnelManager, chatCrawlerDetector)

	go func() {
		addr := fmt.Sprintf(":%s", GetHTTPPortOrPanic())
		logger.Info("starting http server", "addr", addr)
		if err := httpServer.ListenAndServe(addr); err != nil {
			log.Fatal("failed to start http server", err)
		}
	}()

	go func() {
		for range time.Tick(3 * time.Minute) {
			tunnelDelCount := tunnelManager.CleanUp()
			logger.Info("cleaned up tunnels", "deletedTunnels", tunnelDelCount)
		}
	}()

	privateKey := ssh.HostKeyFile(GetPublicKeyOrPanic())
	sshServer := NewSSHServer(redisStore, rateLimiter, tunnelManager, GetHostOrPanic())
	logger.Info("starting ssh server", "addr", fmt.Sprintf(":%s", GetSSHPortOrPanic()))
	if err := sshServer.ListenAndServe(fmt.Sprintf(":%s", GetSSHPortOrPanic()), nil, privateKey); err != nil {
		logger.Error("failed to start ssh server", "err", err)
	}
}
