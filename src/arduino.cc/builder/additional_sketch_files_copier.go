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

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AdditionalSketchFilesCopier struct{}

func (s *AdditionalSketchFilesCopier) Run(context map[string]interface{}) error {
	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	sketchBuildPath := context[constants.CTX_SKETCH_BUILD_PATH].(string)

	err := os.MkdirAll(sketchBuildPath, os.FileMode(0755))
	if err != nil {
		return utils.WrapError(err)
	}

	sketchBasePath := filepath.Dir(sketch.MainFile.Name)

	for _, file := range sketch.AdditionalFiles {
		relativePath, err := filepath.Rel(sketchBasePath, file.Name)
		if err != nil {
			return utils.WrapError(err)
		}

		targetFilePath := filepath.Join(sketchBuildPath, relativePath)
		os.MkdirAll(filepath.Dir(targetFilePath), os.FileMode(0755))

		bytes, err := ioutil.ReadFile(file.Name)
		if err != nil {
			return utils.WrapError(err)
		}

		ioutil.WriteFile(targetFilePath, bytes, os.FileMode(0644))
	}

	return nil
}
