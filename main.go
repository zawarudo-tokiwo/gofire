package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Char      string
	Palette   string
	TickSpeed time.Duration
	Decay     float64
	Flicker   bool
}

var (
	styles     []lipgloss.Style
	styleCache []string
)

var palettes = map[string][]string{
	"red": {
		"#070707", "#1f0707", "#2f0907", "#470907", "#570f07", "#671707", "#771707",
		"#8f2707", "#9f2f07", "#af3f07", "#bf4707", "#c74707", "#df4f07", "#df5707",
		"#e75f07", "#ef6707", "#f76f07", "#f7770f", "#ff7f0f", "#ff8717", "#ff8f17",
		"#ff971f", "#ff9f1f", "#ffa727", "#ffaf27", "#ffb72f", "#ffbf2f", "#ffc737",
		"#ffcf37", "#ffd73f", "#ffdf3f", "#ffe747", "#ffef4f", "#fff75f", "#ffff7f",
		"#ffffaf", "#ffffff",
	},
	"blue": {
		"#000000", "#000614", "#000b21", "#001336", "#001842", "#001d52", "#002466",
		"#002b7a", "#00328f", "#0039a3", "#0040b8", "#0047cc", "#004de0", "#0054f5",
		"#0059ff", "#004ec7", "#0044ad", "#003a94", "#00307a", "#002661", "#001c47",
		"#1a3c8e", "#335cb5", "#4d7ddd", "#669df4", "#80beff", "#99ceff", "#b3deff",
		"#cceeff", "#e6faff", "#f0fbff", "#f5fdff", "#faffff", "#fbffff", "#fdffff",
		"#feffff", "#ffffff",
	},
	"green": {
		"#000000", "#051405", "#0a210a", "#0f360f", "#144214", "#1a571a", "#216b21",
		"#267a26", "#2d8c2d", "#339e33", "#3aaf3a", "#42c242", "#4bd44b", "#54e654",
		"#5df75d", "#66ff66", "#5ce65c", "#52cc52", "#47b347", "#3d993d", "#338033",
		"#47a347", "#5cc75c", "#70eb70", "#85ff85", "#99ff99", "#adffad", "#c2ffc2",
		"#d6ffd6", "#ebffeb", "#f0fff0", "#f5fff5", "#fafffa", "#fbfffb", "#fdfffd",
		"#fefffe", "#ffffff",
	},
	"gray": {
		"#000000", "#0a0a0a", "#141414", "#1e1e1e", "#282828", "#323232", "#3c3c3c",
		"#464646", "#505050", "#5a5a5a", "#646464", "#6e6e6e", "#787878", "#828282",
		"#8c8c8c", "#969696", "#909090", "#888888", "#808080", "#787878", "#707070",
		"#808080", "#909090", "#a0a0a0", "#b0b0b0", "#c0c0c0", "#d0d0d0", "#dcdcdc",
		"#e6e6e6", "#f0f0f0", "#f5f5f5", "#f8f8f8", "#fafafa", "#fcfcfc", "#fdfdfd",
		"#fefefe", "#ffffff",
	},
}

type model struct {
	width      int
	height     int
	firePixels []int
	wind       int
	rnd        *rand.Rand
	config     Config
}

func initialModel(cfg Config) model {
	return model{
		width:      0,
		height:     0,
		firePixels: []int{},
		wind:       0,
		rnd:        rand.New(rand.NewSource(time.Now().UnixNano())),
		config:     cfg,
	}
}

