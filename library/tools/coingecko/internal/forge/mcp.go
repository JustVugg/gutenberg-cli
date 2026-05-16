package forge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type rpcMessage struct {
	JSONRPC string         `json:"jsonrpc,omitempty"`
	ID      any            `json:"id,omitempty"`
	Method  string         `json:"method,omitempty"`
	Params  map[string]any `json:"params,omitempty"`
}

func RunMCP(input io.Reader, output io.Writer) error {
	reader := bufio.NewReader(input)
	for {
		message, err := readRPC(reader)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		response, err := handleRPC(message)
		if err != nil {
			response = map[string]any{"jsonrpc": "2.0", "id": message.ID, "error": map[string]any{"code": -32000, "message": err.Error()}}
		}
		if response != nil {
			if err := writeRPC(output, response); err != nil {
				return err
			}
		}
	}
}

func handleRPC(message rpcMessage) (map[string]any, error) {
	switch message.Method {
	case "initialize":
		return map[string]any{"jsonrpc": "2.0", "id": message.ID, "result": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{"tools": map[string]any{}},
			"serverInfo": map[string]any{"name": "coingecko-mcp", "version": "0.1.0"},
		}}, nil
	case "notifications/initialized":
		return nil, nil
	case "tools/list":
		return map[string]any{"jsonrpc": "2.0", "id": message.ID, "result": map[string]any{"tools": []map[string]any{
			{"name": "coingecko_operations", "description": "List operations", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
			{"name": "coingecko_call", "description": "Call an operation", "inputSchema": map[string]any{"type": "object", "required": []string{"operationId"}, "properties": map[string]any{"operationId": map[string]any{"type": "string"}, "params": map[string]any{"type": "object"}, "path": map[string]any{"type": "object"}, "yes": map[string]any{"type": "boolean"}}}},
			{"name": "coingecko_search_cache", "description": "Search cache", "inputSchema": map[string]any{"type": "object", "required": []string{"query"}, "properties": map[string]any{"query": map[string]any{"type": "string"}}}},
		}}}, nil
	case "tools/call":
		name, _ := message.Params["name"].(string)
		args, _ := message.Params["arguments"].(map[string]any)
		if name == "coingecko_operations" {
			return toolResult(message.ID, Operations()), nil
		}
		if name == "coingecko_search_cache" {
			query, _ := args["query"].(string)
			results, err := SearchCache(query)
			if err != nil {
				return nil, err
			}
			return toolResult(message.ID, results), nil
		}
		if name == "coingecko_call" {
			operationID, _ := args["operationId"].(string)
			result, err := CallOperation(context.Background(), operationID, CallOptions{QueryParams: mapAnyToString(args["params"]), PathParams: mapAnyToString(args["path"]), Body: args["body"], Yes: args["yes"] == true})
			if err != nil {
				return nil, err
			}
			return toolResult(message.ID, result), nil
		}
		return nil, fmt.Errorf("unknown tool: %s", name)
	default:
		return map[string]any{"jsonrpc": "2.0", "id": message.ID, "error": map[string]any{"code": -32601, "message": "method not found"}}, nil
	}
}

func toolResult(id any, value any) map[string]any {
	content, _ := json.MarshalIndent(value, "", "  ")
	return map[string]any{"jsonrpc": "2.0", "id": id, "result": map[string]any{"content": []map[string]any{{"type": "text", "text": string(content)}}}}
}

func readRPC(reader *bufio.Reader) (rpcMessage, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return rpcMessage{}, err
	}
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "{") {
		var message rpcMessage
		return message, json.Unmarshal([]byte(line), &message)
	}
	if !strings.HasPrefix(strings.ToLower(line), "content-length:") {
		return rpcMessage{}, fmt.Errorf("unexpected MCP frame header: %s", line)
	}
	length, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:")))
	if err != nil {
		return rpcMessage{}, err
	}
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return rpcMessage{}, err
		}
		if strings.TrimSpace(header) == "" {
			break
		}
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(reader, body); err != nil {
		return rpcMessage{}, err
	}
	var message rpcMessage
	return message, json.Unmarshal(body, &message)
}

func writeRPC(output io.Writer, message map[string]any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "Content-Length: %d\r\n\r\n%s", len(body), body)
	return err
}

func mapAnyToString(value any) map[string]string {
	output := map[string]string{}
	input, ok := value.(map[string]any)
	if !ok {
		return output
	}
	for key, item := range input {
		output[key] = fmt.Sprint(item)
	}
	return output
}
