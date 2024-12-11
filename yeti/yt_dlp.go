package yeti

import (
	"runtime"
)

const (
	ytDlpLinuxAsset   = "yt-dlp_linux"
	ytDlpMacOsAsset   = "yt-dlp_macos"
	ytDlpWindowsAsset = "yt-dlp.exe"
)

func GetYtDlpBinary() string {
	switch runtime.GOOS {
	case "darwin":
		return ytDlpMacOsAsset
	case "linux":
		return ytDlpLinuxAsset
	case "windows":
		return ytDlpWindowsAsset
	default:
		panic("yet is only supported on Windows, macOS and Linux")
	}
}
