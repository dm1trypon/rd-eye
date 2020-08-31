package tools

import (
	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/models"
)

// LC - logging category
const LC = "TOOLS"

/*
RemoveFromBuffer - removes an element from the stream buffer queue and returns array with boolean status.
Returns <[]models.Package, bool> result. Returns false if an error occurred.
	- packagesBuf <[]models.Package> - array of packages stream buffer.

	- index <int> - Index of the item to delete.
*/
func RemoveFromBuffer(packagesBuf []models.Package, index int) ([]models.Package, bool) {
	if index > len(packagesBuf) {
		logger.Error(LC, "Index greater than array length")
		return packagesBuf[:len(packagesBuf)-1], false
	}

	packagesBuf[len(packagesBuf)-1], packagesBuf[index] = packagesBuf[index], packagesBuf[len(packagesBuf)-1]
	return packagesBuf[:len(packagesBuf)-1], true
}
