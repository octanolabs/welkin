// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// +build none

/*
The ci command is called from Continuous Integration scripts.

Usage: go run build/ci.go <command> <command flags/arguments>

Available commands are:

   install    [ -arch architecture ] [ -cc compiler ] [ packages... ]                          -- builds packages and executables
   test       [ -coverage ] [ packages... ]                                                    -- runs the tests
   lint                                                                                        -- runs certain pre-selected linters
   archive    [ -arch architecture ] [ -type zip|tar ] [ -signer key-envvar ] [ -upload dest ] -- archives build artifacts
   importkeys                                                                                  -- imports signing keys from env
   debsrc     [ -signer key-id ] [ -upload dest ]                                              -- creates a debian source package
   nsis                                                                                        -- creates a Windows NSIS installer
   aar        [ -local ] [ -sign key-id ] [-deploy repo] [ -upload dest ]                      -- creates an Android archive
   xcode      [ -local ] [ -sign key-id ] [-deploy repo] [ -upload dest ]                      -- creates an iOS XCode framework
   xgo        [ -alltools ] [ options ]                                                        -- cross builds according to options
   purge      [ -store blobstore ] [ -days threshold ]                                         -- purges old archives from the blobstore

For all commands, -n prevents execution of external programs (dry run mode).

*/
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/octanolabs/welkin/internal/build"
	"github.com/octanolabs/welkin/params"
)

var (
	// Files that end up in the welkin*.zip archive.
	welkinArchiveFiles = []string{
		"LICENSE.md",
		executablePath("welkin"),
	}

	// Files that end up in the welkin-alltools*.zip archive.
	allToolsArchiveFiles = []string{
		"LICENSE.md",
		executablePath("welkin"),
	}
)

var GOBIN, _ = filepath.Abs(filepath.Join("build", "bin"))

func executablePath(name string) string {
	return filepath.Join(GOBIN, name)
}

func main() {
	log.SetFlags(log.Lshortfile)

	if _, err := os.Stat(filepath.Join("build", "ci.go")); os.IsNotExist(err) {
		log.Fatal("this script must be run from the root of the repository")
	}
	if len(os.Args) < 2 {
		log.Fatal("need subcommand as first argument")
	}
	switch os.Args[1] {
	case "install":
		doInstall(os.Args[2:])
	case "test":
		doTest(os.Args[2:])
	case "lint":
		doLint(os.Args[2:])
	case "archive":
		doArchive(os.Args[2:])
	default:
		log.Fatal("unknown command ", os.Args[1])
	}
}

// Compiling

func doInstall(cmdline []string) {
	var (
		arch = flag.String("arch", "", "Architecture to cross build for")
		cc   = flag.String("cc", "", "C compiler to cross build with")
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	// Check Go version. People regularly open issues about compilation
	// failure with outdated Go. This should save them the trouble.
	if !strings.Contains(runtime.Version(), "devel") {
		// Figure out the minor version number since we can't textually compare (1.10 < 1.9)
		var minor int
		fmt.Sscanf(strings.TrimPrefix(runtime.Version(), "go1."), "%d", &minor)

		if minor < 12 {
			log.Println("You have Go version", runtime.Version())
			log.Println("welkin requires at least Go version 1.12 and cannot")
			log.Println("be compiled with an earlier version. Please upgrade your Go installation.")
			os.Exit(1)
		}
	}
	// Compile packages given as arguments, or everything if there are no arguments.
	packages := []string{"./..."}
	if flag.NArg() > 0 {
		packages = flag.Args()
	}
	packages = build.ExpandPackagesNoVendor(packages)

	if *arch == "" || *arch == runtime.GOARCH {
		goinstall := goTool("install", buildFlags(env)...)
		goinstall.Args = append(goinstall.Args, "-v")
		goinstall.Args = append(goinstall.Args, packages...)
		build.MustRun(goinstall)
		return
	}
	// If we are cross compiling to ARMv5 ARMv6 or ARMv7, clean any previous builds
	if *arch == "arm" {
		os.RemoveAll(filepath.Join(runtime.GOROOT(), "pkg", runtime.GOOS+"_arm"))
		for _, path := range filepath.SplitList(build.GOPATH()) {
			os.RemoveAll(filepath.Join(path, "pkg", runtime.GOOS+"_arm"))
		}
	}
	// Seems we are cross compiling, work around forbidden GOBIN
	goinstall := goToolArch(*arch, *cc, "install", buildFlags(env)...)
	goinstall.Args = append(goinstall.Args, "-v")
	goinstall.Args = append(goinstall.Args, []string{"-buildmode", "archive"}...)
	goinstall.Args = append(goinstall.Args, packages...)
	build.MustRun(goinstall)

	if cmds, err := ioutil.ReadDir("cmd"); err == nil {
		for _, cmd := range cmds {
			pkgs, err := parser.ParseDir(token.NewFileSet(), filepath.Join(".", "cmd", cmd.Name()), nil, parser.PackageClauseOnly)
			if err != nil {
				log.Fatal(err)
			}
			for name := range pkgs {
				if name == "main" {
					gobuild := goToolArch(*arch, *cc, "build", buildFlags(env)...)
					gobuild.Args = append(gobuild.Args, "-v")
					gobuild.Args = append(gobuild.Args, []string{"-o", executablePath(cmd.Name())}...)
					gobuild.Args = append(gobuild.Args, "."+string(filepath.Separator)+filepath.Join("cmd", cmd.Name()))
					build.MustRun(gobuild)
					break
				}
			}
		}
	}
}

func buildFlags(env build.Environment) (flags []string) {
	var ld []string
	if env.Commit != "" {
		ld = append(ld, "-X", "main.gitCommit="+env.Commit)
	}
	if runtime.GOOS == "darwin" {
		ld = append(ld, "-s")
	}

	if len(ld) > 0 {
		flags = append(flags, "-ldflags", strings.Join(ld, " "))
	}
	return flags
}

func goTool(subcmd string, args ...string) *exec.Cmd {
	return goToolArch(runtime.GOARCH, os.Getenv("CC"), subcmd, args...)
}

func goToolArch(arch string, cc string, subcmd string, args ...string) *exec.Cmd {
	cmd := build.GoTool(subcmd, args...)
	cmd.Env = []string{"GOPATH=" + build.GOPATH()}
	if arch == "" || arch == runtime.GOARCH {
		cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
	} else {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
	}
	if cc != "" {
		cmd.Env = append(cmd.Env, "CC="+cc)
	}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	return cmd
}

// Running The Tests
//
// "tests" also includes static analysis tools such as vet.

func doTest(cmdline []string) {
	coverage := flag.Bool("coverage", false, "Whether to record code coverage")
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}
	packages = build.ExpandPackagesNoVendor(packages)

	// Run the actual tests.
	// Test a single package at a time. CI builders are slow
	// and some tests run into timeouts under load.
	gotest := goTool("test", buildFlags(env)...)
	gotest.Args = append(gotest.Args, "-p", "1", "-timeout", "5m")
	if *coverage {
		gotest.Args = append(gotest.Args, "-covermode=atomic", "-cover")
	}

	gotest.Args = append(gotest.Args, packages...)
	build.MustRun(gotest)
}

