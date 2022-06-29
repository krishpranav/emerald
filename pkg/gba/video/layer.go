package video

import (
	"github.com/krishpranav/emerald/pkg/util"
)

type Layers [10]Layer

func (ls *Layers) Len() int { return len(ls) }
func (ls *Layers) Less(i, j int) bool {
	bgi, prioi := ls[i].SortInfo()
	bgj, prioj := ls[j].SortInfo()
	diff := prioj - prioi
	if diff == 0 {
		if bgi && !bgj {
			return true
		} else if !bgi && bgj {
			return false
		}

		return ls[j].Index()-ls[i].Index() <= 0
	}
	return false
}
func (ls *Layers) Swap(i, j int) { ls[i], ls[j] = ls[j], ls[i] }

type Layer interface {
	drawScanline(backing *Backing, start, end uint32)
	Index() int
	SortInfo() (bool, int)
	Enabled() bool
}

type BGLayer struct {
	bg       bool
	index    int
	enabled  bool
	video    *SoftwareRenderer
	vram     *VRAM
	priority int
	mosaic   bool
	color256 uint16

	charBase uint32

	screenBase uint32

	overflow bool

	size uint32

	x uint16
	y uint16

	refx, refy       float64
	dx, dmx, dy, dmy float64

	sx, sy float64

	drawScanlineFunc func(backing *Backing, bg *BGLayer, start, end uint32)
}

func (bgl *BGLayer) pushPixel(layer int, objMap SharedMap, video *SoftwareRenderer, row, x, offset uint32, backing *Backing, mask byte, isTileMode bool) {
	index := uint32(0)
	if isTileMode {
		if bgl.color256 > 0 {
			if x >= 4 {
				x -= 4
			}
			index = (row >> (x << 3)) & 0xff
		} else {
			index = (row >> (x << 2)) & 0xf
		}

		if index == 0 {
			return
		} else if bgl.color256 == 0 {
			index |= uint32(objMap.palette)
		}
	}

	stencil := byte(WRITTEN_MASK)
	oldStencil := backing.stencil[offset]
	blend := video.blendMode
	if video.objwinActive {
		if oldStencil&OBJWIN_MASK > 0 {
			if video.windows[3].enabled[layer] {
				video.setBlendEnabled(layer, video.windows[3].special && video.target1[layer] > 0, blend)
				if video.windows[3].special && video.alphaEnabled {
					mask |= video.target1[layer]
				}
				stencil |= OBJWIN_MASK
			} else {
				return
			}
		} else if video.windows[2].enabled[layer] {
			video.setBlendEnabled(layer, video.windows[2].special && video.target1[layer] > 0, blend)
			if video.windows[2].special && video.alphaEnabled {
				mask |= video.target1[layer]
			}
		} else {
			return
		}
	}

	if (mask&TARGET1_MASK > 0) && (oldStencil&TARGET2_MASK > 0) {
		video.setBlendEnabled(layer, true, 1)
	}

	pixel := uint16(row)
	if isTileMode {
		pixel = video.Palette.accessColor(layer, int(index))
	}

	if mask&TARGET1_MASK > 0 {
		video.setBlendEnabled(layer, blend != 0, blend)
	}

	highPriority := (mask & PRIORITY_MASK) < (oldStencil & PRIORITY_MASK)
	if (mask & PRIORITY_MASK) == (oldStencil & PRIORITY_MASK) {
		highPriority = (mask & BACKGROUND_MASK) > 0
	}

	if (oldStencil & WRITTEN_MASK) == 0 {
		stencil |= mask
	} else if highPriority {
		if (mask&TARGET1_MASK != 0) && (oldStencil&TARGET2_MASK != 0) {
			pixel = video.Palette.mix(video.blendA, pixel, video.blendB, backing.color[offset])
		}
		stencil |= mask & ^byte(TARGET1_MASK)
	} else if (mask & PRIORITY_MASK) > (oldStencil & PRIORITY_MASK) {
		stencil = oldStencil & ^byte(TARGET1_MASK|TARGET2_MASK)
		if (mask&TARGET2_MASK != 0) && (oldStencil&TARGET1_MASK != 0) {
			pixel = video.Palette.mix(video.blendB, pixel, video.blendA, backing.color[offset])
		} else {
			return
		}
	} else {
		return
	}

	if (mask & OBJWIN_MASK) > 0 {
		backing.stencil[offset] |= OBJWIN_MASK

		return
	}
	backing.color[offset] = pixel
	backing.stencil[offset] = stencil
}

