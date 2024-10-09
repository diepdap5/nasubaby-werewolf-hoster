package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"nasubaby-werewolf-hoster/cmd/commands"

	"github.com/bwmarrin/discordgo"
)

func handler(w http.ResponseWriter, r *http.Request) {
	public_key := os.Getenv("PUBLIC_KEY")
	bot_token := os.Getenv("BOT_TOKEN")
	app_id := os.Getenv("APP_ID")
	if public_key == "" || bot_token == "" || app_id == "" {
		http.Error(w, "couldn't get environment variables", http.StatusInternalServerError)
		return
	}

	pubkey_b, err := hex.DecodeString(public_key)
	if err != nil {
		http.Error(w, "couldn't decode the public key", http.StatusInternalServerError)
		return
	}
	pubkey := ed25519.PublicKey(pubkey_b)

	XSig := r.Header.Get("X-Signature-Ed25519")
	if XSig == "" {
		http.Error(w, "Missing 'X-Signature-Ed25519' header", http.StatusBadRequest)
		return
	}

	XTime := r.Header.Get("X-Signature-Timestamp")
	if XTime == "" {
		http.Error(w, "Missing 'X-Signature-Timestamp' header", http.StatusBadRequest)
		return
	}

	body := r.Body
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	SignedData := []byte(XTime + string(bodyBytes))
	XSigB, err := hex.DecodeString(XSig)
	if err != nil {
		http.Error(w, "Invalid 'X-Signature-Ed25519' header", http.StatusUnauthorized)
		return
	}

	if !ed25519.Verify(pubkey, SignedData, XSigB) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Authorized
	var inter discordgo.Interaction
	err = json.Unmarshal(bodyBytes, &inter)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch inter.Type {
	case discordgo.InteractionPing:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type":1}`))
	case discordgo.InteractionApplicationCommand:
		s, err := discordgo.New("Bot " + bot_token)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		handler, ok := commands.Commands[inter.ApplicationCommandData().Name]
		if !ok {
			http.Error(w, "Command not found", http.StatusNotFound)
			return
		}

		interactionCreate := discordgo.InteractionCreate{Interaction: &inter}
		err = handler.Handler(s, &interactionCreate, inter.ApplicationCommandData().Name)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Command executed successfully"}`))
	default:
		http.Error(w, "Unknown interaction type", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/interactions", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Register the /ping command
	bot_token := os.Getenv("BOT_TOKEN")
	app_id := os.Getenv("APP_ID")
	dg, err := discordgo.New("Bot " + bot_token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	_, err = dg.ApplicationCommandCreate(app_id, "", &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	})
	if err != nil {
		log.Fatalf("Error creating command: %v", err)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
