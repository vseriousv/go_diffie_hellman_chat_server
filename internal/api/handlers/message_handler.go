package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/dto"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/services"
	"io"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) GetMessagesByPublicKey(w http.ResponseWriter, r *http.Request) {
	publicKey := chi.URLParam(r, "publicKey")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit int64 = 10
	var offset int64 = 0

	var err error
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	messages, err := h.messageService.GetMessagesByPublicKey(publicKey, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var createDTO dto.CreateMessageDTO

	err = json.Unmarshal(requestBody, &createDTO)
	if err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	id, err := h.messageService.CreateMessage(createDTO)
	if err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}
