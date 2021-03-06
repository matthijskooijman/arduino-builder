/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package types

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/props"
	"arduino.cc/builder/utils"
)

type SketchFile struct {
	Name   string
	Source string
}

type SketchFileSortByName []SketchFile

func (s SketchFileSortByName) Len() int {
	return len(s)
}

func (s SketchFileSortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SketchFileSortByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type Sketch struct {
	MainFile         SketchFile
	OtherSketchFiles []SketchFile
	AdditionalFiles  []SketchFile
}

type Package struct {
	PackageId string
	Platforms map[string]*Platform
}

type Platform struct {
	PlatformId   string
	Folder       string
	DefaultBoard *Board
	Boards       map[string]*Board
	Properties   map[string]string
	Programmers  map[string]map[string]string
}

func (platform *Platform) GetTool(toolName string) map[string]string {
	return props.SubTree(platform.Properties, toolName)
}

type Board struct {
	BoardId    string
	Properties map[string]string
}

type Tool struct {
	Name    string
	Version string
	Folder  string
}

type LibraryLayout uint16

const (
	LIBRARY_FLAT LibraryLayout = 1 << iota
	LIBRARY_RECURSIVE
)

type Library struct {
	Folder     string
	SrcFolder  string
	Layout     LibraryLayout
	Name       string
	Archs      []string
	IsLegacy   bool
	Version    string
	Author     string
	Maintainer string
	Sentence   string
	Paragraph  string
	URL        string
	Category   string
	License    string
}

func (library *Library) SupportsArchitectures(archs []string) bool {
	if utils.SliceContains(archs, constants.LIBRARY_ALL_ARCHS) || utils.SliceContains(library.Archs, constants.LIBRARY_ALL_ARCHS) {
		return true
	}

	for _, libraryArch := range library.Archs {
		if utils.SliceContains(archs, libraryArch) {
			return true
		}
	}

	return false
}

type PlatforKeysRewrite struct {
	Rewrites []PlatforKeyRewrite
}

type PlatforKeyRewrite struct {
	Key      string
	OldValue string
	NewValue string
}

type KeyValuePair struct {
	Key   string
	Value string
}

type Command interface {
	Run(context map[string]interface{}) error
}
