package st7789

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

// DefaultOpts is the recommended default options.
var DefaultOpts = Opts{
	Width:    240,
	Height:   240,
	Rotation: ROTATION_NONE,
}

// Opts defines the options for the device.
type Opts struct {
	Width    int16
	Height   int16
	Rotation Rotation
}

func NewSPI(port spi.Port, dataComm gpio.PinOut, opts *Opts) (*Device, error) {
	if dataComm == gpio.INVALID {
		return nil, errors.New("ssd1306: use nil for dc to use 3-wire mode, do not use gpio.INVALID")
	}
	bits := 8
	if err := dataComm.Out(gpio.Low); err != nil {
		return nil, err
	}
	conn, err := port.Connect(80*physic.MegaHertz, spi.Mode0, bits)
	if err != nil {
		return nil, err
	}

	pin := gpioreg.ByName("GPIO13")
	if err = pin.Out(gpio.Low); err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	if err = pin.Out(gpio.High); err != nil {
		panic(err)
	}

	return newST7789Device(conn, opts, dataComm)
}

type Device struct {
	conn     conn.Conn
	dataComm gpio.PinOut
	rect     image.Rectangle

	rotation                      Rotation
	width                         int16
	height                        int16
	rowOffsetCfg, rowOffset       int16
	columnOffset, columnOffsetCfg int16
	isBGR                         bool
	batchLength                   int32
	backlight                     gpio.PinIO
}

func (d *Device) String() string {
	return fmt.Sprintf("st7789.Device{%s, %s, %s}", d.conn, d.dataComm, d.rect.Max)
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *Device) Bounds() image.Rectangle {
	return d.rect
}

// PowerOff the display
func (d *Device) PowerOff() error {
	return d.backlight.Out(gpio.Low)
}

// PowerOn the display
func (d *Device) PowerOn() error {
	return d.backlight.Out(gpio.High)
}

// Invert the display (black on white vs white on black).
func (d *Device) Invert(blackOnWhite bool) {
	b := byte(0xA6)
	if blackOnWhite {
		b = 0xA7
	}
	d.Command(b)
}

func newST7789Device(conn conn.Conn, opts *Opts, dataComm gpio.PinOut) (*Device, error) {
	d := &Device{
		conn:        conn,
		dataComm:    dataComm,
		rect:        image.Rect(0, 0, int(opts.Width), int(opts.Height)),
		rotation:    opts.Rotation,
		width:       opts.Width,
		height:      opts.Height,
		batchLength: int32(opts.Width),
		backlight:   gpioreg.ByName("GPIO13"),
	}
	d.batchLength = d.batchLength & 1

	d.Command(SWRESET)
	time.Sleep(150 * time.Millisecond)

	d.Command(MADCTL)
	d.Data(MADCTL_MX_RL | MADCTL_MV_REV | MADCTL_ML_BT)

	d.Command(PORCTRL)
	d.SendData(defaultPorchControl())

	d.Command(COLMOD)
	d.Data(COLMOD_CTRL_65K)

	d.Command(GCTRL)
	d.Data(defaultGateControl())

	d.Command(VCOMS)
	d.Data(defaulVCOMSOffsetSet())

	d.Command(LCMCTRL)
	d.Data(LCMCTRL_XBGR | LCMCTRL_XMH | LMCTRL_XMV)

	d.Command(VDVVRHEN)
	d.Data(VDVVRHEN_CMDEN_WRITE)

	d.Command(VRHS)
	d.Data(defaultVRHSet())

	d.Command(VDVS)
	d.Data(defaultVDVSet())

	d.Command(PWCTRL1)
	d.SendData(defaultPowerCtrl())

	d.Command(FRCTRL2)
	d.Data(FRAMERATE_60)

	d.Command(PVGAMCTRL)
	d.SendData(defaultPositiveGammaCtrl())

	d.Command(NVGAMCTRL)
	d.SendData(defaultNegativeGammaCtrl())

	d.Command(INVON)

	d.Command(SLPOUT)

	d.Command(DISPON)

	return d, nil
}

func (d *Device) SetWindow() {
	x1 := d.width - 1
	y1 := d.height - 1
	y0 := 0
	x0 := 0

	d.Command(CASET)
	d.Data(byte(x0 >> 8))
	d.Data(byte(x0 & 0xFF))
	d.Data(byte(x1 >> 8))
	d.Data(byte(x1 & 0xFF))

	d.Command(RASET)
	d.Data(byte(y0 >> 8))
	d.Data(byte(y0 & 0xFF))
	d.Data(byte(y1 >> 8))
	d.Data(byte(y1 & 0xFF))

	d.Command(RAMWR)
	d.Data(0x89)
}

func (d *Device) SendData(c []byte) error {
	if err := d.dataComm.Out(gpio.High); err != nil {
		return err
	}
	return d.conn.Tx(c, nil)
}

func (d *Device) SendCommand(c []byte) error {
	if err := d.dataComm.Out(gpio.Low); err != nil {
		return err
	}
	return d.conn.Tx(c, nil)
}

