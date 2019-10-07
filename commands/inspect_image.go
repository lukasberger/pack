package commands

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/buildpack/pack"
	"github.com/buildpack/pack/config"

	"github.com/buildpack/pack/logging"
	"github.com/buildpack/pack/style"
)

type InspectImageFlags struct {
	BOM bool
}

func InspectImage(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	var flags InspectImageFlags
	cmd := &cobra.Command{
		Use:   "inspect-image <image-name>",
		Short: "Show information about a built image",
		Args:  cobra.ExactArgs(1),
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			img := args[0]

			remote, err := client.InspectImage(img, false)
			if err != nil {
				return errors.Wrapf(err, "inspecting remote image '%s'", img)
			}
			local, err := client.InspectImage(img, true)
			if err != nil {
				return errors.Wrapf(err, "inspecting local image '%s'", img)
			}
			if flags.BOM {
				bom := BOM{
					Remote: remote.BOM,
					Local:  local.BOM,
				}
				rawBOM, err := json.Marshal(bom)
				if err != nil {
					return errors.Wrapf(err, "writing bill of materials for image '%s'", img)
				}
				logger.Infof(string(rawBOM))
				return nil
			}

			logger.Errorf("Inspecting image: %s\n", style.Symbol(img))
			remoteOutput, warnings, err := inspectImageOutput(remote, cfg, false)
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Infof("REMOTE:\n%s\n", remoteOutput)
				for _, w := range warnings {
					logger.Warn(w)
				}
			}

			localOutput, warnings, err := inspectImageOutput(local, cfg, true)
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
	cmd.Flags().BoolVar(&flags.BOM, "bom", false, "print bill of materials")
	return cmd
}

type BOM struct {
	Remote interface{}
	Local  interface{}
}

func inspectImageOutput(
	info *pack.ImageInfo,
	cfg config.Config,
	local bool,
) (output string, warning []string, err error) {
	source := "remote"
	if local {
		source = "local"
	}

	if info == nil {
		return "(not present)", nil, nil
	}

	var buf bytes.Buffer
	warnings, err := generateImageOutput(&buf, cfg, *info)
	if err != nil {
		return "", nil, errors.Wrapf(err, "writing output for %s image", source)
	}

	return buf.String(), warnings, nil
}

func generateImageOutput(writer *bytes.Buffer, cfg config.Config, info pack.ImageInfo) (warnings []string, err error) {
	tpl := template.Must(template.New("").Parse(`
StackID: {{ .Info.StackID }}

Base Image:
  Reference: {{ .Info.Base.Ref }}
  TopLayer: {{ .Info.Base.TopLayer }}

Run Images:
{{- if ne .RunImages "" }}
{{ .RunImages }}
{{- else }}
  (none) 
{{- end }}

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

	runImages, err := runImagesOutput(info.Stack.RunImage.Image, info.Stack.RunImage.Mirrors, cfg)
	if err != nil {
		return nil, err
	}

	return nil, tpl.Execute(writer, &struct {
		Info       pack.ImageInfo
		Buildpacks string
		RunImages  string
	}{
		info,
		bps,
		runImages,
	})
}
