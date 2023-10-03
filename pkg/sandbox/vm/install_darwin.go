package vm

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Code-Hex/vz/v3"
)

var installerDir = os.Expand("$XDG_STATE_HOME/devbox/vm", xdgStateHome)

//go:embed bootstrap
var bootstrapFiles embed.FS

var nixosConfigTmpl = template.Must(template.ParseFS(bootstrapFiles, "bootstrap/install.sh", "bootstrap/configuration.nix"))

func (vm *VM) generateBootstrapScript() (dir string, err error) {
	path, err := vm.dataFilePath(filepath.Join("bootstrap", "install.sh"))
	if err != nil {
		return "", fmt.Errorf("generate install script: %v", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("get current user: %v", err)
	}

	localtime, err := os.Readlink("/etc/localtime")
	if err != nil {
		return "", fmt.Errorf("get current time zone: %v", err)
	}
	tz, ok := strings.CutPrefix(localtime, "/var/db/timezone/zoneinfo/")
	if !ok {
		return "", fmt.Errorf("/etc/localtime symlink missing /var/db/timezone/zoneinfo/ prefix")
	}

	f, err := os.Create(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return "", fmt.Errorf("create directory for install script: %v", err)
		}
		f, err = os.Create(path)
	}
	if err != nil {
		return "", err
	}
	defer f.Close()

	data := struct {
		System   string
		TimeZone string
		User     *user.User
	}{
		System:   "aarch64-linux",
		TimeZone: tz,
		User:     currentUser,
	}
	err = nixosConfigTmpl.Execute(f, data)
	if err != nil {
		return "", fmt.Errorf("execute template: %v", err)
	}
	if err := f.Chmod(0o755); err != nil {
		return "", fmt.Errorf("make install script executable: %v", err)
	}
	return filepath.Dir(path), nil
}

func (vm *VM) installerBootLoader(ctx context.Context) (*vz.LinuxBootLoader, error) {
	kernel, init, initrd, err := downloadInstallerKernel(ctx, vm.nixSystem())
	if err != nil {
		return nil, fmt.Errorf("download installer kernel: %v", err)
	}
	params := fmt.Sprintf("console=hvc0 root=/dev/sda init=%s boot.shell_on_fail", init)
	vm.Logger.Debug("created boot loader", "params", params, "installer", vm.Install)
	return vz.NewLinuxBootLoader(kernel,
		vz.WithInitrd(initrd),
		vz.WithCommandLine(params),
	)
}

func downloadInstallerKernel(ctx context.Context, system string) (kernel, init, initrd string, err error) {
	switch system {
	case "aarch64-linux":
		kernel = "/nix/store/kp0454y12fhlivdnv6vpbc0drdijmh32-nixos-system-nixos-23.05.3701.e9b4b56e5a20/kernel"
		init = "/nix/store/kp0454y12fhlivdnv6vpbc0drdijmh32-nixos-system-nixos-23.05.3701.e9b4b56e5a20/init"
		initrd = "/nix/store/kp0454y12fhlivdnv6vpbc0drdijmh32-nixos-system-nixos-23.05.3701.e9b4b56e5a20/initrd"
	default:
		return "", "", "", fmt.Errorf("unsupported system %s", system)
	}

	cmd := exec.CommandContext(ctx, "nix-store", "--realise", kernel, init, initrd)
	if err := cmd.Run(); err != nil {
		return "", "", "", fmt.Errorf("command nix-store --realise: %v", err)
	}
	return kernel, init, initrd, nil
}

func (vm *VM) installerDisk(ctx context.Context) (vz.StorageDeviceConfiguration, error) {
	iso, err := downloadInstallerISO(ctx, vm.nixSystem())
	if err != nil {
		return nil, fmt.Errorf("download installer iso: %v", err)
	}
	attach, err := vz.NewDiskImageStorageDeviceAttachment(iso, true)
	if err != nil {
		return nil, fmt.Errorf("create disk image storage device: %v", err)
	}
	config, err := vz.NewUSBMassStorageDeviceConfiguration(attach)
	if err != nil {
		return nil, fmt.Errorf("configure disk image as USB mass storage device: %v", err)
	}
	return config, nil
}

func downloadInstallerISO(ctx context.Context, system string) (string, error) {
	url := fmt.Sprintf("https://releases.nixos.org/nixos/23.05/nixos-23.05.3701.e9b4b56e5a20/nixos-minimal-23.05.3701.e9b4b56e5a20-%s.iso", system)
	path := filepath.Join(installerDir, filepath.Base(url))
	flag := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	perm := fs.FileMode(0o644)
	f, err := os.OpenFile(path, flag, perm)
	if errors.Is(err, os.ErrExist) {
		return path, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return "", fmt.Errorf("create directory for ISO: %v", err)
		}
		f, err = os.OpenFile(path, flag, perm)
	}
	if err != nil {
		return "", err
	}
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status %s", resp.Status)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close ISO file: %v", err)
	}
	return path, nil
}

// xdgStateHome returns the path for XDG_STATE_HOME or ~/.local/state if it
// isn't set. It can be used with os.Expand.
func xdgStateHome(s string) string {
	switch s {
	case "XDG_STATE_HOME":
		if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
			return xdg
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return os.TempDir()
		}
		return filepath.Join(home, ".local", "state")
	}
	return ""
}
