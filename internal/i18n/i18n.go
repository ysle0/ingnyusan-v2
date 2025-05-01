package i18n

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Supported languages
const (
	LangEnglish = "en"
	LangKorean  = "ko"
	LangFrench  = "fr"
	LangDutch   = "nl"
)

// Translator manages translations for different languages
type Translator struct {
	translations  map[string]map[string]string
	defaultLang   string
	supportedLang map[string]string
	mutex         sync.RWMutex
}

// NewTranslator creates a new translator instance
func NewTranslator(translationsDir, defaultLang string) *Translator {
	t := &Translator{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
		supportedLang: map[string]string{
			LangEnglish: "English",
			LangKorean:  "한국어",
			LangFrench:  "Français",
			LangDutch:   "Nederlands",
		},
		mutex: sync.RWMutex{},
	}

	// Load translation files
	err := t.loadTranslations(translationsDir)
	if err != nil {
		log.Printf("Warning: failed to load translations: %v", err)
	}

	return t
}

// loadTranslations loads all translation files from the given directory
func (t *Translator) loadTranslations(dir string) error {
	// Make sure the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	// Load each language file
	for lang := range t.supportedLang {
		filePath := filepath.Join(dir, lang+".json")

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: could not load translations for %s: %v", lang, err)
			continue
		}

		langMap := make(map[string]string)
		err = json.Unmarshal(data, &langMap)
		if err != nil {
			return err
		}

		t.mutex.Lock()
		t.translations[lang] = langMap
		t.mutex.Unlock()
	}

	return nil
}

// Translate returns the translation for the given key in the specified language
func (t *Translator) Translate(key, lang string) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// If language doesn't exist, use default
	if _, exists := t.translations[lang]; !exists {
		lang = t.defaultLang
	}

	// Try to get translation for requested language
	if translation, exists := t.translations[lang][key]; exists {
		return translation
	}

	// Fallback to default language
	if lang != t.defaultLang {
		if translation, exists := t.translations[t.defaultLang][key]; exists {
			return translation
		}
	}

	// If all else fails, return the key itself
	return key
}

// GetSupportedLanguages returns a map of language codes to their display names
func (t *Translator) GetSupportedLanguages() map[string]string {
	return t.supportedLang
}

// IsValidLanguage checks if the given language code is supported
func (t *Translator) IsValidLanguage(lang string) bool {
	_, exists := t.supportedLang[lang]
	return exists
}

// GetDefaultLanguage returns the default language code
func (t *Translator) GetDefaultLanguage() string {
	return t.defaultLang
}

// T is a shorthand for Translate
func (t *Translator) T(key, lang string) string {
	return t.Translate(key, lang)
}
