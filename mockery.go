package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shoenig/mockery/libmockery"
)

type flags struct {
	version bool
	iface   string
	stdout  bool
	comment string
	pkgname string
}

type environ struct {
	importPrefix   string
	verifyNoChange bool
}

func main() {
	config := parseFlags(os.Args)

	if config.version {
		fmt.Println("mockery " + libmockery.Version)
		return
	}

	if config.iface == "" {
		fmt.Println("-interface is required")
		os.Exit(1)
	}

	if config.pkgname == "" {
		fmt.Println("-package is required")
		os.Exit(1)
	}

	env := parseEnvironment()

	visitor := &libmockery.GeneratorVisitor{
		Comment:           config.comment,
		OutputProvider:    outputProvider(config),
		OutputPackageName: config.pkgname,
		ImportPrefix:      env.importPrefix,
	}

	walker := libmockery.Walker{
		BaseDir:   ".",
		Interface: config.iface,
	}

	if env.verifyNoChange {
		hasher := libmockery.NewHasher(walker.BaseDir)
		hashBefore, err := hasher.Hash()
		if err != nil {
			fmt.Println("mockery failed to hash file contents:", err)
			os.Exit(1)
		}

		defer func() {
			fmt.Println("ensuring no file contents were changed...")
			hashAfter, err := hasher.Hash()
			if err != nil {
				fmt.Println("mockery failed to hash file contents:", err)
				os.Exit(1)
			}

			if err := libmockery.Same(hashBefore, hashAfter); err != nil {
				fmt.Println("mockery unexpectedly modified files:", err)
				os.Exit(1)
			}
		}()
	}

	generated := walker.Walk(visitor)

	if !generated {
		fmt.Printf("Unable to find interface %q in any go files under this path\n", config.iface)
		os.Exit(1)
	}
}

func outputProvider(config flags) libmockery.OutputStreamProvider {
	if config.stdout {
		return &libmockery.StdoutStreamProvider{}
	}
	return &libmockery.FileOutputStreamProvider{
		BaseDir: config.pkgname,
	}
}

func parseFlags(args []string) flags {
	config := flags{}

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.BoolVar(&config.version, "version", false, "print the version of this mockery executable")
	flagSet.StringVar(&config.iface, "interface", "", "name or matching regular expression of interface to generate mock for")
	flagSet.BoolVar(&config.stdout, "stdout", false, "print the generated mock to stdout instead of writing to disk")
	flagSet.StringVar(&config.comment, "comment", "", "comment to insert into prologue of each generated file")
	flagSet.StringVar(&config.pkgname, "package", "", "package name containing generated mocks")

	flagSet.Parse(args[1:])

	return config
}

func parseEnvironment() environ {
	prefix := os.Getenv("MOCKERY_IMPORT_PREFIX")
	nochange := os.Getenv("MOCKERY_CHECK_NOCHANGE") == "1"
	return environ{
		importPrefix:   prefix,
		verifyNoChange: nochange,
	}
}
