package main

import (
	"log"
)

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func unescapeKString(s []byte) []byte {
	if len(s) < 2 {
		log.Fatalf("K string should begin and end with '\"'. Its length cannot therefore be less than 2. Actual string: %s", s)
	}
	if s[0] != '"' {
		log.Fatalf("K string should begin with '\"'. Actual string: %s", s)
	}
	if s[len(s)-1] != '"' {
		log.Fatalf("K string should end with '\"'. Actual string: %s", s)
	}
	s = s[1 : len(s)-1]
	// TODO: escape sequences \" \n etc.
	return s
}

func unescapeKLabel(s []byte) []byte {
	if len(s) < 2 {
		log.Fatalf("K string should begin and end with '`'. Its length cannot therefore be less than 2. Actual string: %s", s)
	}
	if s[0] != '`' {
		log.Fatalf("K string should begin with '`\"`'. Actual string: %s", s)
	}
	if s[len(s)-1] != '`' {
		log.Fatalf("K string should end with '`'. Actual string: %s", s)
	}
	s = s[1 : len(s)-1]
	// TODO: escape sequences \` \n etc.
	return s
}
