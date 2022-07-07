package main

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	imagev1API "github.com/openshift/api/image/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var logger logrus.FieldLogger

const (
	Version = "v0.0.2"
	// flags
	SrcInternalRegistryFlag = "src-internal-registry"
	DefaultInternalRegistry = "image-registry.openshift-image-registry.svc:5000"

	removeAnnotationsOp = `[
{ "op": "remove", "path": "/tag/annotations"}
]`
)

func main() {
	logger = logrus.New()
	fields := []transform.OptionalFields{
		{
			FlagName: SrcInternalRegistryFlag,
			Help:     "Internal registry hostname[:port] used to determine whether an istag references a local image",
			Example:  "image-registry.openshift-image-registry.svc:5000",
		},
	}
	cli.RunAndExit(cli.NewCustomPlugin("ImagestreamPlugin", Version, fields, Run))
}

type imagestreamOptionalFields struct {
	SrcInternalRegistry string
}

func getOptionalFields(extras map[string]string) (imagestreamOptionalFields, error) {
	// Use default if not provided
	internalRegistry := DefaultInternalRegistry
	if registry, ok := extras[SrcInternalRegistryFlag]; ok {
		internalRegistry = registry
	}

	return imagestreamOptionalFields{
		SrcInternalRegistry: internalRegistry,
	}, nil
}

func Run(request transform.PluginRequest) (transform.PluginResponse, error) {
	u := request.Unstructured
	var patch jsonpatch.Patch
	whiteOut := false
	inputFields, err := getOptionalFields(request.Extras)
	if err != nil {
		return transform.PluginResponse{}, err
	}

	switch u.GetKind() {
	case "ImageStream":
		logger.Info("found ImageStream, adding to whiteout")
		whiteOut = true
	case "ImageTag":
		logger.Info("found ImageTag, adding to whiteout")
		whiteOut = true
	case "ImageStreamTag":
		logger.Info("found ImageStreamTag, processing")
		whiteOut, patch, err = processISTag(u, inputFields)
	}
	if err != nil {
		return transform.PluginResponse{}, err
	}

	return transform.PluginResponse{
		Version:    string(transform.V1),
		IsWhiteOut: whiteOut,
		Patches:    patch,
	}, nil
}

func processISTag(u unstructured.Unstructured, fields imagestreamOptionalFields) (bool, jsonpatch.Patch, error) {
	patch := jsonpatch.Patch{}
	js, err := u.MarshalJSON()
	if err != nil {
		return false, patch, err
	}

	imageStreamTag := &imagev1API.ImageStreamTag{}

	err = json.Unmarshal(js, imageStreamTag)
	if err != nil {
		return false, patch, err
	}

	annotations := imageStreamTag.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}

	dockerImageReference := imageStreamTag.Image.DockerImageReference
	localImage := len(fields.SrcInternalRegistry) > 0 && HasImageRefPrefix(dockerImageReference, fields.SrcInternalRegistry)
	referenceTag := imageStreamTag.Tag != nil && imageStreamTag.Tag.From != nil
	if referenceTag {
		// Removing annotations from the tag, to prevent mismatch
		if imageStreamTag.Tag.Annotations != nil {
			patchJSON := fmt.Sprintf(removeAnnotationsOp)
			patch, err = jsonpatch.DecodePatch([]byte(patchJSON))
			if err != nil {
				return false, patch, err
			}
		}
		if imageStreamTag.Tag.From.Kind == "ImageStreamImage" {
			if imageStreamTag.Tag.From.Namespace == "" || imageStreamTag.Tag.From.Namespace == imageStreamTag.Namespace {
				referenceTag = false
			}
		}
	}

	// Restore the tag if this is a reference tag *or* an external image. Otherwise,
	// image import will create the imagestreamtag automatically.
	if referenceTag || !localImage {
		return false, patch, nil
	}
	// It's a local non-reference tag
	return true, jsonpatch.Patch{}, nil
}

// HasImageRefPrefix returns true if the input image reference begins with
// the input prefix followed by "/"
func HasImageRefPrefix(s, prefix string) bool {
	refSplit := strings.SplitN(s, "/", 2)
	if len(refSplit) != 2 {
		return false
	}
	return refSplit[0] == prefix
}
