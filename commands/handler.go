package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	s   *discordgo.Session
	mtx sync.Mutex
}

func NewHandler(session *discordgo.Session) *Handler {
	return &Handler{
		s:   session,
		mtx: sync.Mutex{},
	}
}