// runs gometalinter on requested packages
func doLint(cmdline []string) {
	flag.CommandLine.Parse(cmdline)

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}
	// Get metalinter and install all supported linters
	build.MustRun(goTool("get", "gopkg.in/alecthomas/gometalinter.v2"))
	build.MustRunCommand(filepath.Join(GOBIN, "gometalinter.v2"), "--install")

	// Run fast linters batched together
	configs := []string{
		"--vendor",
		"--tests",
		"--deadline=2m",
		"--disable-all",
		"--enable=goimports",
		"--enable=varcheck",
		"--enable=vet",
		"--enable=gofmt",
		"--enable=misspell",
		"--enable=goconst",
		"--min-occurrences=6", // for goconst
	}
	build.MustRunCommand(filepath.Join(GOBIN, "gometalinter.v2"), append(configs, packages...)...)

	// Run slow linters one by one
	for _, linter := range []string{"unconvert", "gosimple"} {
		configs = []string{"--vendor", "--tests", "--deadline=10m", "--disable-all", "--enable=" + linter}
		build.MustRunCommand(filepath.Join(GOBIN, "gometalinter.v2"), append(configs, packages...)...)
	}
}

// Release Packaging
func doArchive(cmdline []string) {
	var (
		arch  = flag.String("arch", runtime.GOARCH, "Architecture cross packaging")
		atype = flag.String("type", "zip", "Type of archive to write (zip|tar)")
		ext   string
	)
	flag.CommandLine.Parse(cmdline)
	switch *atype {
	case "zip":
		ext = ".zip"
	case "tar":
		ext = ".tar.gz"
	default:
		log.Fatal("unknown archive type: ", atype)
	}

	var (
		env = build.Env()

		baseArchive = archiveBasename(*arch, params.ArchiveVersion(env.Commit))
		welkin      = "welkin-" + baseArchive + ext
		alltools    = "welkin-alltools-" + baseArchive + ext
	)
	maybeSkipArchive(env)
	if err := build.WriteArchive(welkin, welkinArchiveFiles); err != nil {
		log.Fatal(err)
	}
	if err := build.WriteArchive(alltools, allToolsArchiveFiles); err != nil {
		log.Fatal(err)
	}
}

func archiveBasename(arch string, archiveVersion string) string {
	platform := runtime.GOOS + "-" + arch
	if arch == "arm" {
		platform += os.Getenv("GOARM")
	}
	if arch == "android" {
		platform = "android-all"
	}
	if arch == "ios" {
		platform = "ios-all"
	}
	return platform + "-" + archiveVersion
}

// skips archiving for some build configurations.
func maybeSkipArchive(env build.Environment) {
	if env.IsPullRequest {
		log.Printf("skipping because this is a PR build")
		os.Exit(0)
	}
	if env.IsCronJob {
		log.Printf("skipping because this is a cron job")
		os.Exit(0)
	}
	if env.Branch != "master" && !strings.HasPrefix(env.Tag, "v1.") {
		log.Printf("skipping because branch %q, tag %q is not on the whitelist", env.Branch, env.Tag)
		os.Exit(0)
	}
}

func getenvBase64(variable string) []byte {
	dec, err := base64.StdEncoding.DecodeString(os.Getenv(variable))
	if err != nil {
		log.Fatal("invalid base64 " + variable)
	}
	return []byte(dec)
}

func makeWorkdir(wdflag string) string {
	var err error
	if wdflag != "" {
		err = os.MkdirAll(wdflag, 0744)
	} else {
		wdflag, err = ioutil.TempDir("", "welkin-build-")
	}
	if err != nil {
		log.Fatal(err)
	}
	return wdflag
}

func isUnstableBuild(env build.Environment) bool {
	if env.Tag != "" {
		return false
	}
	return true
}
