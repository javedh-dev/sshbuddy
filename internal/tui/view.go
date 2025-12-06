package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// View renders the main view based on current state
func (m Model) View() string {
	// Fixed width box for 2-column layout
	const boxWidth = 80
	const minHeight = 24

	// Check if terminal is too small
	if m.width < boxWidth+4 || m.height < minHeight {
		errorMsg := lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Align(lipgloss.Center).
			Render("⚠ Terminal Too Small ⚠")

		instruction := lipgloss.NewStyle().
			Foreground(mutedColor).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("Please resize your terminal to at least %dx%d", boxWidth+4, minHeight))

		currentSize := lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("Current: %dx%d", m.width, m.height))

		errorBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(2, 4).
			Render(lipgloss.JoinVertical(lipgloss.Center, errorMsg, "", instruction, "", currentSize))

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, errorBox)
	}

	if m.state == stateForm {
		// Form view (handles its own styling and centering)
		return m.form.View()
	}

	if m.state == stateConfig {
		// Configuration view
		return m.configView.View()
	}

	if m.state == stateTermixAuth {
		// Termix auth view
		return m.termixAuth.View()
	}

	if m.state == stateConfirmDelete {
		// Confirmation dialog
		return m.renderDeleteConfirmation()
	}

	if m.state == stateConfigError {
		// Config error view
		return m.renderConfigError()
	}

	if m.state == stateSourceSelect {
		// Source selection dialog
		return m.renderSourceSelect()
	}

	// ASCII art header
	asciiArt := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Width(boxWidth - 4).
		Align(lipgloss.Center).
		Render(`╔═╗┌─┐┬ ┬  ╔╗ ┬ ┬┌┬┐┌┬┐┬ ┬
╚═╗└─┐├─┤  ╠╩╗│ │ ││ ││└┬┘
╚═╝└─┘┴ ┴  ╚═╝└─┘─┴┘─┴┘ ┴`)

	// Theme indicator
	theme := GetCurrentTheme()
	themeIndicator := lipgloss.NewStyle().
		Foreground(dimColor).
		Width(boxWidth - 4).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Theme: %s", theme.Name))

	separator := lipgloss.NewStyle().
		Foreground(dimColor).
		Width(boxWidth - 4).
		Align(lipgloss.Center).
		Render(strings.Repeat("─", boxWidth-4))

	header := lipgloss.JoinVertical(lipgloss.Left, asciiArt, themeIndicator, separator)

	// Footer with key bindings including ping command and theme switcher
	keyBindings := []string{
		keyStyle.Render("↵") + descStyle.Render(":connect "),
		keyStyle.Render("n") + descStyle.Render(":new "),
		keyStyle.Render("e") + descStyle.Render(":edit "),
		keyStyle.Render("c") + descStyle.Render(":copy "),
		keyStyle.Render("d") + descStyle.Render(":del "),
		keyStyle.Render("f") + descStyle.Render(":fav "),
		keyStyle.Render("p") + descStyle.Render(":ping "),
		keyStyle.Render("s") + descStyle.Render(":settings "),
		keyStyle.Render("/") + descStyle.Render(":search "),
		keyStyle.Render("q") + descStyle.Render(":quit"),
	}
	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(borderColor).
		Width(boxWidth-4).
		Padding(0, 0).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, keyBindings...))

	// Render list in 2 columns
	listView := m.renderTwoColumnList()

	// Add search bar if filtering is active or has filter value
	var searchBar string
	filterState := m.list.FilterState()
	searchQuery := m.list.FilterValue()

	// Show search bar when filtering or when there's a filter value
	if filterState == list.Filtering || filterState == list.FilterApplied {
		if searchQuery == "" {
			searchQuery = "_" // Show cursor when empty
		}
		searchBar = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2).
			Render(fmt.Sprintf("Search: %s", searchQuery))

		searchBar = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(primaryColor).
			Width(boxWidth - 4).
			Render(searchBar)
	}

	// Combine all elements
	var content string
	if searchBar != "" {
		content = lipgloss.JoinVertical(lipgloss.Left,
			header,
			searchBar,
			listView,
			footer,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left,
			header,
			"",
			listView,
			footer,
		)
	}

	// Wrap in a fixed-width box
	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Width(boxWidth).
		Padding(0, 2).
		Render(content)

	// Center the fixed box on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainBox)
}

