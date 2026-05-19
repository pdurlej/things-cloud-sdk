package main

import (
	"encoding/json"
	"testing"
)

func TestHandleInitialize(t *testing.T) {
	server := &mcpServer{}
	resp, ok := server.handle(rpcRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Method:  "initialize",
	})
	if !ok {
		t.Fatal("initialize did not produce a response")
	}
	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("result = %T, want map", resp.Result)
	}
	if result["protocolVersion"] != protocolVersion {
		t.Fatalf("protocolVersion = %v, want %s", result["protocolVersion"], protocolVersion)
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %#v", resp.Error)
	}
}

func TestToolsListIncludesCoreTools(t *testing.T) {
	names := map[string]bool{}
	for _, tool := range tools() {
		names[tool.Name] = true
		if tool.InputSchema["type"] != "object" {
			t.Fatalf("%s input schema type = %v, want object", tool.Name, tool.InputSchema["type"])
		}
	}
	for _, name := range []string{"list_tasks", "search_tasks", "create_task", "complete_task"} {
		if !names[name] {
			t.Fatalf("missing tool %s", name)
		}
	}
}

func TestCreateTaskDryRunDoesNotRequireCloud(t *testing.T) {
	server := &mcpServer{}
	result, err := server.createTask("Dry run task", "note", "today", true)
	if err != nil {
		t.Fatalf("createTask dry-run failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("dry-run returned tool error: %#v", result)
	}

	content := result.Content[0].Text
	var payload struct {
		Status string `json:"status"`
		UUID   string `json:"uuid"`
		Item   struct {
			E string `json:"e"`
		} `json:"item"`
	}
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		t.Fatalf("unmarshal dry-run content failed: %v", err)
	}
	if payload.Status != "dry-run" {
		t.Fatalf("status = %q, want dry-run", payload.Status)
	}
	if payload.UUID == "" {
		t.Fatal("dry-run uuid is empty")
	}
	if payload.Item.E != "Task6" {
		t.Fatalf("item kind = %q, want Task6", payload.Item.E)
	}
}

func TestCompleteTaskDryRunDoesNotRequireCloud(t *testing.T) {
	server := &mcpServer{}
	result, err := server.completeTask("task-1", true)
	if err != nil {
		t.Fatalf("completeTask dry-run failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("dry-run returned tool error: %#v", result)
	}
	var payload struct {
		Status string `json:"status"`
		UUID   string `json:"uuid"`
	}
	if err := json.Unmarshal([]byte(result.Content[0].Text), &payload); err != nil {
		t.Fatalf("unmarshal dry-run content failed: %v", err)
	}
	if payload.Status != "dry-run" || payload.UUID != "task-1" {
		t.Fatalf("payload = %#v, want dry-run for task-1", payload)
	}
}
