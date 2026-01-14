package infra

import (
	"os"

	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

var ImageKit imagekit.Client

func InitImageKit() {
	key := os.Getenv("IMAGEKIT_PRIVATE_KEY")

	ImageKit = imagekit.NewClient(option.WithPrivateKey(key))
}
