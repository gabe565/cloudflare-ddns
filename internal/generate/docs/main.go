package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gabe565.com/cloudflare-ddns/cmd"
	"gabe565.com/cloudflare-ddns/cmd/envs"
	"gabe565.com/cloudflare-ddns/cmd/sources"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		slog.Error("failed to remove existing dir", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		slog.Error("failed to mkdir", "error", err)
		os.Exit(1)
	}

	root := cmd.New(cobrax.WithVersion("beta"), cmd.WithMarkdown())

	if err := errors.Join(
		generateFlagDoc(root, filepath.Join(output, root.Name()+".md")),
		generateEnvDoc(root, filepath.Join(output, root.Name()+"_envs.md")),
		generateSourcesDoc(root, filepath.Join(output, root.Name()+"_sources.md")),
	); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func generateFlagDoc(cmd *cobra.Command, output string) error {
	var buf bytes.Buffer
	if err := doc.GenMarkdown(cmd, &buf); err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	buf.WriteString("### SEE ALSO\n")
	addSeeAlso(&buf, cmd, cmd.Commands()...)

	return os.WriteFile(output, buf.Bytes(), 0o600)
}

func generateEnvDoc(cmd *cobra.Command, output string) error {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{envs.Name})
	if err := cmd.Execute(); err != nil {
		return err
	}

	sourcesCmd, _, err := cmd.Find([]string{sources.Name})
	if err != nil {
		panic(err)
	}

	buf.WriteString("\n### SEE ALSO\n")
	addSeeAlso(&buf, cmd, cmd, sourcesCmd)

	return os.WriteFile(output, buf.Bytes(), 0o600)
}

func generateSourcesDoc(cmd *cobra.Command, output string) error {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{sources.Name})
	if err := cmd.Execute(); err != nil {
		return err
	}

	envCmd, _, err := cmd.Find([]string{envs.Name})
	if err != nil {
		return err
	}

	buf.WriteString("\n### SEE ALSO\n")
	addSeeAlso(&buf, cmd, cmd, envCmd)

	return os.WriteFile(output, buf.Bytes(), 0o600)
}

func addSeeAlso(buf *bytes.Buffer, cmd *cobra.Command, cmds ...*cobra.Command) {
	for _, subcmd := range cmds {
		if subcmd.Name() == "help" {
			continue
		}
		if cmd.Name() == subcmd.Name() {
			fmt.Fprintf(buf, "* [%s](%s.md)  - %s\n", subcmd.Name(), subcmd.Name(), subcmd.Short)
		} else {
			fmt.Fprintf(buf, "* [%s %s](%s_%s.md)  - %s\n", cmd.Name(), subcmd.Name(), cmd.Name(), subcmd.Name(), subcmd.Short)
		}
	}
}
