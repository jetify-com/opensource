package impl

import (
	"os"
	"os/exec"
	"strings"

	"go.jetpack.io/runx/impl/types"
)

func Run(args ...string) error {
	parsed, err := parseArgs(args)
	if err != nil {
		return err
	}
	return run(parsed)
}

// TODO: is this the best name for this struct?
type parsedArgs struct {
	Packages []types.PkgRef
	App      string
	Args     []string
}

func run(args parsedArgs) error {
	paths, err := install(args.Packages...)
	if err != nil {
		return err
	}

	bin, err := lookupBin(paths, args.App)
	if err != nil {
		return err
	}

	cmd := exec.Command(bin, args.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = environ(paths)

	err = cmd.Run()
	if err != nil {
		// If the command failed, we want to return the exit code
		// of the command.
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}

func lookupBin(paths []string, bin string) (string, error) {
	oldPATH := os.Getenv("PATH")
	newPATH := strings.Join(paths, ":")
	os.Setenv("PATH", newPATH)
	defer os.Setenv("PATH", oldPATH)

	path, err := exec.LookPath(bin)
	if err != nil {
		return "", err
	}
	return path, nil
}

func environ(paths []string) []string {
	oldPATH := os.Getenv("PATH")
	allPaths := append(paths, oldPATH)
	newPATH := strings.Join(allPaths, ":")
	os.Setenv("PATH", newPATH)
	defer os.Setenv("PATH", oldPATH)

	return os.Environ()
}

func parseArgs(args []string) (parsedArgs, error) {
	result := parsedArgs{
		Packages: []types.PkgRef{},
		Args:     []string{},
	}

	scanningPackages := true
	for _, arg := range args {
		after, found := strings.CutPrefix(arg, "+")
		if found && scanningPackages {
			ref, err := types.NewPkgRef(after)
			if err != nil {
				return parsedArgs{}, err
			}
			result.Packages = append(result.Packages, ref)
			continue
		}

		if !found && scanningPackages {
			scanningPackages = false
			result.App = arg
			continue
		}

		result.Args = append(result.Args, after)
	}
	return result, nil
}
