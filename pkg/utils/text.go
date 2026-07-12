package utils

import (
	"net/url"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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

func FilterEmptyStrings(in []string) []string {
	n := 0
	for i := 0; i < len(in); i++ {
		if in[i] != "" {
			in[n] = in[i]
			n++
		}
	}
	return in[:n]
}
