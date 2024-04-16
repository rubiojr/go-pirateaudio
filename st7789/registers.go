package st7789

type Rotation uint8

const (
	// REGISTERS

	// System function comand table 1
	NOP       = 0x00
	SWRESET   = 0x01
	RDDID     = 0x04
	RDDST     = 0x09
	RDDPM     = 0x0A
	RDDMADCTL = 0x0B
	RDDCOLMOD = 0x0C
	DDDIM     = 0x0D
	RDDSM     = 0x0E
	RDDSDR    = 0x0F
	SLPIN     = 0x10
	SLPOUT    = 0x11
	PTLON     = 0x12
	NORON     = 0x13
	INVOFF    = 0x20
	INVON     = 0x21
	GAMSET    = 0x26
	DISPOFF   = 0x28
	DISPON    = 0x29
	CASET     = 0x2A
	RASET     = 0x2B
	RAMWR     = 0x2C
	RAMRD     = 0x2E
	PTLAR     = 0x30
	VSCRDEF   = 0x33
	TEOFF     = 0x34
	TEON      = 0x35
	MADCTL    = 0x36
	VSCSAD    = 0x37
	IDMOFF    = 0x38
	IDMON     = 0x39
	COLMOD    = 0x3A
	WRMEMC    = 0x3C
	RDMEMC    = 0x3E
	STE       = 0x44
	GSCAN     = 0x45
	WRDISBV   = 0x51
	RDDISBV   = 0x52
	WRTCTRLD  = 0x53
	RDCTRLD   = 0x54
	WRCACE    = 0x55
	RDCABC    = 0x56
	WRCABCMB  = 0x5E
	RDCABCMB  = 0x5F
	RDABCSDR  = 0x68
	RDID1     = 0xDA
	RDID2     = 0xDB
	RDID3     = 0xDC

	// System function comand table 1
	RAMCTRL   = 0xB0
	RGBCTRL   = 0xB1
	PORCTRL   = 0xB2
	FRCTRL1   = 0xB3
	PARCTRL   = 0xB5
	GCTRL     = 0xB7
	GTADJ     = 0xB8
	DGMEM     = 0xBA
	VCOMS     = 0xBB
	LCMCTRL   = 0xC0
	IDSET     = 0xC1
	VDVVRHEN  = 0xC2
	VRHS      = 0xC3
	VDVS      = 0xC4
	VCMOFSET  = 0xC5
	FRCTRL2   = 0xC6
	CABCCTRL1 = 0xC7
	REGSEL1   = 0xC8
	REGSEL2   = 0xCA
	PWMFRSEL  = 0xCC
	PWCTRL1   = 0xD0
	VAPVANEN  = 0xD2
	CMD2EN    = 0xDF
	PVGAMCTRL = 0xE0
	NVGAMCTRL = 0xE1
	DGMLUTR   = 0xE2
	DGMLUTB   = 0xE3
	GATECTRL  = 0xE4
	SPI2EN    = 0xE7
	PWCTRL2   = 0xE8
	EQCTRL    = 0xE9
	PROMCTRL  = 0xEC
	PROMEN    = 0xFA
	NVMSET    = 0xFC
	PROMACT   = 0xFE

	// COMMAND PARAMETERS

	BG_SPI_CS_BACK  = 0
	BG_SPI_CS_FRONT = 1

	// VDVVRHEN bits
	VDVVRHEN_CMDEN_NVM   = 0x00 // VDV and VRH register value comes from NVM
	VDVVRHEN_CMDEN_WRITE = 0x01 // VDV and VRH register value comes from command write

	// Color mode bits
	COLMOD_RGB_65K   = 0x50 // 16-bit color
	COLMOD_RGB_262K  = 0x60 // 18-bit color
	COLMOD_CTRL_4K   = 0x03 // 12-bit color
	COLMOD_CTRL_65K  = 0x05 // 16-bit color
	COLMOD_CTRL_262K = 0x06 // 18-bit color
	COLMOD_CTRL_16M  = 0x07 // truncates to 18 bits

	// Allowable frame rate codes for FRCTRL2 (Identifier is in Hz)
	FRAMERATE_119 = 0x00
	FRAMERATE_111 = 0x01
	FRAMERATE_105 = 0x02
	FRAMERATE_99  = 0x03
	FRAMERATE_94  = 0x04
	FRAMERATE_90  = 0x05
	FRAMERATE_86  = 0x06
	FRAMERATE_82  = 0x07
	FRAMERATE_78  = 0x08
	FRAMERATE_75  = 0x09
	FRAMERATE_72  = 0x0A
	FRAMERATE_69  = 0x0B
	FRAMERATE_67  = 0x0C
	FRAMERATE_64  = 0x0D
	FRAMERATE_62  = 0x0E
	FRAMERATE_60  = 0x0F
	FRAMERATE_58  = 0x10
	FRAMERATE_57  = 0x11
	FRAMERATE_55  = 0x12
	FRAMERATE_53  = 0x13
	FRAMERATE_52  = 0x14
	FRAMERATE_50  = 0x15
	FRAMERATE_49  = 0x16
	FRAMERATE_48  = 0x17
	FRAMERATE_46  = 0x18
	FRAMERATE_45  = 0x19
	FRAMERATE_44  = 0x1A
	FRAMERATE_43  = 0x1B
	FRAMERATE_42  = 0x1C
	FRAMERATE_41  = 0x1D
	FRAMERATE_40  = 0x1E
	FRAMERATE_39  = 0x1F

	// LCMCTRL bits
	LCMCTRL_XMY  = 0x40 // XOR MY setting in MADCTL
	LCMCTRL_XBGR = 0x20 // XOR RGB setting in MADCTL
	LCMCTRL_XREV = 0x10 // XOR inverse setting in INVON
	LCMCTRL_XMH  = 0x08 // Reverse source output order and only support RGB interface without RAM mode
	LMCTRL_XMV   = 0x04 // XOR MV setting in MADCTL
	LMCTRL_XMX   = 0x02 // XOR MX setting in MADCTL
	LMCTRL_XGS   = 0x01 // XOR GS setting in GATECTRL

	// MADCTL bits
	MADCTL_MY_TB   = 0x00 // Page address order top to bottom
	MADCTL_MY_BT   = 0x80 // Page address order bottom to top
	MADCTL_MX_LR   = 0x00 // Column address order left to right
	MADCTL_MX_RL   = 0x40 // Column address order right to left
	MADCTL_MV_NORM = 0x00 // Page/column order normal
	MADCTL_MV_REV  = 0x20 // Page/column order reverse
	MADCTL_ML_TB   = 0x00 // Line address order LCD refresh top to bottom
	MADCTL_ML_BT   = 0x10 // Line address order LCD refresh bottom to top
	MADCTL_RGB     = 0x00 // RGB order
	MADCTL_BGR     = 0x08 // BGR order
	MADCTL_MH_LR   = 0x00 // Display latch order LCD refresh left to right
	MADCTL_MH_RL   = 0x04 // Display latch order LCD refresh right to left

	MAX_VSYNC_SCANLINES = 254

	ROTATION_NONE Rotation = 0
	ROTATION_90   Rotation = 1
	ROTATION_180  Rotation = 2
	ROTATION_270  Rotation = 3

	SPI_CLOCK_HZ = 16000000
)

