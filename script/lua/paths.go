package lua

import (
	"os"
	"path/filepath"
)

var (
	CurrentDir, _  = os.Getwd()
	SeckillLuaPath = filepath.Join(CurrentDir, "script", "lua", "seckill.lua")
)
