package app

import (
	"context"
	"fmt"
	"go-recorder/facility"
	"go-recorder/handler"
	"go-recorder/repo"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, logger *logrus.Logger, cfg *facility.Config) error {
	repo, err := repo.NewS3Repository(ctx, cfg.S3Region, "", cfg.S3AccessKey, cfg.S3SecretKey)
	if err != nil {
		return err
	}

	handler := handler.NewHandler(repo, logger, cfg)

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.BotToken))
	if err != nil {
		return err
	}

	dg.AddHandler(handler.VoiceRoomHandler)

	err = dg.Open()
	if err != nil {
		return err
	}

	logger.Info("Bot is ready to operate!")
	return nil
}
