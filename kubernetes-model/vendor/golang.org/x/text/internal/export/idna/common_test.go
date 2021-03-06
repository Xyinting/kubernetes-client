/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package idna

// This file contains code that is common between the generation code and the
// package's test code.

import (
	"log"

	"golang.org/x/text/internal/ucd"
)

func catFromEntry(p *ucd.Parser) (cat category) {
	r := p.Rune(0)
	switch s := p.String(1); s {
	case "valid":
		cat = valid
	case "disallowed":
		cat = disallowed
	case "disallowed_STD3_valid":
		cat = disallowedSTD3Valid
	case "disallowed_STD3_mapped":
		cat = disallowedSTD3Mapped
	case "mapped":
		cat = mapped
	case "deviation":
		cat = deviation
	case "ignored":
		cat = ignored
	default:
		log.Fatalf("%U: Unknown category %q", r, s)
	}
	if s := p.String(3); s != "" {
		if cat != valid {
			log.Fatalf(`%U: %s defined for %q; want "valid"`, r, s, p.String(1))
		}
		switch s {
		case "NV8":
			cat = validNV8
		case "XV8":
			cat = validXV8
		default:
			log.Fatalf("%U: Unexpected exception %q", r, s)
		}
	}
	return cat
}

var joinType = map[string]info{
	"L": joiningL,
	"D": joiningD,
	"T": joiningT,
	"R": joiningR,
}