func (bgl *BGLayer) drawScanline(backing *Backing, start, end uint32) {
	bgl.drawScanlineFunc(backing, bgl, start, end)
}

func (bgl *BGLayer) Index() int            { return bgl.index }
func (bgl *BGLayer) SortInfo() (bool, int) { return bgl.bg, bgl.priority }
func (bgl *BGLayer) Enabled() bool         { return bgl.enabled }

type ObjLayer struct {
	video           *SoftwareRenderer
	bg              bool
	index, priority int
	enabled         bool
	objwin          uint16
}

func NewObjLayer(video *SoftwareRenderer, index int) *ObjLayer {
	return &ObjLayer{
		video:    video,
		index:    LAYER_OBJ,
		priority: index,
	}
}

func (ol *ObjLayer) drawScanline(backing *Backing, start, end uint32) {
	y := ol.video.Vcount
	if start >= end {
		return
	}

	objs := ol.video.OAM.objs
	for i := 0; i < len(objs); i++ {
		obj := objs[i]
		if obj.disable {
			continue
		}
		if (obj.mode & OBJWIN_MASK) != ol.objwin {
			continue
		}
		if (obj.mode&OBJWIN_MASK) == 0 && (ol.priority != obj.priority) {
			continue
		}

		wrappedY := int32(obj.y)
		if obj.y >= VERTICAL_PIXELS {
			wrappedY = int32(obj.y) - 256
		}

		totalHeight := obj.cachedHeight
		if obj.scalerot {
			totalHeight = obj.cachedHeight << util.BoolToInt(obj.doublesize)
		}

		mosaicY := y
		if obj.mosaic {
			mosaicY = y - (y % ol.video.objMosaicY)
		}

		if wrappedY <= int32(y) && wrappedY+int32(totalHeight) > int32(y) {
			obj.drawScanline(backing, uint32(mosaicY), wrappedY, start, end)
		}
	}
}

func (ol *ObjLayer) Index() int            { return ol.index }
func (ol *ObjLayer) SortInfo() (bool, int) { return ol.bg, ol.priority }
func (ol *ObjLayer) Enabled() bool         { return ol.enabled }

type Backdrop struct {
	video           *SoftwareRenderer
	bg              bool
	index, priority int
	enabled         bool
}

func NewBackdrop(video *SoftwareRenderer) *Backdrop {
	return &Backdrop{
		video:    video,
		bg:       true,
		index:    LAYER_BACKDROP,
		priority: -1,
		enabled:  true,
	}
}

func (b *Backdrop) drawScanline(backing *Backing, start, end uint32) {
	for x := start; x < end; x++ {
		if backing.stencil[x]&WRITTEN_MASK == 0 {
			backing.color[x] = b.video.Palette.accessColor(b.index, 0)
			backing.stencil[x] = WRITTEN_MASK
		} else if backing.stencil[x]&TARGET1_MASK > 0 {
			backing.color[x] = b.video.Palette.mix(b.video.blendB, b.video.Palette.accessColor(b.index, 0), b.video.blendA, backing.color[x])
			backing.stencil[x] = WRITTEN_MASK
		}
	}
}

func (b *Backdrop) Index() int            { return b.index }
func (b *Backdrop) SortInfo() (bool, int) { return b.bg, b.priority }
func (b *Backdrop) Enabled() bool         { return b.enabled }
