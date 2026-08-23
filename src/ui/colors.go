package ui

type Color struct {
	R, G, B, A uint8
}

func color(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func (c Color) ABGR() uint32 {
	return uint32(c.A)<<24 | uint32(c.B)<<16 | uint32(c.G)<<8 | uint32(c.R)
}

func containsIgnoreCase(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 = c1 + 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 = c2 + 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func getBikeManufacturerColor(bikeName string) Color {
	name := bikeName
	if len(name) > 0 {
		name = bikeName[:min(len(bikeName), 20)]
	}
	switch {
	case containsIgnoreCase(name, "Fantic"):
		return color(32, 0, 128, 255)
	case containsIgnoreCase(name, "GasGas"):
		return color(255, 255, 255, 255)
	case containsIgnoreCase(name, "Honda"):
		return color(255, 0, 0, 255)
	case containsIgnoreCase(name, "Husqvarna"):
		return color(224, 224, 224, 255)
	case containsIgnoreCase(name, "Kawasaki"):
		return color(0, 255, 0, 255)
	case containsIgnoreCase(name, "KTM"):
		return color(255, 128, 0, 255)
	case containsIgnoreCase(name, "Suzuki"):
		return color(255, 255, 0, 255)
	case containsIgnoreCase(name, "TM MX"):
		return color(0, 176, 255, 255)
	case containsIgnoreCase(name, "Triumph"):
		return color(255, 227, 63, 255)
	case containsIgnoreCase(name, "Yamaha"):
		return color(0, 0, 255, 255)
	case containsIgnoreCase(name, "Beta"):
		return color(192, 0, 0, 255)
	default:
		return color(128, 128, 128, 255)
	}
}
