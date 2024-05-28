package pretrained

import (
	"torch/torch"
)

// Model is an interface for pretrained model.
// It has only one method `Load(string) error` to load model
// from local or remote file.
type Model interface {
	Load(modelNamOrPath string, config interface{ Config }, params map[string]interface{}, device torch.Device) error
}
