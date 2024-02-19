// Akamai Provider for Terraform
// Copyright© Akamai Technologies, Inc.  All Rights Reserved.
// Akamai Provider for Terraform is licensed under the Mozilla Public License Version 2.0,
// a copy of which is reproduced in the LICENSE file.
// Akamai Provider for Terraform uses a version of dnsjava that was modified by Akamai.
// dnsjava is used under the terms of the BSD 3-clause license, as shown in the notice below.
//
// ========================
// dnsjava
// https://github.com/dnsjava/dnsjava
//
// Copyright (c) 1998-2019, Brian Wellington
// Copyright (c) 2005 VeriSign. All rights reserved.
// Copyright (c) 2019-2021, dnsjava authors
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
// EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ======================
// dnsjava version used - v3.5.2 (785639f66732d6c10db89abe277981039d367e18)

package txtrecord

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// NormalizeTarget tries to normalize txt record target. If it cannot be normalized in the
// provided format, it escapes it and retries.
func NormalizeTarget(r string) (string, error) {
	normalized, ok := normalizeTarget(r)
	if ok {
		return normalized, nil
	}

	normalized, ok = normalizeTarget(fmt.Sprintf("%q", r))
	if ok {
		return normalized, nil
	}

	return "", fmt.Errorf("normalizing txt record targed '%s' failed", r)
}

// normalizeTarget is a txt record target normalization func compliant with akamai api
func normalizeTarget(in string) (string, bool) {
	var newRdata strings.Builder
	for _, ch := range in {
		if isSafeASCII(ch) {
			newRdata.WriteRune(ch)
		} else {
			newRdata.WriteString(fmt.Sprintf("\\%03d", ch))
		}
	}
	in = newRdata.String()

	tok := newTokenizer(in)
	stgs, err := rdataFromString(tok)
	if err != nil {
		return "", false
	}
	return rrToString(stgs), true
}

func isSafeASCII(ch rune) bool {
	return ch >= 32 && ch <= 127
}

