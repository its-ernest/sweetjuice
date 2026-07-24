package ui

import (
	"encoding/json"
	"testing"
)

func TestWidgetSerialization(t *testing.T) {
	text := NewText("Hello World").WithTextColor("#FF0000").WithFontSize(18)
	button := NewButton("Click Me").WithBackgroundColor("#00FF00")
	col := NewColumn(text, button).WithPadding(10, 10, 10, 10)

	serialized, err := col.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize column: %v", err)
	}

	if serialized["type"] != "column" {
		t.Errorf("Expected type 'column', got '%v'", serialized["type"])
	}

	// Marshal and unmarshal to check generic JSON representation
	bytes, err := json.Marshal(serialized)
	if err != nil {
		t.Fatalf("Failed to marshal serialized map: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal bytes: %v", err)
	}

	chList, ok := parsed["children"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'children' to be a slice")
	}

	if len(chList) != 2 {
		t.Errorf("Expected 2 children, got %d", len(chList))
	}

	first := chList[0].(map[string]interface{})
	if first["type"] != "text" {
		t.Errorf("Expected first child to be 'text', got '%v'", first["type"])
	}

	firstText := first["text"].(string)
	if firstText != "Hello World" {
		t.Errorf("Expected text to be 'Hello World', got '%s'", firstText)
	}
}
