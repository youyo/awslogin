package jsonout

import (
	"encoding/json"
	"errors"
	"io"
)

// Writer は構造化 JSON 出力を管理する
type Writer struct {
	stdout io.Writer
	stderr io.Writer
}

// New は新しい Writer を作成する
func New(stdout, stderr io.Writer) *Writer {
	return &Writer{stdout: stdout, stderr: stderr}
}

// resultEnvelope は正常系レスポンスの JSON エンベロープ
type resultEnvelope struct {
	Result any `json:"result"`
}

// errorEnvelope はエラーレスポンスの JSON エンベロープ
type errorEnvelope struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WriteResult は stdout に正常系結果を JSON で出力する
func (w *Writer) WriteResult(v any) error {
	return json.NewEncoder(w.stdout).Encode(resultEnvelope{Result: v})
}

// WriteError は stdout に構造化エラーを JSON で出力する
// AppError の場合はそのコードを使い、それ以外は INTERNAL_ERROR にフォールバックする
func (w *Writer) WriteError(err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return json.NewEncoder(w.stdout).Encode(errorEnvelope{
			Error: errorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
	}
	return json.NewEncoder(w.stdout).Encode(errorEnvelope{
		Error: errorDetail{
			Code:    ErrInternal,
			Message: err.Error(),
		},
	})
}

// WriteEvent は stderr にイベントを JSON で出力する（NDJSON: 1行1JSON）
func (w *Writer) WriteEvent(event any) error {
	return json.NewEncoder(w.stderr).Encode(event)
}

// Stderr は stderr の io.Writer を返す（イベント出力用）
func (w *Writer) Stderr() io.Writer {
	return w.stderr
}
