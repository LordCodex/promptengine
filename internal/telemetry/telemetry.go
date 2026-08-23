package telemetry

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Telemetry holds diagnostic usage details (fully anonymous, disabled by default)
type Telemetry struct {
	Enabled     bool   `json:"enabled"`
	UserConsent bool   `json:"user_consent"`
	Endpoint    string `json:"endpoint"`
	wg          sync.WaitGroup
}

type Event struct {
	Command   string            `json:"command"`
	Duration  int64             `json:"duration_ms"`
	Timestamp string            `json:"timestamp"`
	OS        string            `json:"os"`
	Arch      string            `json:"arch"`
	Context   map[string]string `json:"context,omitempty"`
}

func NewTelemetry() *Telemetry {
	// Telemetry is strictly disabled by default.
	// Users must set user_consent to true in ~/.promptengine/config.json
	return &Telemetry{
		Enabled:     false,
		UserConsent: false,
		Endpoint:    "https://telemetry.promptengine.dev/v1/event",
	}
}

// Track logs anonymised event telemetry
func (t *Telemetry) Track(event Event) {
	if !t.Enabled || !t.UserConsent {
		return
	}

	event.Timestamp = time.Now().UTC().Format(time.RFC3339)

	t.wg.Add(1)
	// Send non-blocking payload
	go func() {
		defer t.wg.Done()

		data, err := json.Marshal(event)
		if err != nil {
			return
		}

		req, err := http.NewRequest("POST", t.Endpoint, strings.NewReader(string(data)))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()
}

// Flush blocks until all pending telemetry requests are completed.
func (t *Telemetry) Flush() {
	t.wg.Wait()
}

// SetConsent registers user tracking preferences
func (t *Telemetry) SetConsent(consent bool) {
	t.UserConsent = consent
	t.Enabled = consent
}

func init() {
	// Environment variable fallback override
	if os.Getenv("PROMPTENGINE_TELEMETRY_OPT_IN") == "true" {
		// Opt-in explicitly declared
	}
}
