package pack

import "github.com/buildpack/pack/builder"

type ImageInfo struct {
	Stack      string
	Buildpacks []builder.BuildpackMetadata
	Base       ImageBase
}

type ImageBase struct {
	Ref string
	TopLayer string
}

func (c *Client) InspectImage(name string, daemon bool) (*ImageInfo, error) {
	return nil, nil
}
