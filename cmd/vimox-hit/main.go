package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var cliVersion = "0.2.0"

type envelope struct {
	SchemaVersion string `json:"schemaVersion"`
	OK            bool   `json:"ok"`
	Data          any    `json:"data,omitempty"`
	Error         any    `json:"error,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		writeJSON(envelope{SchemaVersion: "1", OK: false, Error: apiError{Code: "COMMAND_FAILED", Message: err.Error()}})
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "version" || args[0] == "--version" {
		writeJSON(envelope{SchemaVersion: "1", OK: true, Data: map[string]string{"name": "vimox-hit", "version": cliVersion}})
		return nil
	}
	switch args[0] {
	case "context":
		return showContext(args[1:])
	case "project":
		return initProject(args[1:])
	case "ai-log":
		return aiLog(args[1:])
	case "doctor":
		return doctor()
	default:
		return errors.New("用法: vimox-hit context show --project ID | project init --name NAME | ai-log list|show|diagnose|copy")
	}
}

func showContext(args []string) error {
	if len(args) < 1 || args[0] != "show" {
		return errors.New("用法: vimox-hit context show --project PROJECT_ID")
	}
	projectID := flagValue(args[1:], "--project")
	if projectID == "" {
		return errors.New("必须指定 --project")
	}
	return getJSON("/projects/"+url.PathEscape(projectID), nil)
}

func aiLog(args []string) error {
	if len(args) < 1 {
		return errors.New("用法: vimox-hit ai-log list|show|diagnose|copy [LOG_ID]")
	}
	switch args[0] {
	case "list":
		query := url.Values{"limit": {"30"}}
		if status := flagValue(args[1:], "--status"); status != "" {
			query.Set("status", status)
		}
		return getJSON("/ai-logs", query)
	case "show", "diagnose", "copy":
		if len(args) < 2 {
			return errors.New("必须指定日志 ID")
		}
		value, err := requestJSON("/ai-logs/"+url.PathEscape(args[1]), nil)
		if err != nil {
			return err
		}
		if args[0] == "show" {
			writeJSON(envelope{SchemaVersion: "1", OK: true, Data: value})
			return nil
		}
		diagnostic, err := marshalDiagnostic(value)
		if err != nil {
			return err
		}
		if args[0] == "diagnose" {
			writeRaw(diagnostic)
			return nil
		}
		if err := copyClipboard(diagnostic); err != nil {
			return err
		}
		writeJSON(envelope{SchemaVersion: "1", OK: true, Data: map[string]string{"message": "脱敏诊断已复制", "logId": args[1]}})
		return nil
	default:
		return fmt.Errorf("不支持的 AI 日志操作: %s", args[0])
	}
}

func doctor() error {
	baseURL := apiBaseURL()
	request, err := http.NewRequest(http.MethodGet, strings.TrimSuffix(baseURL, "/api/v1")+"/health", nil)
	if err != nil {
		return err
	}
	response, err := (&http.Client{Timeout: 15 * time.Second}).Do(request)
	if err != nil {
		return fmt.Errorf("Hit 服务连接失败: %w", err)
	}
	defer response.Body.Close()
	var health json.RawMessage
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&health); err != nil {
		return err
	}
	writeJSON(envelope{SchemaVersion: "1", OK: response.StatusCode >= 200 && response.StatusCode < 300, Data: map[string]any{"api": health, "authenticated": os.Getenv("VIMOX_HIT_TOKEN") != ""}})
	return nil
}

func initProject(args []string) error {
	if len(args) < 1 || args[0] != "init" {
		return errors.New("用法: vimox-hit project init --name NAME")
	}
	name := flagValue(args[1:], "--name")
	if name == "" {
		return errors.New("必须指定 --name")
	}
	root := os.Getenv("VIMOX_HIT_PROJECTS_DIR")
	if root == "" {
		root = filepath.Join(userHome(), "VimoxProjects")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	path := filepath.Join(root, time.Now().Format("2006-01-02")+"_"+slug(name)+"_"+randomSuffix())
	for _, dir := range []string{"input", "analysis", "materials", "output", "metadata", "logs"} {
		if err := os.MkdirAll(filepath.Join(path, dir), 0700); err != nil {
			return err
		}
	}
	manifest := map[string]any{"schemaVersion": 1, "name": name, "status": "draft", "createdAt": time.Now().UTC().Format(time.RFC3339), "files": map[string]string{}}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(path, "project.json"), append(data, '\n'), 0600); err != nil {
		return err
	}
	writeJSON(envelope{SchemaVersion: "1", OK: true, Data: map[string]string{"path": path}})
	return nil
}

func getJSON(path string, query url.Values) error {
	value, err := requestJSON(path, query)
	if err != nil {
		return err
	}
	writeJSON(envelope{SchemaVersion: "1", OK: true, Data: value})
	return nil
}

func requestJSON(path string, query url.Values) (any, error) {
	token := os.Getenv("VIMOX_HIT_TOKEN")
	if token == "" {
		return nil, errors.New("尚未登录：请在 hit.vimox.cn 的 CLI 安装页完成授权并设置 VIMOX_HIT_TOKEN")
	}
	endpoint := apiBaseURL() + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
	if err != nil {
		return nil, fmt.Errorf("Hit API 请求失败: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("Hit API 返回 HTTP %d", response.StatusCode)
	}
	var value any
	if err := json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(&value); err != nil {
		return nil, fmt.Errorf("读取 API 响应失败: %w", err)
	}
	return value, nil
}

func apiBaseURL() string {
	if value := strings.TrimRight(os.Getenv("VIMOX_HIT_API_URL"), "/"); value != "" {
		return value
	}
	return "https://hit.vimox.cn/api/v1"
}

func marshalDiagnostic(value any) ([]byte, error) {
	sanitized := sanitize(value)
	return json.MarshalIndent(map[string]any{"schemaVersion": 1, "source": "vimox-hit ai-log", "log": sanitized}, "", "  ")
}

func sanitize(value any) any {
	switch item := value.(type) {
	case map[string]any:
		clean := make(map[string]any, len(item))
		for key, child := range item {
			lower := strings.ToLower(key)
			if lower == "token" || strings.Contains(lower, "access_token") || strings.Contains(lower, "refresh_token") || strings.Contains(lower, "authorization") || strings.Contains(lower, "api_key") || strings.Contains(lower, "apikey") || strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "cookie") || strings.Contains(lower, "base64") || strings.Contains(lower, "image_data") {
				clean[key] = "[REDACTED]"
				continue
			}
			clean[key] = sanitize(child)
		}
		return clean
	case []any:
		clean := make([]any, len(item))
		for i, child := range item {
			clean[i] = sanitize(child)
		}
		return clean
	case string:
		trimmed := strings.TrimSpace(item)
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			var parsed any
			if json.Unmarshal([]byte(trimmed), &parsed) == nil {
				return sanitize(parsed)
			}
		}
		return scrubText(item)
	default:
		return value
	}
}

func scrubText(value string) string {
	if strings.Contains(value, "data:") && strings.Contains(value, ";base64,") {
		return "[MEDIA DATA REDACTED]"
	}
	value = regexp.MustCompile(`(?i)(bearer\s+)[A-Za-z0-9._~+/=-]+`).ReplaceAllString(value, "${1}[REDACTED]")
	value = regexp.MustCompile(`(?i)((?:access_token|refresh_token|api[_-]?key|secret|password|signature|sign)=)[^&\s\"]+`).ReplaceAllString(value, "${1}[REDACTED]")
	return regexp.MustCompile(`https?://[^\s\"<>]+`).ReplaceAllStringFunc(value, func(raw string) string {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "[URL REDACTED]"
		}
		parsed.RawQuery = ""
		parsed.Fragment = ""
		return parsed.String()
	})
}

func copyClipboard(data []byte) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("clip")
	case "darwin":
		command = exec.Command("pbcopy")
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			command = exec.Command("wl-copy")
		} else if _, err := exec.LookPath("xclip"); err == nil {
			command = exec.Command("xclip", "-selection", "clipboard")
		} else {
			return errors.New("找不到 wl-copy 或 xclip；请使用 ai-log diagnose 输出诊断内容")
		}
	}
	command.Stdin = bytes.NewReader(data)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("复制到系统剪贴板失败: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func flagValue(args []string, flag string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == flag {
			return args[i+1]
		}
	}
	return ""
}

func userHome() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return "."
}

func slug(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`/\:*?"<>|`, r) {
			return '-'
		}
		return r
	}, value)
	return strings.Trim(value, " .-")
}

func randomSuffix() string { return fmt.Sprintf("%x", time.Now().UnixNano()) }

func writeJSON(value envelope) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
}

func writeRaw(value []byte) { _, _ = os.Stdout.Write(append(value, '\n')) }