func tick(speed time.Duration) tea.Cmd {
	return tea.Tick(speed, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick(m.config.TickSpeed)
}

type tickMsg time.Time

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "h", "left":
			m.wind--
		case "l", "right":
			m.wind++

		case "f":
			m.config.Flicker = !m.config.Flicker

		case "k", "up":
			m.config.Decay -= 0.5
			if m.config.Decay < 0 {
				m.config.Decay = 0
			}
		case "j", "down":
			m.config.Decay += 0.5

		case "]":
			m.config.TickSpeed -= 10 * time.Millisecond
			if m.config.TickSpeed < 1*time.Millisecond {
				m.config.TickSpeed = 1 * time.Millisecond
			}
		case "[":
			m.config.TickSpeed += 10 * time.Millisecond

		case "0", " ": // Reset config
			m.wind = 0
			m.config.Decay = 6.0
			m.config.TickSpeed = 40 * time.Millisecond
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.firePixels = make([]int, m.width*m.height)
		m.igniteSource()

	case tickMsg:
		if m.width > 0 && m.height > 0 {
			m.spreadFire()
		}
		return m, tick(m.config.TickSpeed)
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	var s strings.Builder
	for y := 0; y < m.height-1; y++ {
		for x := 0; x < m.width; x++ {
			idx := y*m.width + x
			heat := m.firePixels[idx]

			if heat > 0 {
				s.WriteString(styleCache[heat])
			} else {
				s.WriteString(" ") // No ANSI codes needed for empty space
			}
		}
		if y < m.height-2 {
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func (m *model) igniteSource() {
	if m.height == 0 || m.width == 0 {
		return
	}
	// Last row of the grid
	startIdx := (m.height - 1) * m.width
	for x := 0; x < m.width; x++ {
		m.firePixels[startIdx+x] = 36 // Max temp
	}
}

func (m *model) spreadFire() {
	for x := 0; x < m.width; x++ {
		for y := 1; y < m.height; y++ {
			srcIndex := y*m.width + x
			pixelHeat := m.firePixels[srcIndex]

			// Erase completely cold
			if pixelHeat == 0 {
				targetIndex := (y-1)*m.width + x // Pixel above the cold one
				if targetIndex >= 0 && targetIndex < len(m.firePixels) {
					m.firePixels[targetIndex] = 0
				}
				continue
			}

			decay := int(m.rnd.Float64() * m.config.Decay) // Cool pixel by a value from 0 to 6

			randomFlicker := 0
			if m.config.Flicker {
				randomFlicker = int(m.rnd.Float64()*3.0) - 1
			}
			totalWind := randomFlicker + m.wind
			targetX := x + totalWind

			if targetX < 0 || targetX >= m.width {
				continue
			}

			targetIndex := (y-1)*m.width + targetX // Pixel above + wind adjustment
			newHeat := max(pixelHeat-decay, 0)     // LSP suggested it (heat can't be negative)

			if targetIndex >= 0 && targetIndex < len(m.firePixels) {
				m.firePixels[targetIndex] = newHeat
			}
		}
	}
}

func main() {
	charFlag := flag.String("char", "█", "The character used to draw the fire")
	speedFlag := flag.Duration("speed", 50*time.Millisecond, "Tick speed (e.g. 30ms, 100ms)")
	paletteFlag := flag.String("palette", "red", "Color palette: red, green, blue, gray")
	decayFlag := flag.Float64("decay", 6.0, "Heat decay intensity (higher value, shorter flame)")
	noFlicker := flag.Bool("no-flicker", false, "Disable flicker")

	flag.Parse()

	isFlickerEnabled := !*noFlicker

	colors, ok := palettes[*paletteFlag]
	if !ok {
		var available []string
		for k := range palettes {
			available = append(available, k)
		}
		fmt.Printf("Unknown palette '%s'. Available: %s\n", *paletteFlag, strings.Join(available, ", "))
		os.Exit(1)
	}

	styleCache = make([]string, len(colors))
	for i, c := range colors {
		// Pre-render the character into the style immediately
		styleCache[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(*charFlag)
	}

	cfg := Config{
		Char:      *charFlag,
		Palette:   *paletteFlag,
		TickSpeed: *speedFlag,
		Decay:     *decayFlag,
		Flicker:   isFlickerEnabled,
	}

	p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