func rrToString(strs []string) string {
	if len(strs) == 0 {
		// always return at least an empty quoted String
		return "\"\""
	}
	var sb strings.Builder
	for i, str := range strs {
		sb.WriteString(byteArrayToString([]byte(str), true))
		if i != len(strs)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func byteArrayToString(array []byte, quote bool) string {
	var sb strings.Builder
	if quote {
		sb.WriteString(`"`)
	}
	for _, value := range array {
		b := value & 0xFF
		if b < 0x20 || b >= 0x7f {
			sb.WriteString("\\")
			sb.WriteString(fmt.Sprintf("%03d", b))
		} else if b == '"' || b == '\\' {
			sb.WriteString("\\")
			sb.WriteRune(rune(b))
		} else {
			sb.WriteRune(rune(b))
		}
	}
	if quote {
		sb.WriteString(`"`)
	}
	return sb.String()
}

func rdataFromString(st *tokenizer) ([]string, error) {
	strings := make([]string, 0, 2)
	for {
		t, err := st.get(false, false)
		if err != nil {
			return nil, err
		}
		if !t.isString() {
			break
		}
		byteArray, err := byteArrayFromString(t.value.String())
		if err != nil {
			return nil, fmt.Errorf("tokenizer exception: %v", err)
		}
		strings = append(strings, string(byteArray))
	}
	return strings, st.unget()
}

func byteArrayFromString(s string) ([]byte, error) {
	array := []byte(s)
	escaped := false
	hasEscapes := false

	for _, item := range array {
		if item == '\\' {
			hasEscapes = true
			break
		}
	}
	if !hasEscapes {
		if len(array) > 255 {
			return nil, errors.New("text string too long")
		}
		return array, nil
	}

	var os bytes.Buffer

	digits := 0
	intval := 0
	for _, value := range array {
		if escaped {
			b := value
			if b >= '0' && b <= '9' {
				digits++
				intval *= 10
				intval += int(b - '0')
				if intval > 255 {
					return nil, errors.New("bad escape")
				}
				if digits < 3 {
					continue
				}
				b = byte(intval)
			} else if digits > 0 {
				return nil, errors.New("bad escape")
			}
			os.WriteByte(b)
			escaped = false
		} else if value == '\\' {
			escaped = true
			digits = 0
			intval = 0
		} else {
			os.WriteByte(value)
		}
	}
	if digits > 0 && digits < 3 {
		return nil, errors.New("bad escape")
	}
	array = os.Bytes()
	if len(array) > 255 {
		return nil, errors.New("text string too long")
	}

	return os.Bytes(), nil
}

const (
	// EOF ...
	EOF int = 0

	// EOL - End of line */
	EOL int = 1

	// WHITESPACE - Whitespace; only returned when wantWhitespace is set */
	WHITESPACE int = 2

	// IDENTIFIER - An identifier (unquoted string) */
	IDENTIFIER int = 3

	// QUOTEDSTRING - A quoted string */
	QUOTEDSTRING int = 4

	// COMMENT - A comment; only returned when wantComment is set */
	COMMENT int = 5
)

type tokenizer struct {
	ungottenToken bool
	current       *token
	line          int
	multiline     int
	is            stream
	sb            *strings.Builder
	quoting       bool
	delimiters    string
}

type stream struct {
	data string
}

const defaultDelimiters string = " \t\n;()\""
const quotes string = "\""

func (s *stream) read() (int, error) {
	if len(s.data) == 0 {
		return -1, io.EOF
	}
	rr := int(s.data[0])
	s.data = s.data[1:]
	return rr, nil
}

func (s *stream) unread(c int) {
	s.data = fmt.Sprintf("%c%s", c, s.data)
}

type token struct {
	tokenType int
	value     *strings.Builder
}

func (t token) isString() bool {
	return t.tokenType == IDENTIFIER || t.tokenType == QUOTEDSTRING
}

func (t *tokenizer) ungetChar(c int) {
	if c == -1 {
		return
	}
	t.is.unread(c)
	if c == '\n' {
		t.line--
	}
}

func (t *tokenizer) getChar() (int, error) {
	b, err := t.is.read()
	if err != nil {
		return -1, err
	}
	c := int(b)
	if c == '\r' {
		b, err = t.is.read()
		if err != nil && err != io.EOF {
			return 0, err
		}
		next := int(b)
		if next != '\n' {
			t.is.unread(next)
		}
		c = '\n'
	}
	if c == '\n' {
		t.line++
	}
	return c, nil
}

func (t *tokenizer) skipWhitespace() (int, error) {
	skipped := 0
	for {
		c, err := t.getChar()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if c != ' ' && c != '\t' && !(c == '\n' && t.multiline > 0) {
			t.ungetChar(c)
			return skipped, nil
		}
		skipped++
	}
}

func newTokenizer(in string) *tokenizer {
	return &tokenizer{
		is: stream{
			data: in,
		},
		delimiters: defaultDelimiters,
		sb:         &strings.Builder{},
		line:       1,
	}
}

func (t *tokenizer) checkUnbalancedParens() error {
	if t.multiline > 0 {
		return errors.New("unbalanced parentheses")
	}
	return nil
}

func (t *tokenizer) setCurrentToken(tokenType int, value *strings.Builder) *token {
	current := &token{
		tokenType: tokenType,
		value:     value,
	}
	t.current = current
	return current
}

func (t *tokenizer) unget() error {
	if t.ungottenToken {
		return errors.New("Cannot unget multiple tokens")
	}
	if t.current.tokenType == EOL {
		t.line--
	}
	t.ungottenToken = true
	return nil
}

//nolint:gocyclo
func (t *tokenizer) get(wantWhitespace bool, wantComment bool) (*token, error) {
	var tokenType int
	var c int

	if t.ungottenToken {
		t.ungottenToken = false
		if t.current.tokenType == WHITESPACE {
			if wantWhitespace {
				return t.current, nil
			}
		} else if t.current.tokenType == COMMENT {
			if wantComment {
				return t.current, nil
			}
		} else {
			if t.current.tokenType == EOL {
				t.line++
			}
			return t.current, nil
		}
	}
	skipped, err := t.skipWhitespace()
	if err != nil {
		return nil, err
	}
	if skipped > 0 && wantWhitespace {
		return t.setCurrentToken(WHITESPACE, nil), nil
	}
	tokenType = IDENTIFIER
	t.sb.Reset()
	for {
		c, _ = t.getChar()
		if c == -1 || strings.ContainsRune(t.delimiters, rune(c)) {
			if c == -1 {
				if t.quoting {
					return nil, errors.New("EOF in quoted string")
				} else if t.sb.Len() == 0 {
					return t.setCurrentToken(EOF, nil), nil
				}
				return t.setCurrentToken(tokenType, t.sb), nil
			}
			if t.sb.Len() == 0 && tokenType != QUOTEDSTRING {
				if c == '(' {
					t.multiline++
					_, err := t.skipWhitespace()
					if err != nil {
						return nil, err
					}
					continue
				} else if c == ')' {
					if t.multiline <= 0 {
						return nil, errors.New("invalid close parenthesis")
					}
					t.multiline--
					_, err := t.skipWhitespace()
					if err != nil {
						return nil, err
					}
					continue
				} else if c == '"' {
					if !t.quoting {
						t.quoting = true
						t.delimiters = quotes
						tokenType = QUOTEDSTRING
					} else {
						t.quoting = false
						t.delimiters = defaultDelimiters
						_, err := t.skipWhitespace()
						if err != nil {
							return nil, err
						}
					}
					continue
				} else if c == '\n' {
					return t.setCurrentToken(EOL, nil), nil
				} else if c == ';' {
					for {
						c, _ = t.getChar()
						if c == '\n' || c == -1 {
							break
						}
						t.sb.WriteRune(rune(c))
					}
					if wantComment {
						t.ungetChar(c)
						return t.setCurrentToken(COMMENT, t.sb), nil
					} else if c == -1 && tokenType != QUOTEDSTRING {
						err := t.checkUnbalancedParens()
						if err != nil {
							return nil, err
						}
						return t.setCurrentToken(EOF, nil), nil
					} else if t.multiline > 0 {
						_, err := t.skipWhitespace()
						if err != nil {
							return nil, err
						}
						t.sb.Reset()
						continue
					}
					return t.setCurrentToken(EOL, nil), nil
				}
				return nil, errors.New("Illegal state")
			}
			t.ungetChar(c)
			break
		} else if c == '\\' {
			c, _ = t.getChar()
			if c == -1 {
				return nil, errors.New("unterminated escape sequence")
			}
			t.sb.WriteRune('\\')
		} else if t.quoting && c == '\n' {
			return nil, errors.New("newline in quoted string")
		}
		t.sb.WriteRune(rune(c))
	}
	if t.sb.Len() == 0 && tokenType != QUOTEDSTRING {
		err := t.checkUnbalancedParens()
		if err != nil {
			return nil, err
		}
		return t.setCurrentToken(EOF, nil), nil
	}
	return t.setCurrentToken(tokenType, t.sb), nil
}
