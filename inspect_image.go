package pack

import (
	"context"
	"encoding/json"

	"github.com/buildpack/lifecycle/metadata"
	"github.com/pkg/errors"

	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/image"
)

type ImageInfo struct {
	StackID    string
	Buildpacks []builder.BuildpackMetadata
	Base       ImageBase
	BOM        interface{}
	Stack      metadata.StackMetadata
}

type ImageBase struct {
	Ref      string
	TopLayer string
}

func (c *Client) InspectImage(name string, daemon bool) (*ImageInfo, error) {
	img, err := c.imageFetcher.Fetch(context.Background(), name, daemon, false)
	if err != nil {
		if errors.Cause(err) == image.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	layersMd, err := metadata.GetLayersMetadata(img)
	if err != nil {
		return nil, err
	}
	rawBuildMd, err := metadata.GetRawMetadata(img, metadata.BuildMetadataLabel)
	if err != nil {
		return nil, err
	}
	var buildMD metadata.BuildMetadata
	if err := json.Unmarshal([]byte(rawBuildMd), &buildMD); err != nil {
		return nil, err
	}

	buildpacks := make([]builder.BuildpackMetadata, 0, len(buildMD.Buildpacks))
	for _, bp := range buildMD.Buildpacks {
		buildpacks = append(buildpacks, builder.BuildpackMetadata{
			BuildpackInfo: builder.BuildpackInfo{
				ID:      bp.ID,
				Version: bp.Version,
			},
		})
	}

	stackID, err := metadata.GetRawMetadata(img, metadata.StackMetadataLabel)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		StackID:    stackID,
		Stack:      layersMd.Stack,
		Buildpacks: buildpacks,
		Base: ImageBase{
			Ref:      layersMd.RunImage.Reference,
			TopLayer: layersMd.RunImage.TopLayer,
		},
		BOM: buildMD.BOM,
	}, nil
}
