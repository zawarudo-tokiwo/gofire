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
}

var styles []lipgloss.Style

// TODO: Add more palettes
var palettes = map[string][]string{
	"red": {
		"#070707", "#1f0707", "#2f0907", "#470907", "#570f07", "#671707", "#771707",
		"#8f2707", "#9f2f07", "#af3f07", "#bf4707", "#c74707", "#df4f07", "#df5707",
		"#e75f07", "#ef6707", "#f76f07", "#f7770f", "#ff7f0f", "#ff8717", "#ff8f17",
		"#ff971f", "#ff9f1f", "#ffa727", "#ffaf27", "#ffb72f", "#ffbf2f", "#ffc737",
		"#ffcf37", "#ffd73f", "#ffdf3f", "#ffe747", "#ffef4f", "#fff75f", "#ffff7f",
		"#ffffaf", "#ffffff",
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
		case "0", "space":
			m.wind = 0
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

			char := " "
			if heat > 0 {
				char = m.config.Char
			}
			s.WriteString(styles[heat].Render(char))
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

			decay := int(m.rnd.Float64() * 6.0) // Cool pixel by a value from 0 to 6

			randomFlicker := int(m.rnd.Float64()*3.0) - 1
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
	paletteFlag := flag.String("palette", "red", "Color palette: red (will add more in the future)")

	flag.Parse()

	colors, ok := palettes[*paletteFlag]
	if !ok {
		fmt.Printf("Unknown palette '%s'. Available: red", *paletteFlag)
		os.Exit(1)
	}

	styles = make([]lipgloss.Style, len(colors))
	for i, c := range colors {
		styles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(c))
	}

	cfg := Config{
		Char:      *charFlag,
		Palette:   *paletteFlag,
		TickSpeed: *speedFlag,
	}

	p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
