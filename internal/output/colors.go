package output

type Color string

const (
	cReset Color = "\x1b[0m"

	cDim    Color = "\x1b[2m"
	cRed    Color = "\x1b[31m"
	cGreen  Color = "\x1b[32m"
	cYellow Color = "\x1b[33m"
	cBlue   Color = "\x1b[34m"
	cMag    Color = "\x1b[35m"
	cCyan   Color = "\x1b[36m"
)

func Wrap(noColor bool, c Color, s string) string {
	if noColor {
		return s
	}
	return string(c) + s + string(cReset)
}

func ColorizeError(noColor bool, s string) string {
	return Wrap(noColor, cRed, s)
}
