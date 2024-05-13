package handler

import (
	"context"
	"fmt"
	"go-recorder/facility"
	"go-recorder/repo"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/sirupsen/logrus"
)

func NewHandler(s3 repo.S3, logger *logrus.Logger, cfg *facility.Config) *Handler {
	return &Handler{
		s3:  s3,
		log: logger,
		cfg: cfg,
	}
}

type Handler struct {
	s3  repo.S3
	log *logrus.Logger
	cfg *facility.Config
}

func (h *Handler) VoiceRoomHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type == discordgo.ChannelTypeGuildVoice {
		h.voiceRoomListener(s, c.GuildID, c.ID)
	}
}
func (h *Handler) voiceRoomListener(s *discordgo.Session, gID, cID string) {
	vc, err := s.ChannelVoiceJoin(gID, cID, true, false)
	if err != nil {
		h.log.Errorf("error joining voice room with guild id %s, channel id %s\n", gID, cID)
		return
	}
	defer vc.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.cfg.TimeLimit)*time.Minute)
	defer cancel()

	end := make(chan string)

	go func(ctx context.Context, vc *discordgo.VoiceConnection) {
		file, err := os.CreateTemp("", "temp-*.ogg")
		if err != nil {
			h.log.Errorln("error creating temporary file", err)
			return
		}
		defer os.Remove(file.Name())
		defer file.Close()

		audioWriter, err := oggwriter.NewWith(file, 48000, 2)
		if err != nil {
			h.log.Errorln("error creating audio file", err)
			return
		}
		defer audioWriter.Close()

	loop:
		for {
			select {
			case p := <-vc.OpusRecv:
				packet := createPionRTPPacket(p)
				if err := audioWriter.WriteRTP(packet); err != nil {
					h.log.Errorln("error recording to audio file", err)
					return
				}
			case <-ctx.Done():
				break loop
			}
		}

		filename := fmt.Sprintf("%s_%v", vc.ChannelID, time.Now())

		if _, err := file.Seek(0, 0); err != nil {
			h.log.Errorln(err)
			return
		}

		_, err = h.s3.Upload(context.Background(), file, filename, h.cfg.S3BucketName)
		if err != nil {
			h.log.Errorln("error uploading file", err)
			return
		}

		h.log.Info("Successfully recorded and uploaded audio file!")

		url, err := h.s3.GetURL(context.Background(), h.cfg.S3BucketName, filename)
		if err != nil {
			h.log.Errorln("error getting url file", err)
			return
		}

		end <- url

	}(ctx, vc)

	url := <-end

	if _, err := s.ChannelMessageSend(vc.ChannelID, url); err != nil {
		h.log.Errorln("error sending message", err)
		return
	}

	h.log.Println("EVENT SUCCESSFUL")
}

func createPionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}
