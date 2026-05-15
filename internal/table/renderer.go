package table

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

type Renderer interface {
	RenderContainers(containers []docker.Container, columns []string, log *slog.Logger)
	RenderImages(images []docker.Image, columns []string, log *slog.Logger)
	RenderVolumes(volumes []docker.Volume, columns []string, log *slog.Logger)
}

type renderer struct {
	styles Styles
}

type Styles struct {
	Header  lipgloss.Style
	Cell    lipgloss.Style
	Border  lipgloss.Border
	Running lipgloss.Style
	Exited  lipgloss.Style
	Other   lipgloss.Style
}

func NewRenderer() Renderer {
	return &renderer{
		styles: Styles{
			Header: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")),
			Cell: lipgloss.NewStyle().
				Padding(0, 1),
			Border:  lipgloss.RoundedBorder(),
			Running: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
			Exited:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")),
			Other:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")),
		},
	}
}

func (r *renderer) RenderContainers(containers []docker.Container, columns []string, log *slog.Logger) {
	if len(containers) == 0 {
		fmt.Println("No containers found.")
		return
	}

	width := getTerminalWidth()

	headers := columns
	rows := make([][]string, len(containers))

	for i, c := range containers {
		row := make([]string, len(columns))
		for j, col := range columns {
			row[j] = r.formatCell(col, c)
		}
		rows[i] = row
	}

	t := table.New().
		Border(r.styles.Border).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return r.styles.Header
			}
			return r.styles.Cell
		}).
		Headers(headers...).
		Rows(rows...).
		Width(width)

	fmt.Println(t.Render())
}

func (r *renderer) formatCell(col string, c docker.Container) string {
	switch strings.ToUpper(col) {
	case "ID":
		return c.ID
	case "NAME":
		return c.Name
	case "IMAGE":
		return c.Image
	case "COMMAND":
		return truncate(c.Command, 30)
	case "CREATED":
		return humanizeTime(c.Created)
	case "STATUS":
		return r.colorStatus(c.Status)
	case "PORTS":
		return c.Ports
	default:
		return ""
	}
}

func (r *renderer) colorStatus(status string) string {
	clean := strings.Split(status, " (")[0]
	lower := strings.TrimSpace(strings.ToLower(clean))
	switch {
	case strings.Contains(lower, "up"):
		return r.styles.Running.Render(lower)
	case strings.Contains(lower, "exited"):
		return r.styles.Exited.Render(lower)
	default:
		return r.styles.Other.Render(lower)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func humanizeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d.Hours() < 1:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d.Hours() < 24:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 60 {
		return 120
	}
	return width
}

func (r *renderer) RenderImages(images []docker.Image, columns []string, log *slog.Logger) {
	if len(images) == 0 {
		fmt.Println("No images found.")
		return
	}

	width := getTerminalWidth()
	headers := columns
	rows := make([][]string, len(images))

	for i, img := range images {
		row := make([]string, len(columns))
		for j, col := range columns {
			row[j] = r.formatImageCell(col, img)
		}
		rows[i] = row
	}

	t := table.New().
		Border(r.styles.Border).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return r.styles.Header
			}
			return r.styles.Cell
		}).
		Headers(headers...).
		Rows(rows...).
		Width(width - 4)

	fmt.Println(t.Render())
}

func (r *renderer) formatImageCell(col string, img docker.Image) string {
	switch strings.ToUpper(col) {
	case "REPOSITORY":
		return truncate(img.Repository, 30)
	case "TAG":
		return img.Tag
	case "IMAGE ID":
		return img.ID
	case "CREATED":
		return humanizeTime(img.Created)
	case "SIZE":
		return humanizeSize(img.Size)
	default:
		return ""
	}
}

func humanizeSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (r *renderer) RenderVolumes(volumes []docker.Volume, columns []string, log *slog.Logger) {
	if len(volumes) == 0 {
		fmt.Println("No volumes found.")
		return
	}

	width := getTerminalWidth()
	headers := columns
	rows := make([][]string, len(volumes))

	for i, v := range volumes {
		row := make([]string, len(columns))
		for j, col := range columns {
			row[j] = r.formatVolumeCell(col, v)
		}
		rows[i] = row
	}

	t := table.New().
		Border(r.styles.Border).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return r.styles.Header
			}
			return r.styles.Cell
		}).
		Headers(headers...).
		Rows(rows...).
		Width(width - 4)

	fmt.Println(t.Render())
}

func (r *renderer) formatVolumeCell(col string, v docker.Volume) string {
	switch strings.ToUpper(col) {
	case "NAME":
		return truncate(v.Name, 30)
	case "DRIVER":
		return v.Driver
	case "MOUNTPOINT":
		return truncate(v.Mountpoint, 40)
	case "CREATED":
		return humanizeTime(v.Created)
	default:
		return ""
	}
}
