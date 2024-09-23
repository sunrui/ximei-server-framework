/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-06 02:08:03
 */

package hypersonic

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/fs"
	"os"
)

// I18n 国际化
type I18n struct {
	dict map[string]map[string]any // 字典
}

// NewI18n 创建国际化
func NewI18n(dir string) (*I18n, error) {
	i18n := I18n{
		dict: make(map[string]map[string]any),
	}

	if err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		var pathByte []byte
		if pathByte, err = os.ReadFile(dir + "/" + path); err != nil {
			return err
		} else {
			var t any
			if err = toml.Unmarshal(pathByte, &t); err != nil {
				return err
			} else {
				i18n.dict[path] = t.(map[string]any)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &i18n, nil
}

// 翻译
func (i18n I18n) t(lang string, key string) (string, bool) {
	errMsg := fmt.Sprintf("Translation key '%s' for language '%s' not found.", key, lang)

	v := i18n.dict[lang+".toml"]
	if v == nil {
		return errMsg, false
	}

	value, ok := v[key].(string)
	if !ok || value == "" {
		v[key] = errMsg
		return v[key].(string), false
	}

	return value, true
}

// T 翻译
func (i18n I18n) T(lang string, key string) string {
	value, _ := i18n.t(lang, key)
	return value
}

// Tf 翻译并格式化
func (i18n I18n) Tf(lang string, key string, args ...any) string {
	value, found := i18n.t(lang, key)
	if found {
		return fmt.Sprintf(value, args...)
	}
	return value
}