// renderTwoColumnList renders the host list in a two-column layout
func (m *Model) renderTwoColumnList() string {
	items := m.list.VisibleItems()
	if len(items) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			Padding(2, 0).
			Render("No hosts configured. Press 'n' to add a new host.")
		return emptyMsg
	}

	const columnWidth = 34 // Each column width
	const columnGap = 2    // Gap between columns
	const itemHeight = 3   // Title + Description + Tags
	const listHeight = 3   // Number of items visible per column

	var leftColumn, rightColumn []string

	// Get the current cursor position
	cursor := m.list.Index()
	startIdx := 0

	// Calculate scroll offset to keep cursor visible
	itemsPerScreen := listHeight * 2 // Two columns
	if cursor >= itemsPerScreen {
		startIdx = ((cursor / itemsPerScreen) * itemsPerScreen)
	}

	// Split items into two columns with scrolling
	endIdx := min(startIdx+itemsPerScreen, len(items))

	// Helper function to render an item or empty placeholder
	renderItemAtIndex := func(i int) string {
		// Check if we have an actual item at this position
		if i >= len(items) {
			// Return empty placeholder
			return lipgloss.NewStyle().
				Width(columnWidth).
				Height(itemHeight).
				Render("")
		}

		if itm, ok := items[i].(item); ok {
			isSelected := i == cursor

			// Format the item with status
			var statusText string
			if itm.pinging {
				statusText = statusPingingStyle.Render("●")
			} else {
				switch itm.status {
				case "🟢":
					statusText = statusOnlineStyle.Render("●")
				case "🔴":
					statusText = statusOfflineStyle.Render("●")
				default:
					statusText = statusUnknownStyle.Render("○")
				}
			}

			// Title line - build with alias and ping time
			alias := itm.host.Alias

			// Truncate alias to fit with ping time
			maxAliasLen := 15
			if len(alias) > maxAliasLen {
				alias = alias[:maxAliasLen-3] + "..."
			}

			// Style the alias with primary color
			styledAlias := lipgloss.NewStyle().Foreground(primaryColor).Render(alias)

			pingTimeStr := ""
			if itm.pingTime != "" {
				pingTimeStr = lipgloss.NewStyle().Foreground(dimColor).Render(fmt.Sprintf(" (%s)", itm.pingTime))
			}

			port := itm.host.Port
			if port == "" {
				port = "22"
			}

			// Description line - truncate to fit
			hostInfo := fmt.Sprintf("%s@%s:%s", itm.host.User, itm.host.Hostname, port)
			if len(hostInfo) > 28 {
				hostInfo = hostInfo[:25] + "..."
			}

			// Source line - render with colors and favorite indicator
			sourceLine := renderSource(itm.host.Source, itm.host.AvailableIn, itm.host.Favorite, columnWidth-2, isSelected)

			var titleLine, descLine string
			if isSelected {
				// Selected item with thick border - need to account for border width
				titleLine = lipgloss.NewStyle().
					Bold(true).
					BorderLeft(true).
					BorderStyle(lipgloss.ThickBorder()).
					BorderForeground(primaryColor).
					Padding(0, 0, 0, 1).
					Width(columnWidth - 2). // Subtract border + padding
					Render(fmt.Sprintf("%s %s%s", statusText, styledAlias, pingTimeStr))

				descLine = lipgloss.NewStyle().
					Foreground(mutedColor).
					BorderLeft(true).
					BorderStyle(lipgloss.ThickBorder()).
					BorderForeground(primaryColor).
					Padding(0, 0, 0, 1).
					Width(columnWidth - 2). // Subtract border + padding
					Render(hostInfo)

				sourceLine = lipgloss.NewStyle().
					BorderLeft(true).
					BorderStyle(lipgloss.ThickBorder()).
					BorderForeground(primaryColor).
					Padding(0, 0, 0, 1).
					Width(columnWidth - 2).
					Render(sourceLine)
			} else {
				// Normal item without border - use full width with padding
				titleLine = lipgloss.NewStyle().
					Padding(0, 0, 0, 2).
					Width(columnWidth - 2). // Subtract padding
					Render(fmt.Sprintf("%s %s%s", statusText, styledAlias, pingTimeStr))

				descLine = lipgloss.NewStyle().
					Foreground(dimColor).
					Padding(0, 0, 0, 2).
					Width(columnWidth - 2). // Subtract padding
					Render(hostInfo)

				sourceLine = lipgloss.NewStyle().
					Padding(0, 0, 0, 2).
					Width(columnWidth - 2).
					Render(sourceLine)
			}

			// Wrap in a fixed-width container to prevent shifting
			titleLine = lipgloss.NewStyle().Width(columnWidth).Render(titleLine)
			descLine = lipgloss.NewStyle().Width(columnWidth).Render(descLine)
			sourceLine = lipgloss.NewStyle().Width(columnWidth).Render(sourceLine)

			return lipgloss.JoinVertical(lipgloss.Left, titleLine, descLine, sourceLine)
		}

		return lipgloss.NewStyle().Width(columnWidth).Height(itemHeight).Render("")
	}

	// Render items row-wise: fill left column first, then right column for each row
	// Only render rows that have at least one item
	for row := 0; row < listHeight; row++ {
		leftIdx := startIdx + (row * 2)      // 0, 2, 4, 6...
		rightIdx := startIdx + (row * 2) + 1 // 1, 3, 5, 7...

		// Stop if we've rendered all available items
		if leftIdx >= len(items) {
			break
		}

		leftColumn = append(leftColumn, renderItemAtIndex(leftIdx))

		// Only add right column if there's an item for it
		if rightIdx < len(items) {
			rightColumn = append(rightColumn, renderItemAtIndex(rightIdx))
		} else {
			// Add empty space to maintain layout
			rightColumn = append(rightColumn, lipgloss.NewStyle().
				Width(columnWidth).
				Height(itemHeight).
				Render(""))
		}
	}

	// Join columns side by side with gap
	var rows []string
	for i := 0; i < len(leftColumn); i++ {
		row := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn[i], rightColumn[i])
		row_space := lipgloss.JoinHorizontal(lipgloss.Top, "")
		rows = append(rows, row)
		rows = append(rows, row_space)
	}

	listContent := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Add scroll indicator if needed
	if len(items) > itemsPerScreen {
		scrollInfo := lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			Render(fmt.Sprintf("  %d-%d of %d (↑↓ scroll)", startIdx+1, min(endIdx, len(items)), len(items)))
		listContent = lipgloss.JoinVertical(lipgloss.Left, listContent, scrollInfo)
	}

	return listContent
}

