package gl_utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// Camera2D a Camera based on an orthogonal projection
type Camera2D struct {
	x                  float32
	y                  float32
	width              float32
	halfWidth          float32
	height             float32
	halfHeight         float32
	zoom               float32
	minZoom            float32
	maxZoom            float32
	centered           bool
	flipVertical       bool
	near               float32
	far                float32
	projectionMatrix   mgl32.Mat4
	inverseMatrix      mgl32.Mat4
	matrixDirty        bool
}

// NewCamera2D sets up an orthogonal projection camera
func NewCamera2D(width int, height int, zoom float32) *Camera2D {
	c := &Camera2D{
		width:      float32(width),
		halfWidth:  float32(width) / 2,
		height:     float32(height),
		halfHeight: float32(height) / 2,
		zoom:       zoom,
		minZoom:    0.01,
		maxZoom:    20,
	}
	c.far = -2
	c.near = 2
	c.matrixDirty = true
	c.rebuildMatrix()

	return c
}

func (c *Camera2D) Width() float32  { return c.width }
func (c *Camera2D) Height() float32 { return c.height }

// ProjectionMatrix returns the projection matrix of the camera
func (c *Camera2D) ProjectionMatrix() *mgl32.Mat4 {
	c.rebuildMatrix()
	return &c.projectionMatrix
}

// SetPosition sets the current position of the camera. If the camera is centered, the center will be moving
func (c *Camera2D) SetPosition(x float32, y float32) {
	c.x = x
	c.y = y
	c.matrixDirty = true
}

// Translate move the camera position by the specified amount
func (c *Camera2D) Translate(x float32, y float32) {
	if c.flipVertical {
		y = -y
	}
	c.x += x
	c.y += y
	c.matrixDirty = true
}

// Zoom returns the current zoom level
func (c *Camera2D) Zoom() float32 { return c.zoom }

// SetZoom sets the zoom factor
func (c *Camera2D) SetZoom(zoom float32) {
	zoom = mgl32.Clamp(zoom, c.minZoom, c.maxZoom)
	c.zoom = zoom
	c.matrixDirty = true
}

// MinZoom returns the minimum zoom level allowed
func (c *Camera2D) MinZoom() float32 { return c.minZoom }

// MaxZoom returns the maximum zoom level allowed
func (c *Camera2D) MaxZoom() float32 { return c.maxZoom }

// SetZoomRange sets the minimum and maximum zoom factors allowed
func (c *Camera2D) SetZoomRange(minZoom float32, maxZoom float32) {
	c.minZoom = minZoom
	c.maxZoom = maxZoom
	if c.zoom > c.maxZoom || c.zoom < c.minZoom {
		c.SetZoom(c.zoom)
	}
}

// SetCentered sets the center of the camera to the center of the screen
func (c *Camera2D) SetCentered(centered bool) {
	c.centered = centered
	c.matrixDirty = true
}

// SetFlipVertical sets the orientation of the vertical axis. Pass true to have a cartesian coordinate system
func (c *Camera2D) SetFlipVertical(flip bool) {
	c.flipVertical = flip
	c.matrixDirty = true
}

// SetVisibleArea configures the camera to make the specified area completely visible, position and zoom are changed accordingly
func (c *Camera2D) SetVisibleArea(x1 float32, y1 float32, x2 float32, y2 float32) {
	width := math.Abs(float64(x2 - x1))
	height := math.Abs(float64(y2 - y1))
	zoom := float32(math.Min(float64(c.width)/width, float64(c.height)/height))
	c.SetZoom(zoom)

	x := math.Min(float64(x1), float64(x2))
	y := math.Min(float64(y1), float64(y2))
	if c.centered {
		c.SetPosition(float32(x+width/2), float32(y+height/2))
	} else {
		c.SetPosition(float32(x), float32(y))
	}
}

func (c *Camera2D) rebuildMatrix() {
	if !c.matrixDirty {
		return
	}
	var left, right, top, bottom float32

	if c.centered {
		halfWidth := c.halfWidth / c.zoom
		halfHeight := c.halfHeight / c.zoom
		left = -halfWidth
		right = halfWidth
		top = halfHeight
		bottom = -halfHeight
	} else {
		right = c.width / c.zoom
		top = c.height / c.zoom
	}

	left += c.x
	right += c.x
	top += c.y
	bottom += c.y

	if c.flipVertical {
		bottom, top = top, bottom
	}

	c.projectionMatrix = mgl32.Ortho(left, right, top, bottom, c.near, c.far)
	c.inverseMatrix = c.projectionMatrix.Inv()
	c.matrixDirty = false
}

func (c *Camera2D) ScreenToWorld(vec mgl32.Vec2) mgl32.Vec3 {
	if c.flipVertical {
		vec[1] = c.height - vec[1]
	}
	x := (vec[0] - c.halfWidth) / c.halfWidth
	y := (vec[1] - c.halfHeight) / c.halfHeight
	return mgl32.TransformCoordinate(mgl32.Vec3{x, y, 0}, c.inverseMatrix)
}

func (c *Camera2D) WorldToScreen(vec mgl32.Vec3) mgl32.Vec2 {
	ret := mgl32.TransformCoordinate(vec, c.projectionMatrix)
	ret[0] = ret[0]*c.halfWidth + c.halfWidth
	ret[1] = ret[1]*c.halfHeight + c.halfHeight
	if c.flipVertical {
		ret[1] = c.height - ret[1]
	}
	return mgl32.Vec2{ret[0], ret[1]}
}
