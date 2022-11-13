package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/google/uuid"
	"github.com/russross/blackfriday/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx          context.Context
	storePath    string
	snippetsPath string
	tags         map[string]int32
	snippets     map[string]Snippet
}

type Config struct {
	SavePosition bool `json:"save_position"`
	Width        int  `json:"width"`
	Height       int  `json:"height"`
	PositionX    int  `json:"position_x"`
	PositionY    int  `json:"position_y"`
}

type Snippet struct {
	Id      string    `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func CodeClipper() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	homeDir, _ := os.UserHomeDir()

	a.ctx = ctx
	a.storePath = homeDir + "/code-clipper"
	a.snippetsPath = a.storePath + "/snippets"

	if _, err := os.Stat(a.storePath); os.IsNotExist(err) {
		err := os.Mkdir(a.storePath, 0644)
		checkErr(err, "Cannot create required directory")
	}

	loadConfig(ctx, a.storePath)
	a.snippets = loadSnippets(ctx, a.snippetsPath, ".json")
	a.tags = processTags(a.snippets)
}

func (a *App) beforeShutdown(ctx context.Context) (prevent bool) {
	height, width := runtime.WindowGetSize(ctx)
	x, y := runtime.WindowGetPosition(ctx)
	config := Config{
		SavePosition: true,
		Width:        height,
		Height:       width,
		PositionX:    x,
		PositionY:    y,
	}

	jsonStruct, _ := json.MarshalIndent(config, "", " ")
	_ = os.WriteFile(a.storePath+"/config.json", jsonStruct, 0644)

	return false
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func loadConfig(ctx context.Context, root string) {
	file, err := os.ReadFile(root + "/config.json")
	if os.IsNotExist(err) {
		return
	}
	data := Config{}
	_ = json.Unmarshal([]byte(file), &data)
	if data.SavePosition {
		runtime.WindowSetSize(ctx, data.Width, data.Height)
		runtime.WindowSetPosition(ctx, data.PositionX, data.PositionY)
	}
}

func loadSnippets(ctx context.Context, root string, ext string) map[string]Snippet {
	var snippets = make(map[string]Snippet)
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		checkErr(e, "Could not load snippets")
		if filepath.Ext(d.Name()) == ".json" {
			file, _ := os.ReadFile(root + "/" + d.Name())
			data := Snippet{
				Id: strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())),
			}
			_ = json.Unmarshal([]byte(file), &data)
			snippets[data.Id] = data
		}
		return nil
	})
	checkErr(err, "Could not load snippets")
	return snippets
}

func processTags(snippets map[string]Snippet) map[string]int32 {
	var allTags = make(map[string]int32)
	for _, s := range snippets {
		for _, t := range s.Tags {
			if _, ok := allTags[t]; ok {
				allTags[t] += 1
			} else {
				allTags[t] = 1
			}
		}
	}
	return allTags
}

// thanks to https://github.com/zupzup/markdown-code-highlight-chroma/blob/main/main.go
func processCodeRendering(mdContent []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(mdContent))
	if err != nil {
		return "", err
	}
	var hlErr error
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		if hlErr != nil {
			hlErr = err
			return
		}
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")
		oldCode := s.Text()
		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)
		formatter := html.New(html.WithClasses(false))
		iterator, err := lexer.Tokenise(nil, string(oldCode))
		if err != nil {
			hlErr = err
			return
		}
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		if err := formatter.Format(buf, styles.Monokai, iterator); err != nil {
			hlErr = err
			return
		}
		if err := buf.Flush(); err != nil {
			hlErr = err
			return
		}
		s.SetHtml(b.String())
	})
	if hlErr != nil {
		return "", hlErr
	}
	new, err := doc.Html()
	if err != nil {
		return "", err
	}
	return new, nil
}

func (a *App) GetUUID() string {
	return uuid.New().String()
}

func (a *App) GetSnippets() map[string]Snippet {
	return a.snippets
}

func (a *App) GetTags() map[string]int32 {
	return a.tags
}

func (a *App) SaveSnippet(id string, body string, title string, tags []string, created string) Snippet {
	dtc, _ := time.Parse(time.RFC3339, created)
	data := Snippet{
		Id:      id,
		Title:   title,
		Body:    body,
		Tags:    tags,
		Created: dtc,
		Updated: time.Now(),
	}

	jsonStruct, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile(a.snippetsPath+"/"+id+".json", jsonStruct, 0644)

	a.snippets[id] = data
	a.tags = processTags(a.snippets)

	return data
}

func (a *App) DeleteSnippet(id string) {
	_ = os.Remove(a.snippetsPath + "/" + id + ".json")
	delete(a.snippets, id)
	a.tags = processTags(a.snippets)
}

func (a *App) RenderSnippet(id string) string {
	htmlContent := blackfriday.Run([]byte(a.snippets[id].Body))
	rendered, e := processCodeRendering(htmlContent)
	checkErr(e, "Could not render snippet")
	return rendered
}
