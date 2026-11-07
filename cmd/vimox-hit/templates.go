package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type mediaFiles []string

func (m *mediaFiles) String() string         { return strings.Join(*m, ",") }
func (m *mediaFiles) Set(value string) error { *m = append(*m, value); return nil }

func templateCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("用法: vimox-hit template list|upload|publish|stats --kind video|suite")
	}
	flags := flag.NewFlagSet("template "+args[0], flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	kind := flags.String("kind", "video", "video or suite")
	file := flags.String("file", "", "template JSON file")
	id := flags.String("id", "", "template ID")
	private := flags.Bool("private", false, "keep uploaded template private")
	origin := flags.String("origin", "", "system or mine")
	query := flags.String("query", "", "search text")
	cursor := flags.String("cursor", "", "next page cursor")
	limit := flags.Int("limit", 24, "page size (1-60)")
	var media mediaFiles
	flags.Var(&media, "media", "video/image path; repeat in metadata index order")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("模板命令不接受位置参数，请使用 --id 或 --file")
	}
	if *kind != "video" && *kind != "suite" {
		return errors.New("--kind 必须是 video 或 suite")
	}
	kindValue := strings.ToUpper(*kind)
	var result any
	var err error
	switch args[0] {
	case "list":
		if *limit < 1 || *limit > 60 {
			return errors.New("--limit 必须在 1-60 之间")
		}
		values := url.Values{"limit": {fmt.Sprint(*limit)}}
		if *query != "" {
			values.Set("q", *query)
		}
		if *cursor != "" {
			values.Set("cursor", *cursor)
		}
		if *origin != "" {
			if *origin != "system" && *origin != "mine" {
				return errors.New("--origin 必须是 system 或 mine")
			}
			if *kind == "suite" {
				if *origin == "system" {
					values.Set("origin", "builtin")
				} else {
					values.Set("origin", "user")
				}
			} else {
				values.Set("origin", *origin)
			}
		}
		path := "/video-templates"
		if *kind == "suite" {
			path = "/suites"
		}
		result, err = requestJSON(path, values)
	case "stats":
		if *id == "" {
			return errors.New("必须指定 --id")
		}
		result, err = requestJSON("/template-stats", url.Values{"kind": {kindValue}, "templateId": {*id}})
	case "publish":
		if *id == "" {
			return errors.New("必须指定 --id")
		}
		result, err = publishTemplate(kindValue, *id)
	case "upload":
		if *file == "" {
			return errors.New("必须指定 --file 模板 JSON 路径")
		}
		data, readErr := os.ReadFile(*file)
		if readErr != nil {
			return fmt.Errorf("读取模板失败: %w", readErr)
		}
		if len(data) > 200_000 || !json.Valid(data) {
			return errors.New("模板必须是小于 200 KB 的有效 JSON")
		}
		if *kind == "suite" {
			if len(media) > 0 {
				return errors.New("套图媒体请通过模板 provenance.sourceImages 中的 HTTPS 链接提供；--media 仅用于视频模板")
			}
			result, err = requestAPI(http.MethodPost, "/suites", nil, bytes.NewReader(data), "application/json", time.Minute)
		} else if len(media) == 0 {
			result, err = requestAPI(http.MethodPost, "/video-templates/import", nil, bytes.NewReader(data), "application/json", time.Minute)
		} else {
			result, err = uploadVideoTemplate(data, media)
		}
		if err != nil {
			return err
		}
		record, ok := result.(map[string]any)
		if !ok {
			return errors.New("上传响应缺少模板记录")
		}
		templateID, ok := record["id"].(string)
		if !ok || templateID == "" {
			return errors.New("上传响应缺少模板 ID")
		}
		if !*private {
			publication, publishErr := publishTemplate(kindValue, templateID)
			if publishErr != nil {
				return fmt.Errorf("模板已保存为 %s，发布失败: %w；可用 template publish --kind %s --id %s 重试", templateID, publishErr, *kind, templateID)
			}
			result = map[string]any{"template": record, "publication": publication}
		}
	default:
		return errors.New("未知模板命令；可用 list、upload、publish、stats")
	}
	if err != nil {
		return err
	}
	writeJSON(envelope{SchemaVersion: "1", OK: true, Data: result})
	return nil
}

func publishTemplate(kind, id string) (any, error) {
	data, err := json.Marshal(map[string]string{"kind": kind, "templateId": id})
	if err != nil {
		return nil, err
	}
	return requestAPI(http.MethodPost, "/template-publications", nil, bytes.NewReader(data), "application/json", time.Minute)
}

func uploadVideoTemplate(metadata []byte, files []string) (any, error) {
	if len(files) > 20 {
		return nil, errors.New("媒体最多 20 个")
	}
	var total int64
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() || info.Size() > 128<<20 {
			return nil, errors.New("媒体必须是最大 128 MiB 的普通文件")
		}
		total += info.Size()
		if !strings.HasPrefix(mediaType(path), "image/") && !strings.HasPrefix(mediaType(path), "video/") {
			return nil, fmt.Errorf("不支持的媒体格式: %s", filepath.Ext(path))
		}
	}
	if total > 192<<20 {
		return nil, errors.New("媒体总大小不能超过 192 MiB")
	}
	reader, writer := io.Pipe()
	form := multipart.NewWriter(writer)
	defer reader.Close()
	go func() {
		err := form.WriteField("metadata", string(metadata))
		for _, path := range files {
			if err != nil {
				break
			}
			var file *os.File
			file, err = os.Open(path)
			if err != nil {
				break
			}
			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{"name": "media", "filename": filepath.Base(path)}))
			header.Set("Content-Type", mediaType(path))
			var part io.Writer
			part, err = form.CreatePart(header)
			if err == nil {
				_, err = io.Copy(part, file)
			}
			file.Close()
		}
		if err == nil {
			err = form.Close()
		}
		writer.CloseWithError(err)
	}()
	return requestAPI(http.MethodPost, "/video-templates/import", nil, reader, form.FormDataContentType(), 10*time.Minute)
}

func mediaType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp4", ".m4v":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".webm":
		return "video/webm"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
