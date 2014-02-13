package gote

import (
	"strings"
)

type TemplateDictionary struct {
	dict          map[string]string
	subs          map[string][]*TemplateDictionary
	shownSections []string
	parent        *TemplateDictionary
}

func NewTemplateDictionary() *TemplateDictionary {
	dict := make(map[string]string)
	return &TemplateDictionary{dict: dict}
}

func (td *TemplateDictionary) Put(key, val string) {
	key = strings.ToUpper(key)
	_, ok := td.dict[key]
	if !ok {
		td.dict[key] = val
	}
}

func (td *TemplateDictionary) Get(key string) (val string, ok bool) {
	key = strings.ToUpper(key)
	val, ok = td.dict[key]
	return
}

func (td *TemplateDictionary) ChildDicts(key string) []*TemplateDictionary {
	return td.subs[strings.ToUpper(key)]
}

func (td *TemplateDictionary) AddChildDict(key string) *TemplateDictionary {
	dict := &TemplateDictionary{parent: td}
	key = strings.ToUpper(key)
	dicts, ok := td.subs[key]
	if ok {
		td.subs[key] = append(dicts, dict)
	} else {
		td.subs[key] = []*TemplateDictionary{dict}
	}
	return dict
}

func (td *TemplateDictionary) AddChildDictAndShowSection(section string) *TemplateDictionary {
	td.ShowSection(section)
	return td.AddChildDict(section)
}

func (td *TemplateDictionary) ShowSection(section string) {
	for _, sec := range td.shownSections {
		if sec == section {
			return
		}
	}
	td.shownSections = append(td.shownSections, section)
}

func (td *TemplateDictionary) isHidenSection(section string) bool {
	for _, sec := range td.shownSections {
		if sec == section {
			return false
		}
	}
	return true
}
