package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/dto"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api/services"
	"io"
	"net/http"
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

	messages, err := h.messageService.GetMessagesByPublicKey(publicKey)
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
