package deploy

import (
	"crypto/sha256"
	"encoding/hex"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed patcher/copilot-patcher.js
var patcherContent []byte

var patcherHash string

func init() {
	h := sha256.Sum256(patcherContent)
	patcherHash = hex.EncodeToString(h[:])
}

func FindGameDir() (string, error) {
	if path := os.Getenv("LT2_GAME_PATH"); path != "" {
		path = filepath.Clean(path)
		if validGameDir(path) {
			return path, nil
		}
		return "", fmt.Errorf("LT2_GAME_PATH (%s) does not point to Legion TD 2", path)
	}

	if runtime.GOOS == "windows" {
		if path := findSteamRegistry(); path != "" && validGameDir(path) {
			return path, nil
		}
	}

	candidates := []string{
		`C:\Program Files (x86)\Steam\steamapps\common\Legion TD 2`,
		`D:\Steam\steamapps\common\Legion TD 2`,
		`E:\Steam\steamapps\common\Legion TD 2`,
		`C:\Program Files\Steam\steamapps\common\Legion TD 2`,
	}
	for _, p := range candidates {
		if validGameDir(p) {
			return p, nil
		}
	}

	return "", fmt.Errorf("Legion TD 2 installation not found; set LT2_GAME_PATH env var")
}

func validGameDir(path string) bool {
	info, err := os.Stat(filepath.Join(path, "Legion TD 2.exe"))
	return err == nil && !info.IsDir()
}

func targetPath(gameDir string) string {
	return filepath.Join(gameDir, "Legion TD 2_Data", "uiresources", "AeonGT", "hud", "js", "copilot-patcher.js")
}

func DeployPatcher(gameDir string) error {
	target := targetPath(gameDir)
	dir := filepath.Dir(target)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	existing, err := os.ReadFile(target)
	if err == nil {
		h := sha256.Sum256(existing)
		if hex.EncodeToString(h[:]) == patcherHash {
			slog.Info("patcher is up-to-date", "path", target)
			return nil
		}
		slog.Info("patcher outdated, updating", "path", target)
	} else {
		slog.Info("patcher not found, installing", "path", target)
	}

	if err := os.WriteFile(target, patcherContent, 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	slog.Info("patcher deployed", "path", target)
	return nil
}

func findSteamRegistry() string {
	appID := "489520"
	key := `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App ` + appID

	cmd := exec.Command("reg", "query", key, "/v", "InstallLocation")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "InstallLocation") {
			parts := strings.SplitN(line, "REG_SZ", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
