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
	"arduino.cc/builder/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CTagsTargetFileSaver struct {
	SourceField string
}

func (s *CTagsTargetFileSaver) Run(context map[string]interface{}) error {
	source := context[s.SourceField].(string)

	preprocPath := context[constants.CTX_PREPROC_PATH].(string)
	err := os.MkdirAll(preprocPath, os.FileMode(0755))
	if err != nil {
		return utils.WrapError(err)
	}

	ctagsTargetFileName := filepath.Join(preprocPath, constants.FILE_CTAGS_TARGET)
	err = ioutil.WriteFile(ctagsTargetFileName, []byte(source), os.FileMode(0644))
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_CTAGS_TEMP_FILE_NAME] = ctagsTargetFileName

	return nil
}
