package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kukymbr/i18n"
	"os"
	"os/signal"
	"syscall"
)

var opts = struct {
	DataType  string
	Dir       string
	Recursive bool
	Verbose   bool
}{
	DataType:  string(i18n.YAML),
	Dir:       ".",
	Recursive: true,
	Verbose:   false,
}

func main() {
	flag.StringVar(&opts.DataType, "type", opts.DataType, "Data type ("+i18n.YAML.String()+"|"+i18n.JSON.String()+")")
	flag.StringVar(&opts.Dir, "dir", opts.Dir, "Path to the directory with the bundle files")
	flag.BoolVar(&opts.Recursive, "recursive", opts.Recursive, "Process directory recursively")
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "Verbose output")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "i18nvalidator: A CLI tool to validate i18n bundle files for missing translations.\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Usage:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")

		flag.PrintDefaults()
	}

	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if opts.Verbose {
		fmt.Println("🚀 Validating i18n bundle files:")
		fmt.Printf("   Data type:      %s\n", opts.DataType)
		fmt.Printf("   Directory: %s\n", opts.Dir)
		fmt.Printf("   Recursive: %t\n\n", opts.Recursive)
	}

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error()+"\n")

		cancel()

		//nolint:gocritic // FP for exitAfterDefer.
		os.Exit(1)
	}

	if opts.Verbose {
		fmt.Println("👍 Bundle has no missing translations!")
	}
}

func run(ctx context.Context) error {
	dt, err := i18n.ParseDataType(opts.DataType)
	if err != nil {
		return err
	}

	bundle, err := i18n.NewBundle(i18n.English, i18n.FromDirs(dt, opts.Recursive, opts.Dir))
	if err != nil {
		return fmt.Errorf("initialize bundle: %w", err)
	}

	errs := i18n.DetectMissingTranslations(ctx, bundle)
	if len(errs) > 0 {
		return errs
	}

	return nil
}
