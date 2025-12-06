package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color scheme
type Theme struct {
	Name        string
	Primary     lipgloss.TerminalColor
	Accent      lipgloss.TerminalColor
	Error       lipgloss.TerminalColor
	Text        lipgloss.TerminalColor
	Muted       lipgloss.TerminalColor
	Dim         lipgloss.TerminalColor
	Border      lipgloss.TerminalColor
	PingingWarn lipgloss.TerminalColor
}

// Available themes - optimized for both light and dark backgrounds via AdaptiveColor
var themes = map[string]Theme{
	"purple": {
		Name:        "Purple Dream",
		Primary:     lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}, // Purple
		Accent:      lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#60A5FA"}, // Blue
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"}, // Dark Gray / Light Gray
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}, // Medium Gray / Lighter Gray
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"}, // Light Gray / Darker Gray
		Border:      lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#8B5CF6"}, // Purple
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
	},
	"blue": {
		Name:        "Ocean Blue",
		Primary:     lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#60A5FA"}, // Blue
		Accent:      lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#22D3EE"}, // Cyan
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"},
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"},
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"},
		Border:      lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#3B82F6"}, // Blue
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
	},
	"green": {
		Name:        "Matrix Green",
		Primary:     lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}, // Green
		Accent:      lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#22D3EE"}, // Cyan
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"},
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"},
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"},
		Border:      lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"}, // Green
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
	},
	"pink": {
		Name:        "Bubblegum Pink",
		Primary:     lipgloss.AdaptiveColor{Light: "#DB2777", Dark: "#F472B6"}, // Pink
		Accent:      lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}, // Purple
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"},
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"},
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"},
		Border:      lipgloss.AdaptiveColor{Light: "#DB2777", Dark: "#EC4899"}, // Pink
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
	},
	"amber": {
		Name:        "Sunset Amber",
		Primary:     lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
		Accent:      lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"},
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"},
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"},
		Border:      lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#F59E0B"}, // Amber
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#F59E0B"}, // Amber
	},
	"cyan": {
		Name:        "Cyber Cyan",
		Primary:     lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#22D3EE"}, // Cyan
		Accent:      lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}, // Purple
		Error:       lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}, // Red
		Text:        lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F3F4F6"},
		Muted:       lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"},
		Dim:         lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"},
		Border:      lipgloss.AdaptiveColor{Light: "#0891B2", Dark: "#06B6D4"}, // Cyan
		PingingWarn: lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}, // Amber
	},
}

var currentTheme = themes["purple"]

var (
	// Minimal color palette
	primaryColor = currentTheme.Primary
	accentColor  = currentTheme.Accent
	errorColor   = currentTheme.Error
	textColor    = currentTheme.Text
	mutedColor   = currentTheme.Muted
	dimColor     = currentTheme.Dim
	borderColor  = currentTheme.Border

	// Clean title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(1)

	// Form label style
	labelStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Width(10)

	labelFocusedStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Width(10)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			MarginTop(1)

	// Minimal box style for forms
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2)

	// Focused item style
	focusedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Instructions style
	instructionsStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Padding(1, 0)

	// Key binding styles
	keyStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	descStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// Status indicator styles (text-based) - consistent across all themes
	statusOnlineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"}) // Green

	statusOfflineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}) // Red

	statusUnknownStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#4B5563", Dark: "#9CA3AF"}) // Gray

	statusPingingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"}) // Yellow/Amber
)

// ApplyTheme updates all styles with the selected theme
func ApplyTheme(themeName string) {
	theme, exists := themes[themeName]
	if !exists {
		theme = themes["purple"] // Default fallback
	}

	currentTheme = theme

	// Update color variables
	primaryColor = theme.Primary
	accentColor = theme.Accent
	errorColor = theme.Error
	textColor = theme.Text
	mutedColor = theme.Muted
	dimColor = theme.Dim
	borderColor = theme.Border

	// Update all styles
	titleStyle = titleStyle.Foreground(primaryColor)
	subtitleStyle = subtitleStyle.Foreground(mutedColor)
	labelStyle = labelStyle.Foreground(textColor)
	labelFocusedStyle = labelFocusedStyle.Foreground(primaryColor)
	helpStyle = helpStyle.Foreground(dimColor)
	boxStyle = boxStyle.BorderForeground(borderColor)
	focusedStyle = focusedStyle.Foreground(primaryColor)
	instructionsStyle = instructionsStyle.Foreground(dimColor)
	keyStyle = keyStyle.Foreground(primaryColor)
	descStyle = descStyle.Foreground(dimColor)
	// Status indicators remain consistent across themes
	statusOnlineStyle = statusOnlineStyle.Foreground(lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"})
	statusOfflineStyle = statusOfflineStyle.Foreground(lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"})
	statusUnknownStyle = statusUnknownStyle.Foreground(lipgloss.AdaptiveColor{Light: "#4B5563", Dark: "#9CA3AF"})
	statusPingingStyle = statusPingingStyle.Foreground(lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FBBF24"})
}

// GetThemeNames returns a list of available theme names
func GetThemeNames() []string {
	return []string{"purple", "blue", "green", "pink", "amber", "cyan"}
}

// GetCurrentTheme returns the current theme
func GetCurrentTheme() Theme {
	return currentTheme
}