// FillRectangle fills a rectangle at a given coordinates with a color
func (d *Device) FillRectangle(x, y, width, height int16, c color.RGBA) error {
	k, i := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= k || (x+width) > k || y >= i || (y+height) > i {
		return errors.New("rectangle coordinates outside display area")
	}
	d.SetWindow()
	c565 := RGBATo565(c)
	c1 := uint8(c565)
	c2 := uint8(c565 >> 8)

	data := make([]uint8, d.PixelCount())
	for i := int32(0); i < int32(d.width); i++ {
		data[i*2] = c1
		data[i*2+1] = c2
	}
	j := int32(width) * int32(height)
	for j > 0 {
		if j >= int32(d.height) {
			d.SendData(data)
		} else {
			d.SendData(data[:j*2])
		}
		j -= int32(d.height)
	}
	return nil
}

// Size returns the current size of the display.
func (d *Device) Size() (int16, int16) {
	if d.rotation == ROTATION_NONE || d.rotation == ROTATION_180 {
		return d.width, d.height
	}
	return d.height, d.width
}

// PixelCount returns the number of pixels in the display
func (d *Device) PixelCount() uint32 {
	return uint32(d.width) * uint32(d.height)
}

// RGBATo565 converts a color.RGBA to uint16 used in the display (bits r:5, g:6, b:5)
func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

// SetPixel sets a pixel in the screen
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	if x < 0 || y < 0 ||
		(((d.rotation == ROTATION_NONE || d.rotation == ROTATION_180) && (x >= d.width || y >= d.height)) ||
			((d.rotation == ROTATION_90 || d.rotation == ROTATION_270) && (x >= d.height || y >= d.width))) {
		return
	}
	d.FillRectangle(x, y, 1, 1, c)
}

// FillScreen fills the screen with a given color
func (d *Device) FillScreen(c color.RGBA) {
	if d.rotation == ROTATION_NONE || d.rotation == ROTATION_180 {
		d.FillRectangle(0, 0, d.width, d.height, c)
	} else {
		d.FillRectangle(0, 0, d.height, d.width, c)
	}
}

// SetRotation changes the rotation of the device (clock-wise)
func (d *Device) SetRotation(rotation Rotation) {
	madctl := uint8(0)
	vscsad := verticalScrollOffset(0)
	switch rotation % 4 {
	case ROTATION_NONE:
		madctl = MADCTL_MX_RL | MADCTL_MY_TB | MADCTL_MV_REV
		d.rowOffset = d.rowOffsetCfg
		d.columnOffset = d.columnOffsetCfg
	case ROTATION_90:
		madctl = MADCTL_MX_RL | MADCTL_MY_BT | MADCTL_MV_NORM
		vscsad = verticalScrollOffset(320 - int(d.width))
		d.rowOffset = d.columnOffsetCfg
		d.columnOffset = d.rowOffsetCfg
	case ROTATION_180:
		madctl = MADCTL_MX_LR | MADCTL_MY_BT | MADCTL_MV_REV
		vscsad = verticalScrollOffset(320 - int(d.width))
		d.rowOffset = 0
		d.columnOffset = 0
	case ROTATION_270:
		madctl = MADCTL_MX_LR | MADCTL_MY_TB | MADCTL_MV_NORM
		d.rowOffset = 0
		d.columnOffset = 0
	}
	if d.isBGR {
		madctl |= MADCTL_BGR
	}

	// Set the display orientation
	d.Command(MADCTL)
	d.Data(madctl)

	// Set vertical scroll offset so that images are located correctly on 240x240 displays
	d.Command(VSCSAD)
	d.SendData(vscsad)
}

// IsBGR changes the color mode (RGB/BGR)
func (d *Device) IsBGR(bgr bool) {
	d.isBGR = bgr
}

// InverColors inverts the colors of the screen
func (d *Device) InvertColors(invert bool) {
	if invert {
		d.Command(INVON)
	} else {
		d.Command(INVOFF)
	}
}

// Command sends a command to the device
func (d *Device) Command(cmd uint8) {
	d.SendCommand([]byte{cmd})
}

// Data sends data to the device
func (d *Device) Data(data uint8) {
	d.SendData([]byte{data})
}

// DrawFastVLine draws a vertical line faster than using SetPixel
func (d *Device) DrawFastVLine(x, y0, y1 int16, c color.RGBA) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	d.FillRectangle(x, y0, 1, y1-y0+1, c)
}

// DrawFastHLine draws a horizontal line faster than using SetPixel
func (d *Device) DrawFastHLine(x0, x1, y int16, c color.RGBA) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	d.FillRectangle(x0, y, x1-x0+1, 1, c)
}

func (d *Device) DrawImage(reader io.Reader) {
	d.SetWindow()
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	d.DrawRAW(img)
}

func (d *Device) DrawRAW(img image.Image) {
	d.SetWindow()
	rect := img.Bounds()
	rgbaimg := image.NewRGBA(rect)
	draw.Draw(rgbaimg, rect, img, rect.Min, draw.Src)

	np := []uint8{}
	for i := 0; i < int(d.width); i++ {
		for j := 0; j < int(d.height); j++ {
			rgba := rgbaimg.At(int(d.width)-i, j).(color.RGBA)
			c565 := RGBATo565(rgba)
			c1 := uint8(c565)
			c2 := uint8(c565 >> 8)
			np = append(np, c1, c2)
		}
	}

	for i := 0; i < len(np); i += 4096 {
		d.SendData(np[i : i+4096])
	}
}
