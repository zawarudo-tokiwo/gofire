# 🔥 gofire
> **Fire simulation for the terminal.**
![gofire_preview_gif](https://media.giphy.com/media/ToNeUURXHkED93Q1bF/giphy.gif)

`gofire` is an implementation of the **Doom Fire Algorithm** for terminal, written in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for style definitions.
## Features

* 🎨 **Palette Support:** Includes Red, Blue, Green, Purple, Gray, and a **TTY** mode.
* 🌬️ **Interactive Physics:** Control wind direction, fire height, and tick speed.
* 🖥️ **Adaptive Rendering:** Automatically resizes and handles window changes.
## Installation

### Requirements

* **Go 1.22+** (For the `range integer` syntax used in the code).

### Build from Source

```bash
git clone https://github.com/zawarudo-tokiwo/gofire.git
cd gofire
go build -o gofire main.go
```

## Usage

Simply run the command to start the fire:

```bash
gofire
```

### Command Line Flags

You can customize the simulation at startup using flags:

| Flag          | Default | Description                                                             |
| ------------- | ------- | ----------------------------------------------------------------------- |
| `-palette`    | `red`   | Choose color scheme (`red`, `green`, `blue`, `purple`, `gray`, `tty` ). |
| `-char`       | `█`     | The character used to render pixels (try `*`, `#`, or `@`).             |
| `-speed`      | `50ms`  | The update tick rate. Lower is faster (e.g., `-speed 20ms`).            |
| `-decay`      | `6.0`   | How fast the fire cools. Higher = shorter flames.                       |
| `-no-flicker` | `false` | Disable the horizontal wind flicker effect.                             |

**Example:**

```bash
# A fast, blue fire using '#' characters
gofire -palette blue -speed 30ms -char "#"
```
## Controls

You can control the fire in real-time using your keyboard:

| Key                | Action                                          |
| ------------------ | ----------------------------------------------- |
| `h` / `Left`   | Blow wind to the **Left**                       |
| `l` / `Right`  | Blow wind to the **Right**                      |
| `k` / `Up`     | **Decrease Decay** (Make flames taller/hotter)  |
| `j` / `Down`   | **Increase Decay** (Make flames shorter/cooler) |
| `[`            | **Increase Speed** (Run simulation faster)      |
| `]`            | **Decrease Speed** (Run simulation slower)      |
| `f`            | Toggle **Flicker** effect                       |
| `Space` / `0`  | **Reset** wind and decay settings               |
| `q` / `Ctrl+C` | **Quit**                                        |

## 🎨 Palettes

* **Red / Blue / Green / Purple / Gray:** Standard gradients.
* **TTY:** Uses your terminal's ANSI color codes (0-15).
