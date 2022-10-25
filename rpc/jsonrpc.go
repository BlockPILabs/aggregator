package rpc

import (
	"encoding/json"
)

type JsonRpcRequest struct {
	Id      any         `json:"id"`
	JSONRpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JsonRpcResponse struct {
	Id      any                   `json:"id"`
	JSONRpc string                `json:"jsonrpc,omitempty"`
	Error   *JsonRpcResponseError `json:"error,omitempty"`
	Result  interface{}           `json:"result,omitempty"`
}

func (r *JsonRpcResponse) Marshal() []byte {
	data, _ := json.Marshal(r)
	return data
}

type JsonRpcResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewJsonRpcResponseError(code int, message string, data any) *JsonRpcResponseError {
	return &JsonRpcResponseError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewJsonRpcResponse(id any, result any, err *JsonRpcResponseError) *JsonRpcResponse {
	return &JsonRpcResponse{
		Id:      id,
		JSONRpc: "2.0",
		Error:   err,
		Result:  result,
	}
}

func Error(id any, code int, msg string) *JsonRpcResponse {
	return NewJsonRpcResponse(id, nil, NewJsonRpcResponseError(code, msg, nil))
}

func ErrorServerError(id any, msg string) *JsonRpcResponse {
	return NewJsonRpcResponse(id, nil, NewJsonRpcResponseError(-32000, msg, nil))
}

func ErrorInvalidRequest(id any, msg string) *JsonRpcResponse {
	return NewJsonRpcResponse(id, nil, NewJsonRpcResponseError(-32600, msg, nil))
}

func ErrorMethodNotFound(id any, msg string) *JsonRpcResponse {
	return NewJsonRpcResponse(id, nil, NewJsonRpcResponseError(-32601, msg, nil))
}

func ErrorInvalidParams(id any, msg string) *JsonRpcResponse {
	return NewJsonRpcResponse(id, nil, NewJsonRpcResponseError(-32602, msg, nil))
}

func MustUnmarshalJsonRpcRequest(data []byte) *JsonRpcRequest {
	req := &JsonRpcRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil
	}
	return req
}
