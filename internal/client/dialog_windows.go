package client

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Wails 2.9.1 panics inside OpenMultipleFilesDialog and OpenDirectoryDialog
// (a nil result is type-asserted after cancel, and the multi-file dialog can
// take the single-file path). These pickers call the system dialogs directly
// and treat cancel as an empty result.

var (
	comdlg32                 = windows.NewLazySystemDLL("comdlg32.dll")
	procGetOpenFileNameW     = comdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
	shell32                  = windows.NewLazySystemDLL("shell32.dll")
	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
)

const (
	ofnAllowMultiSelect = 0x00000200
	ofnExplorer         = 0x00080000
	ofnFileMustExist    = 0x00001000
	ofnPathMustExist    = 0x00000800
	ofnNoChangeDir      = 0x00000008
	ofnHideReadOnly     = 0x00000004
	fnerrBufferTooSmall = 0x3003

	bifReturnOnlyFSDirs = 0x00000001
	bifNewDialogStyle   = 0x00000040
)

type openFileName struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        unsafe.Pointer
	dwReserved        uint32
	flagsEx           uint32
}

type browseInfo struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

// pickFiles opens a multi-select file dialog. Cancel returns nil, nil.
func pickFiles(title string) ([]string, error) {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return nil, err
	}
	filter := utf16Z("All files (*.*)", "*.*")
	owner := uintptr(windows.GetForegroundWindow())
	for n := 32 * 1024; n <= 1024*1024; n *= 2 {
		buf := make([]uint16, n)
		ofn := openFileName{
			hwndOwner:   owner,
			lpstrFilter: &filter[0],
			lpstrFile:   &buf[0],
			nMaxFile:    uint32(len(buf)),
			lpstrTitle:  titlePtr,
			flags:       ofnAllowMultiSelect | ofnExplorer | ofnFileMustExist | ofnPathMustExist | ofnNoChangeDir | ofnHideReadOnly,
		}
		ofn.lStructSize = uint32(unsafe.Sizeof(ofn))
		r, _, callErr := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
		if r != 0 {
			return ofnSelections(buf), nil
		}
		code, _, _ := procCommDlgExtendedError.Call()
		if code == 0 {
			return nil, nil
		}
		if code == fnerrBufferTooSmall {
			continue
		}
		if callErr != nil && callErr != syscall.Errno(0) {
			return nil, fmt.Errorf("не удалось открыть выбор файлов (код %d): %w", code, callErr)
		}
		return nil, fmt.Errorf("не удалось открыть выбор файлов (код %d)", code)
	}
	return nil, fmt.Errorf("слишком много выбранных файлов")
}

// pickFolder opens a folder dialog. Cancel returns "", nil.
func pickFolder(title string) (string, error) {
	release := initCOM()
	if release {
		defer windows.CoUninitialize()
	}
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return "", err
	}
	display := make([]uint16, 260)
	info := browseInfo{
		hwndOwner:      uintptr(windows.GetForegroundWindow()),
		pszDisplayName: &display[0],
		lpszTitle:      titlePtr,
		ulFlags:        bifReturnOnlyFSDirs | bifNewDialogStyle,
	}
	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&info)))
	if pidl == 0 {
		return "", nil
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(pidl))
	pathBuf := make([]uint16, 32768)
	ok, _, _ := procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
	if ok == 0 {
		return "", fmt.Errorf("не удалось прочитать выбранную папку")
	}
	return windows.UTF16ToString(pathBuf), nil
}

func initCOM() bool {
	err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED)
	if err == nil {
		return true
	}
	errno, ok := err.(syscall.Errno)
	// S_FALSE: COM was already initialized on this thread.
	return ok && errno == 1
}

// utf16Z joins strings as a double-NUL terminated UTF-16 list.
func utf16Z(parts ...string) []uint16 {
	var buf []uint16
	for _, part := range parts {
		encoded, err := windows.UTF16FromString(part)
		if err != nil {
			continue
		}
		buf = append(buf, encoded...)
	}
	return append(buf, 0)
}

// ofnSelections parses a GetOpenFileName multi-select buffer.
// One selection is a full path. Several selections are a directory followed by file names.
func ofnSelections(buf []uint16) []string {
	var parts []string
	start := 0
	for i, c := range buf {
		if c != 0 {
			continue
		}
		if i == start {
			break
		}
		parts = append(parts, string(utf16.Decode(buf[start:i])))
		start = i + 1
	}
	if len(parts) <= 1 {
		return parts
	}
	dir := parts[0]
	out := make([]string, 0, len(parts)-1)
	for _, name := range parts[1:] {
		out = append(out, filepath.Join(dir, name))
	}
	return out
}
