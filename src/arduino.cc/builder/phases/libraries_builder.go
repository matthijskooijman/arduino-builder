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

package phases

import (
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
	"path/filepath"
)

type LibrariesBuilder struct{}

func (s *LibrariesBuilder) Run(context map[string]interface{}) error {
	librariesBuildPath := context[constants.CTX_LIBRARIES_BUILD_PATH].(string)
	buildProperties := context[constants.CTX_BUILD_PROPERTIES].(map[string]string)
	includes := context[constants.CTX_INCLUDE_FOLDERS].([]string)
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	libraries := context[constants.CTX_IMPORTED_LIBRARIES].([]*types.Library)
	verbose := context[constants.CTX_VERBOSE].(bool)
	warningsLevel := context[constants.CTX_WARNINGS_LEVEL].(string)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	err := os.MkdirAll(librariesBuildPath, os.FileMode(0755))
	if err != nil {
		return utils.WrapError(err)
	}

	objectFiles, err := compileLibraries(libraries, librariesBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_OBJECT_FILES_LIBRARIES] = objectFiles

	return nil
}

func compileLibraries(libraries []*types.Library, buildPath string, buildProperties map[string]string, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	var err error
	var objectFiles []string
	for _, library := range libraries {
		objectFiles, err = compileLibrary(objectFiles, library, buildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	}

	return objectFiles, nil

}

func compileLibrary(objectFiles []string, library *types.Library, buildPath string, buildProperties map[string]string, includes []string, verbose bool, warningsLevel string, logger i18n.Logger) ([]string, error) {
	libraryBuildPath := filepath.Join(buildPath, library.Name)

	err := os.MkdirAll(libraryBuildPath, os.FileMode(0755))
	if err != nil {
		return nil, utils.WrapError(err)
	}

	if library.Layout == types.LIBRARY_RECURSIVE {
		objectFiles, err = builder_utils.CompileFilesRecursive(objectFiles, library.SrcFolder, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	} else {
		utilitySourcePath := filepath.Join(library.SrcFolder, constants.LIBRARY_FOLDER_UTILITY)
		_, utilitySourcePathErr := os.Stat(utilitySourcePath)
		if utilitySourcePathErr == nil {
			includes = append(includes, utils.WrapWithHyphenI(utilitySourcePath))
		}
		objectFiles, err = builder_utils.CompileFiles(objectFiles, library.SrcFolder, false, libraryBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
		if err != nil {
			return nil, utils.WrapError(err)
		}

		if utilitySourcePathErr == nil {
			utilityBuildPath := filepath.Join(libraryBuildPath, constants.LIBRARY_FOLDER_UTILITY)
			objectFiles, err = builder_utils.CompileFiles(objectFiles, utilitySourcePath, false, utilityBuildPath, buildProperties, includes, verbose, warningsLevel, logger)
			if err != nil {
				return nil, utils.WrapError(err)
			}
		}
	}

	return objectFiles, nil
}