// renderSource renders the source label with icons for all available sources and favorite indicator
func renderSource(primarySource string, availableIn []string, isFavorite bool, maxWidth int, isSelected bool) string {
	if primarySource == "" {
		primarySource = "manual"
	}

	// Helper to get icon for a source
	getIcon := func(source string) string {
		switch source {
		case "manual", "sshbuddy":
			return "◆" // Diamond for manual/sshbuddy
		case "ssh-config":
			return "■" // Square for config file
		case "termix":
			return "▲" // Triangle for API/cloud
		default:
			return "○"
		}
	}

	// Helper to get display name
	getDisplayName := func(source string) string {
		switch source {
		case "manual", "sshbuddy":
			return "sshbuddy"
		case "ssh-config":
			return "config"
		case "termix":
			return "termix"
		default:
			return source
		}
	}

	// Build the source display showing all sources
	var sourceParts []string

	if len(availableIn) == 0 {
		// Fallback to primary source if AvailableIn is empty
		icon := getIcon(primarySource)
		displayName := getDisplayName(primarySource)
		sourceParts = append(sourceParts, icon+" "+displayName)
	} else {
		// Show all sources with icon + name
		for _, src := range availableIn {
			icon := getIcon(src)
			displayName := getDisplayName(src)
			sourceParts = append(sourceParts, icon+" "+displayName)
		}
	}

	// Join all source parts with spacing and apply consistent color
	sourceText := strings.Join(sourceParts, "  ")
	sourceStyle := lipgloss.NewStyle().Foreground(dimColor)
	sourceText = sourceStyle.Render(sourceText)

	// Add filled heart icon for favorites beside source
	if isFavorite {
		favoriteIcon := lipgloss.NewStyle().Foreground(errorColor).Render(" ❤")
		return sourceText + favoriteIcon
	}

	return sourceText
}

// renderDeleteConfirmation renders the delete confirmation dialog
func (m Model) renderDeleteConfirmation() string {
	if m.deleteConfirmHost == nil {
		return ""
	}

	host := m.deleteConfirmHost

	// Warning icon and title
	warningIcon := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("⚠ Delete Host?")

	// Host details
	hostDetails := lipgloss.NewStyle().
		Foreground(textColor).
		MarginTop(1).
		MarginBottom(1).
		Render(fmt.Sprintf("Alias: %s\nHost: %s@%s", host.Alias, host.User, host.Hostname))

	// Confirmation message
	confirmMsg := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render("This action cannot be undone.")

	// Action buttons
	yesButton := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("Y")

	noButton := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render("N")

	actions := lipgloss.NewStyle().
		MarginTop(1).
		Render(yesButton + descStyle.Render(" Yes  ") + noButton + descStyle.Render(" No  ") +
			keyStyle.Render("esc") + descStyle.Render(" Cancel"))

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left,
		warningIcon,
		hostDetails,
		confirmMsg,
		actions,
	)

	// Wrap in a dialog box
	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(2, 4).
		Width(50).
		Render(content)

	// Center on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}

