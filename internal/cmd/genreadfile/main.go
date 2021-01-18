package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/lestrrat-go/codegen"
	"github.com/pkg/errors"
)

type definition struct {
	Filename         string // go >= 1.16
	FallbackFilename string // go < 1.16
	Package          string
	ReturnType       string
	ParseOptions     bool
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func _main() error {
	definitions := []definition{
		{
			Package:          "jwk",
			ReturnType:       "Set",
			Filename:         "jwk/io.go",
			FallbackFilename: "jwk/io_go1.15.go",
		},
		{
			Package:          "jws",
			ReturnType:       "*Message",
			Filename:         "jws/io.go",
			FallbackFilename: "jws/io_go1.15.go",
		},
		{
			Package:          "jwe",
			ReturnType:       "*Message",
			Filename:         "jwe/io.go",
			FallbackFilename: "jwe/io_go1.15.go",
		},
		{
			Package:          "jwt",
			ReturnType:       "Token",
			Filename:         "jwt/io.go",
			FallbackFilename: "jwt/io_go1.15.go",
			ParseOptions:     true,
		},
	}

	for _, def := range definitions {
		if err := generateFile(def); err != nil {
			return err
		}
		if err := generateFallbackFile(def); err != nil {
			return err
		}
	}
	return nil
}

func generateFile(def definition) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// +build go1.16")
	fmt.Fprintf(&buf, "\n\n// Automatically generated by internal/cmd/genreadfile/main.go. DO NOT EDIT")
	fmt.Fprintf(&buf, "\n\npackage %s", def.Package)
	fmt.Fprintf(&buf, "\n\nimport \"github.com/lestrrat-go/jwx/internal/fs\"")
	fmt.Fprintf(&buf, "\n\n// ReadFileOption describes an option that can be passed to `ReadFile`")
	fmt.Fprintf(&buf, "\ntype ReadFileOption = fs.OpenOption")
	fmt.Fprintf(&buf, "\n\n// LocalFS is used to explicitly tell ReadFile to read from the")
	fmt.Fprintf(&buf, "\n// local file system")
	fmt.Fprintf(&buf, "\ntype LocalFS = fs.Local")
	fmt.Fprintf(&buf, "\n\n// WithFS creates an option that specifies where to load files from")
	fmt.Fprintf(&buf, "\n// when `ReadFile()` is called. For example, you can specify an")
	fmt.Fprintf(&buf, "\n// instance of `embed.FS` to load files from an embedded file")
	fmt.Fprintf(&buf, "\n// system in a compiled binary.")
	fmt.Fprintf(&buf, "\nfunc WithFS(v fs.FS) ReadFileOption {")
	fmt.Fprintf(&buf, "\nreturn fs.WithFS(v)")
	fmt.Fprintf(&buf, "\n}")
	fmt.Fprintf(&buf, "\n\n// ReadFile reads a JWK set at the specified location.")
	fmt.Fprintf(&buf, "\n//")
	fmt.Fprintf(&buf, "\n// For go >= 1.16 where io/fs is supported, you may pass `WithFS()` option")
	fmt.Fprintf(&buf, "\n// to provide an alternate location to load the files from. This means")
	fmt.Fprintf(&buf, "\n// you can embed your JWK in your compiled program and read it back from")
	fmt.Fprintf(&buf, "\n// the embedded resource.")
	fmt.Fprintf(&buf, "\n//")
	fmt.Fprintf(&buf, "\n// By default files are read from the local file system by using `os.Open`")
	fmt.Fprintf(&buf, "\nfunc ReadFile(path string, options ...ReadFileOption) (%s, error) {", def.ReturnType)
	if def.ParseOptions {
		fmt.Fprintf(&buf, "\nvar parseOptions []ParseOption")
		fmt.Fprintf(&buf, "\nvar readFileOptions []ReadFileOption")
		fmt.Fprintf(&buf, "\nfor _, option := range options {")
		fmt.Fprintf(&buf, "\nswitch option := option.(type) {")
		fmt.Fprintf(&buf, "\ncase ParseOption:")
		fmt.Fprintf(&buf, "\nparseOptions = append(parseOptions, option)")
		fmt.Fprintf(&buf, "\ndefault:")
		fmt.Fprintf(&buf, "\nreadFileOptions = append(readFileOptions, option)")
		fmt.Fprintf(&buf, "\n}")
		fmt.Fprintf(&buf, "\n}")
		fmt.Fprintf(&buf, "\n\nf, err := fs.Open(path, readFileOptions...)")
		fmt.Fprintf(&buf, "\nif err != nil {")
		fmt.Fprintf(&buf, "\nreturn nil, err")
		fmt.Fprintf(&buf, "\n}")
		fmt.Fprintf(&buf, "\n\ndefer f.Close()")
		fmt.Fprintf(&buf, "\nreturn ParseReader(f, parseOptions...)")
		fmt.Fprintf(&buf, "\n}")
	} else {
		fmt.Fprintf(&buf, "\nf, err := fs.Open(path, options...)")
		fmt.Fprintf(&buf, "\nif err != nil {")
		fmt.Fprintf(&buf, "\nreturn nil, err")
		fmt.Fprintf(&buf, "\n}")
		fmt.Fprintf(&buf, "\n\ndefer f.Close()")
		fmt.Fprintf(&buf, "\nreturn ParseReader(f)")
		fmt.Fprintf(&buf, "\n}")
	}
	if err := codegen.WriteFile(def.Filename, &buf, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return errors.Wrapf(err, `failed to write to %s`, def.Filename)
	}
	return nil
}

