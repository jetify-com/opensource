# Experimental Devbox VM

Experimental support for Devbox virtual machines on macOS.

## Usage

The `dxvm` command acts like `devbox shell` except that it launches the Devbox
environment in a VM.

To create a new VM, run the following:

	cd ~/my/project
	dxvm -install
	# Wait for the prompt. This might appear to hang the first time it's run
	# while downloading the NixOS installer.

	mkdir bootstrap
	sudo mount -t virtiofs bootstrap bootstrap
	sudo bootstrap/install.sh
	sudo shutdown now
	# ^C to exit

Now that the VM is bootstrapped, you can launch it any time with:

	cd ~/my/project
	dxvm

The NixOS installer files are cached in `~/.local/state/devbox/vm`. You can
monitor the ISO in this directory to estimate how far along the download is. The
final size should be around 800 MiB.

The first time `dxvm` is run in a Devbox project, it creates a `.devbox/vm`
directory that contains the VM's state and configuration files:

- `log` - error and debug log messages
- `disk.img` - main disk image, typically mounted as root
- `id` - an opaque Virtualization Framework machine ID

The following files can be edited (for example, `echo 4 > cpu`) to adjust the
VM's resources:

- `mem` - the amount of memory (in bytes) to allocate to the VM
- `cpu` - the number of CPUs to allocate to the VM

## Building

This package uses the macOS Virtualization Framework, and therefore needs CGO.
Devbox and Nix are unable to download the macOS SDK directly, so some extra
setup is required:

- macOS Ventura (13) or later
- XCode command line tools (open Xcode at least once to accept the EULA)

To compile and sign `dxvm` run:

	devbox run build

The `devbox run build` script uses `./cmd/dxvmsign` to sign the Go binary, which
allows it to use the Virtualization Framework. It's a small wrapper around
Apple's `codesign` utility.

## Limitations

- Mounting the Devbox project directory was temporarily removed while cleaning
things up. Needs to be brought back.
- Only aarch64-linux is implemented right now. Other systems have been tested,
but they aren't an option in the dxvm command.
- Using ctrl-c to exit has the unfortunate side-effect of making it impossible
to interrupt a program in the VM.
- The host terminal has no way of telling the guest when it has resized (usually
this is done with SIGWINCH). Running less/vim/etc. in the VM might look messed
up. Run `stty cols X rows Y` in the VM to manually set the size of your terminal
window.

# Todo/Ideas

- Support macOS and x86_64-linux.
- macOS Sonoma added support for VM suspend/resume. This would probably make VM
start times a lot faster (maybe instant?).
- Clipboard sharing.
- Expose sockets for services.
- Mount /nix/store as an overlay to share packages between VMs.
- Communicate directly with the host Nix daemon?
- Disk resizing.
- Memory balloon (adjust VM memory at runtime).
- Multiple consoles.

## Docs

Some useful links for learning more about the Virtualization Framework:

- `vz` - Go bindings for Apple's Virtualization Framework
	- <https://github.com/Code-Hex/vz>
	- <https://github.com/Code-Hex/vz/wiki>
	- <https://pkg.go.dev/github.com/Code-Hex/vz/v3>
- Virtualization Framework
	- <https://developer.apple.com/documentation/virtualization>
