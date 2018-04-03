// RIF - (HTTP) Requests In Files

// Copyright (C) 2017 - Jonathan Lloyd (jonathan@thisisjonathan.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package validation

import (
	"errors"

	"github.com/jonathanlloyd/rif/internal/app/fileversions"
)

// ValidateRifYamlFile takes a struct that represents an unmarshalled RIF file
// and returns a slice of errors listing any validation errors detected
func ValidateRifYamlFile(rFile fileversions.RifYamlFileV0) []error {
	errs := []error{}

	if rFile.RifVersion == nil {
		errs = append(errs, errors.New("rif_version is required"))
	}

	return errs
}
