package commands

import (
	"bytes"
	"github.com/buildpack/pack"
	"github.com/buildpack/pack/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"text/template"

	"github.com/buildpack/pack/logging"
	"github.com/buildpack/pack/style"
)

func InspectImage(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect-image <builder-image-name>",
		Short: "Show information about a built image",
		Args:  cobra.ExactArgs(1),
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			img := args[0]

			logger.Infof("Inspecting image: %s\n", style.Symbol(img))

			remoteOutput, warnings, err := inspectImageOutput(client, img, cfg, false)
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Infof("REMOTE:\n%s\n", remoteOutput)
				for _, w := range warnings {
					logger.Warn(w)
				}
			}

			localOutput, warnings, err := inspectImageOutput(client, img, cfg, true)
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Infof("\nLOCAL:\n%s\n", localOutput)
				for _, w := range warnings {
					logger.Warn(w)
				}
			}

			return nil
		}),
	}
	AddHelpFlag(cmd, "inspect-image")
	return cmd
}

func inspectImageOutput(client PackClient, img string, cfg config.Config, local bool) (output string, warning []string, err error) {
	source := "remote"
	if local {
		source = "local"
	}

	info, err := client.InspectImage(img, local)
	if err != nil {
		return "", nil, errors.Wrapf(err, "inspecting %s image '%s'", source, img)
	}
	
	if info == nil {
		return "(not present)", nil, nil
	}
	
	var buf bytes.Buffer
	warnings, err := generateImageOutput(&buf, img, cfg, *info)
	if err != nil {
		return "", nil, errors.Wrapf(err, "writing output for %s image '%s'", source, img)
	}
	
	return buf.String(), warnings, nil
	return "", nil, nil
}

func generateImageOutput(writer *bytes.Buffer, imageName string, cfg config.Config, info pack.ImageInfo) (warnings []string, err error) {
	tpl := template.Must(template.New("").Parse(`
Stack: {{ .Info.Stack }}

Buildpacks:
{{- if .Info.Buildpacks }}
{{ .Buildpacks }}
{{- else }}
  (none) 
{{- end }}`,
	))

	bps, err := buildpacksOutput(info.Buildpacks)
	if err != nil {
		return nil, err
	}

	return nil, tpl.Execute(writer, &struct {
		Info       pack.ImageInfo
		Buildpacks string
	}{
		info,
		bps,
	})
}
