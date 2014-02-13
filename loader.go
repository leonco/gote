package gote

import (
	"io/ioutil"
	"os"
	"path"
)

type TemplateLoader interface {
	GetTemplate(filename string) (*Template, error)
}

type TemplateLoaderContext struct {
	Loader    TemplateLoader
	Directory string
}

func NewTemplateLoaderContext(loader TemplateLoader, directory string) *TemplateLoaderContext {
	return &TemplateLoaderContext{loader, directory}
}

type TemplateCache struct {
	templates   map[string]*Template
	lastUpdated map[string]int64
	basePath    string
	parser      TemplateParser
}

func (cache *TemplateCache) GetTemplate(filename string) (*Template, error) {
	filename = path.Join(cache.basePath, filename)

	fi, err := os.Stat(filename)

	var lastModified int64 = 0
	if err == nil {
		lastModified = fi.ModTime().Unix()
	}

	if cache.inCache(filename, lastModified) {
		//return &cache.templates[filename], nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	results, err := parseTemplate(cache.parser, contents)

	if err != nil {
		return nil, err
	}

	results.context = NewTemplateLoaderContext(cache, path.Base(filename))
	cache.updateCache(filename, results, lastModified)
	return results, nil
}

func (cache *TemplateCache) inCache(filename string, lastModified int64) bool {
	if v, ok := cache.lastUpdated[filename]; ok {
		return v >= lastModified
	}
	return false
}

func (cache *TemplateCache) updateCache(filename string, results *Template, lastModified int64) {
	cache.templates[filename] = results
	cache.lastUpdated[filename] = lastModified
}

func NewTemplateCache(parser TemplateParser, dir string) *TemplateCache {
	templates := make(map[string]*Template)
	updated := make(map[string]int64)
	return &TemplateCache{templates, updated, dir, parser}
}
