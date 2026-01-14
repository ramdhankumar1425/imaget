package model

type TransformParams interface {
	isTransformParams()
}

type BlurParams struct {
	Sigma float64
}

type SharpenParams struct {
	Sigma float64
}

type ResizeParams struct {
	Width  int
	Height int
}

type CropParams struct {
	Width  int
	Height int
	Anchor string
}

type GrayscaleParams struct {
}

func (BlurParams) isTransformParams()      {}
func (SharpenParams) isTransformParams()   {}
func (ResizeParams) isTransformParams()    {}
func (CropParams) isTransformParams()      {}
func (GrayscaleParams) isTransformParams() {}
