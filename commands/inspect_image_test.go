package commands_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/buildpack/pack"
	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/commands"
	cmdmocks "github.com/buildpack/pack/commands/mocks"
	"github.com/buildpack/pack/config"
	ilogging "github.com/buildpack/pack/internal/logging"
	"github.com/buildpack/pack/logging"
	h "github.com/buildpack/pack/testhelpers"
)

func TestInspectImageCommand(t *testing.T) {
	color.Disable(true)
	defer func() { color.Disable(false) }()
	spec.Run(t, "Commands", testInspectImageCommand, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testInspectImageCommand(t *testing.T, when spec.G, it spec.S) {

	var (
		command        *cobra.Command
		logger         logging.Logger
		outBuf         bytes.Buffer
		mockController *gomock.Controller
		mockClient     *cmdmocks.MockPackClient
		cfg            config.Config
	)

	it.Before(func() {
		cfg = config.Config{
			RunImages: []config.RunImage{
				{Image: "some/run-image", Mirrors: []string{"first/local", "second/local"}},
			},
		}
		mockController = gomock.NewController(t)
		mockClient = cmdmocks.NewMockPackClient(mockController)
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)

		command = commands.InspectImage(logger, cfg, mockClient)
	})

	it.After(func() {
		mockController.Finish()
	})

	when("#InspectImage", func() {
		when("image cannot be found", func() {
			it("logs 'Not present'", func() {
				mockClient.EXPECT().InspectImage("some/image", false).Return(nil, nil)
				mockClient.EXPECT().InspectImage("some/image", true).Return(nil, nil)

				command.SetArgs([]string{"some/image"})
				h.AssertNil(t, command.Execute())

				h.AssertContains(t, outBuf.String(), "REMOTE:\n(not present)\n\nLOCAL:\n(not present)\n")
			})
		})

		when("inspector returns an error", func() {
			it("logs the error message", func() {
				mockClient.EXPECT().InspectImage("some/image", false).Return(nil, errors.New("some remote error"))
				mockClient.EXPECT().InspectImage("some/image", true).Return(nil, errors.New("some local error"))

				command.SetArgs([]string{"some/image"})
				h.AssertNil(t, command.Execute())

				h.AssertContains(t, outBuf.String(), `ERROR: inspecting remote image 'some/image': some remote error`)
				h.AssertContains(t, outBuf.String(), `ERROR: inspecting local image 'some/image': some local error`)
			})
		})

		when.Pend("the image has empty fields in info", func() {
			it.Before(func() {
				mockClient.EXPECT().InspectBuilder("some/image", false).Return(&pack.BuilderInfo{
					Stack: "test.stack.id",
				}, nil)

				mockClient.EXPECT().InspectBuilder("some/image", true).Return(&pack.BuilderInfo{
					Stack: "test.stack.id",
				}, nil)

				command.SetArgs([]string{"some/image"})
			})

			it("missing creator info is skipped", func() {
				h.AssertNil(t, command.Execute())
				h.AssertNotContains(t, outBuf.String(), "Created By:")
			})

			it("missing description is skipped", func() {
				h.AssertNil(t, command.Execute())
				h.AssertNotContains(t, outBuf.String(), "Description:")
			})

			it("missing buildpacks logs a warning", func() {
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Buildpacks:\n  (none)")
				h.AssertContains(t, outBuf.String(), "Warning: 'some/image' has no buildpacks")
				h.AssertContains(t, outBuf.String(), "Users must supply buildpacks from the host machine")
			})

			it("missing groups logs a warning", func() {
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Detection Order:\n  (none)")
				h.AssertContains(t, outBuf.String(), "Warning: 'some/image' does not specify detection order")
				h.AssertContains(t, outBuf.String(), "Users must build with explicitly specified buildpacks")
			})

			it("missing run image logs a warning", func() {
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Run Images:\n  (none)")
				h.AssertContains(t, outBuf.String(), "Warning: 'some/image' does not specify a run image")
				h.AssertContains(t, outBuf.String(), "Users must build with an explicitly specified run image")
			})

			it("missing lifecycle version prints assumed", func() {
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Lifecycle:\n  Version: 0.3.0")
			})
		})

		when("is successful", func() {
			var (
				buildpack1Info = builder.BuildpackInfo{ID: "test.bp.one", Version: "1.0.0"}
				buildpack2Info = builder.BuildpackInfo{ID: "test.bp.two", Version: "2.0.0"}
				buildpacks     = []builder.BuildpackMetadata{
					{BuildpackInfo: buildpack1Info, Latest: true},
					{BuildpackInfo: buildpack2Info, Latest: false},
				}
				remoteInfo = &pack.ImageInfo{
					StackID:    "test.stack.id.remote",
					Buildpacks: buildpacks,
					Base:       pack.ImageBase{},
				}
				localInfo = &pack.ImageInfo{
					StackID:    "test.stack.id.local",
					Buildpacks: buildpacks,
				}
			)

			when("an image arg is passed", func() {
				it.Before(func() {
					command.SetArgs([]string{"some/image"})
					mockClient.EXPECT().InspectImage("some/image", false).Return(remoteInfo, nil)
					mockClient.EXPECT().InspectImage("some/image", true).Return(localInfo, nil)
				})

				it("displays image information for local and remote", func() {
					h.AssertNil(t, command.Execute())
					h.AssertContains(t, outBuf.String(), "Inspecting image: 'some/image'")
					h.AssertContains(t, outBuf.String(), `
REMOTE:

Stack: test.stack.id.remote

Buildpacks:
  ID                 VERSION
  test.bp.one        1.0.0
  test.bp.two        2.0.0
`)

					h.AssertContains(t, outBuf.String(), `
LOCAL:

Stack: test.stack.id.local

Buildpacks:
  ID                 VERSION
  test.bp.one        1.0.0
  test.bp.two        2.0.0
`)
				})
			})
		})
	})
}
