package AliAgent

import (
	"os"
	"path/filepath"
)

const defaultPath = defaultSNName

func SNPath() string {
	return filepath.Join(os.Getenv("ProgramData"), defaultSNName)
}