func generateFallbackFile(def definition) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// +build !go1.16")
	fmt.Fprintf(&buf, "\n\n// Automatically generated by internal/cmd/genreadfile/main.go. DO NOT EDIT")
	fmt.Fprintf(&buf, "\n\npackage %s", def.Package)
	fmt.Fprintf(&buf, "\n\nimport \"github.com/lestrrat-go/jwx/internal/fs\"")
	fmt.Fprintf(&buf, "\n\n// ReadFileOption describes an option that can be passed to `ReadFile`")
	fmt.Fprintf(&buf, "\ntype ReadFileOption = fs.OpenOption")
	fmt.Fprintf(&buf, "\n\n// ReadFile reads a JWK set at the specified location.")
	fmt.Fprintf(&buf, "\n//")
	fmt.Fprintf(&buf, "\n\n// for go >= 1.16 where io/fs is supported, you may pass `WithFS()` option")
	fmt.Fprintf(&buf, "\n// to provide an alternate location to load the files from to provide an ")
	fmt.Fprintf(&buf, "\n// alternate location to load the files from (if you are reading")
	fmt.Fprintf(&buf, "\n// this message, your go (or your go doc) is probably running go < 1.16)")
	fmt.Fprintf(&buf, "\nfunc ReadFile(path string) (%s, error) {", def.ReturnType)
	if def.ParseOptions {
		fmt.Fprintf(&buf, "\nvar parseOptions []ParseOption")
		fmt.Fprintf(&buf, "\nfor _, option := range options {")
		fmt.Fprintf(&buf, "\nswitch option := option.(type) {")
		fmt.Fprintf(&buf, "\ncase ParseOption:")
		fmt.Fprintf(&buf, "\nparseOptions = append(parseOptions, option)")
		fmt.Fprintf(&buf, "\n}")
		fmt.Fprintf(&buf, "\n}")
	}
	fmt.Fprintf(&buf, "\nf, err := os.Open(path)")
	fmt.Fprintf(&buf, "\nif err != nil {")
	fmt.Fprintf(&buf, "\nreturn nil, errors.Wrapf(err, `failed to open %%s`, path)")
	fmt.Fprintf(&buf, "\n}")
	fmt.Fprintf(&buf, "\n\ndefer f.Close()")
	if def.ParseOptions {
		fmt.Fprintf(&buf, "\nreturn ParseReader(f, parseOptions...)")
	} else {
		fmt.Fprintf(&buf, "\nreturn ParseReader(f)")
	}
	fmt.Fprintf(&buf, "\n}")
	if err := codegen.WriteFile(def.FallbackFilename, &buf, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return errors.Wrapf(err, `failed to write to %s`, def.FallbackFilename)
	}
	return nil
}