// renderConfigError renders the config validation error screen
func (m Model) renderConfigError() string {
	// Error icon and title
	errorIcon := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("⚠ Configuration Errors")

	// Error count - determine source of error
	errorSource := "configuration"
	if len(m.configErrors) > 0 {
		// Check if error is from termix by looking at the error message
		firstError := m.configErrors[0].Error()
		if strings.Contains(firstError, "termix") {
			errorSource = "Termix"
		} else if strings.Contains(firstError, "Config:") {
			errorSource = "configuration"
		} else {
			errorSource = "sshbuddy.json"
		}
	}

	errorCount := lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginTop(1).
		Render(fmt.Sprintf("Found %d error(s) in %s:", len(m.configErrors), errorSource))

	// List errors (limit to first 10)
	var errorLines []string
	maxErrors := 10
	for i, err := range m.configErrors {
		if i >= maxErrors {
			remaining := len(m.configErrors) - maxErrors
			errorLines = append(errorLines, lipgloss.NewStyle().
				Foreground(dimColor).
				Italic(true).
				Render(fmt.Sprintf("... and %d more error(s)", remaining)))
			break
		}

		errorLine := lipgloss.NewStyle().
			Foreground(errorColor).
			Render(fmt.Sprintf("• %s", err.Error()))
		errorLines = append(errorLines, errorLine)
	}

	errorList := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(strings.Join(errorLines, "\n"))

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Please fix the errors in your config file.")

	// Action buttons
	ignoreButton := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render("I")

	quitButton := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("Q")

	actions := lipgloss.NewStyle().
		MarginTop(1).
		Render(ignoreButton + descStyle.Render(" Ignore & Continue  ") +
			quitButton + descStyle.Render(" Quit"))

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left,
		errorIcon,
		errorCount,
		errorList,
		instructions,
		actions,
	)

	// Wrap in a dialog box
	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Padding(2, 4).
		Width(70).
		Render(content)

	// Center on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}

// renderSourceSelect renders the source selection dialog
func (m Model) renderSourceSelect() string {
	if m.pendingConnectHost == nil {
		return ""
	}

	host := m.pendingConnectHost

	// Title
	titleText := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render("⚡ Select Source")

	// Host info
	hostInfo := lipgloss.NewStyle().
		Foreground(textColor).
		MarginTop(1).
		MarginBottom(1).
		Render(fmt.Sprintf("Host: %s (%s@%s)", host.Alias, host.User, host.Hostname))

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render("This host is available in multiple sources. Choose which one to use:")

	// Source list
	var sourceItems []string
	for i, source := range host.AvailableIn {
		// Get icon and name
		var icon, name string
		switch source {
		case "termix":
			icon = "▲"
			name = "Termix"
		case "ssh-config":
			icon = "■"
			name = "SSH Config"
		case "manual", "sshbuddy":
			icon = "◆"
			name = "SSHBuddy"
		default:
			icon = "○"
			name = source
		}

		// Format item
		var item string
		if i == m.selectedSourceIdx {
			// Selected item - highlighted
			item = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(fmt.Sprintf("▸ %s %s", icon, name))
		} else {
			// Normal item
			item = lipgloss.NewStyle().
				Foreground(dimColor).
				Render(fmt.Sprintf("  %s %s", icon, name))
		}
		sourceItems = append(sourceItems, item)
	}

	sourceList := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Render(strings.Join(sourceItems, "\n"))

	// Action buttons
	enterButton := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Render("↵")

	escButton := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("esc")

	actions := lipgloss.NewStyle().
		MarginTop(1).
		Render(enterButton + descStyle.Render(" Connect  ") +
			keyStyle.Render("↑↓") + descStyle.Render(" Navigate  ") +
			escButton + descStyle.Render(" Cancel"))

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleText,
		hostInfo,
		instructions,
		sourceList,
		actions,
	)

	// Wrap in a dialog box
	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(2, 4).
		Width(60).
		Render(content)

	// Center on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
