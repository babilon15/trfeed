package utils

import (
	"net/url"
	"regexp"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func IsAlphanumLower(str string) bool { return regexp.MustCompile(`^[a-z0-9]*$`).MatchString(str) }

func IsValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

func RemoveDiacritics(str string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, err := transform.String(t, str)
	if err != nil {
		return "", err
	}
	return res, nil
}
