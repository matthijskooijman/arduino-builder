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
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"os"
)

type RewriteHardwareKeys struct{}

func (s *RewriteHardwareKeys) Run(context map[string]interface{}) error {
	if !utils.MapHas(context, constants.CTX_PLATFORM_KEYS_REWRITE) {
		return nil
	}

	hardware := context[constants.CTX_HARDWARE].(map[string]*types.Package)
	platformKeysRewrite := context[constants.CTX_PLATFORM_KEYS_REWRITE].(types.PlatforKeysRewrite)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)

	warn := utils.DebugLevel(context) > 0

	for _, aPackage := range hardware {
		for _, platform := range aPackage.Platforms {
			if platform.Properties[constants.REWRITING] != constants.REWRITING_DISABLED {
				for _, rewrite := range platformKeysRewrite.Rewrites {
					if platform.Properties[rewrite.Key] != constants.EMPTY_STRING && platform.Properties[rewrite.Key] == rewrite.OldValue {
						platform.Properties[rewrite.Key] = rewrite.NewValue
						if warn {
							logger.Fprintln(os.Stderr, constants.MSG_WARNING_PLATFORM_OLD_VALUES, platform.Properties[constants.PLATFORM_NAME], rewrite.Key+"="+rewrite.OldValue, rewrite.Key+"="+rewrite.NewValue)
						}
					}
				}
			}
		}
	}

	return nil
}