func defaultGateControl() uint8 {
	return 0x14 // VGH (V) = 12.54, VGL (V) = -9.6
}

func defaultNegativeGammaCtrl() []byte {
	return []byte{0xD0, 0x04, 0x0C, 0x11, 0x13, 0x2C, 0x3F, 0x44, 0x51, 0x2F, 0x1F, 0x1F, 0x20, 0x23}
}

func defaultPorchControl() []byte {
	return []byte{0x0C, 0x0C, 0x00, 0x33, 0x33}
}

func defaultPositiveGammaCtrl() []byte {
	return []byte{0xD0, 0x04, 0x0D, 0x11, 0x13, 0x2B, 0x3F, 0x54, 0x4C, 0x18, 0x0D, 0x0B, 0x1F, 0x23}
}

func defaultPowerCtrl() []byte {
	return []byte{0xA4, 0xA1}
}

func defaulVCOMSOffsetSet() uint8 {
	return 0x37 // VCOMS OFFSET(V) = 0.575
}

func defaultVDVSet() uint8 {
	return 0x20 // VDV (V) = 0
}

func defaultVRHSet() uint8 {
	return 0x12 // VAP(GVDD) (V) = 4.45+(vcom+vcom offset+0.5vdv)
}

func verticalScrollOffset(offset int) []byte {
	return []byte{0x00, uint8(offset)}
}
