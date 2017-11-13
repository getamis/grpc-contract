// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func ToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func WriteFile(content, filePath, fileName string) {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create folder %v failed: %v\n", filePath, err)
		os.Exit(-1)
	}

	err = ioutil.WriteFile(path.Join(filePath, fileName), []byte(content), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write %v file: %v\n", fileName, err)
		os.Exit(-1)
	}
}
