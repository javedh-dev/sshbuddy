package tui

import (
	"fmt"
	"sshbuddy/internal/config"
	"sshbuddy/pkg/models"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateList sessionState = iota
	stateForm
	stateConfirmDelete
	stateConfigError
	stateConfig
	stateTermixAuth
	stateSourceSelect
)

type item struct {
	host     models.Host
	status   string // Ping status indicator
	pinging  bool   // Is currently being pinged
	pingTime string // Ping time in ms
}

func (i item) Title() string {
	// Colored dot based on ping status
	var statusText string
	if i.pinging {
		// Yellow dot for pinging in progress
		statusText = statusPingingStyle.Render("●")
	} else {
		switch i.status {
		case "🟢":
			statusText = statusOnlineStyle.Render("●")
		case "🔴":
			statusText = statusOfflineStyle.Render("●")
		default:
			statusText = statusUnknownStyle.Render("○")
		}
	}

	// Add ping time if available
	if i.pingTime != "" {
		return fmt.Sprintf("%s %s %s", statusText, i.host.Alias,
			lipgloss.NewStyle().Foreground(dimColor).Render(fmt.Sprintf("(%s)", i.pingTime)))
	}
	return fmt.Sprintf("%s %s", statusText, i.host.Alias)
}

func (i item) Description() string {
	port := i.host.Port
	if port == "" {
		port = "22"
	}
	return fmt.Sprintf("%s@%s:%s", i.host.User, i.host.Hostname, port)
}

func (i item) FilterValue() string { return i.host.Alias + i.host.Hostname }

type Model struct {
	list               list.Model
	form               FormModel
	configView         ConfigViewModel
	termixAuth         TermixAuthModel
	state              sessionState
	config             *models.Config
	pingStatus         map[string]bool   // track ping status for each host
	pinging            map[string]bool   // track which hosts are currently being pinged
	pingTimes          map[string]string // track ping times for each host
	width              int
	height             int
	selectedHost       *models.Host             // Host to connect to after quitting
	editingIndex       int                      // Index of host being edited (-1 if adding new)
	deleteConfirmHost  *models.Host             // Host pending deletion confirmation
	deleteConfirmIdx   int                      // Index of host pending deletion
	configErrors       []models.ValidationError // Config validation errors
	pendingConnectHost *models.Host             // Host pending source selection
	selectedSourceIdx  int                      // Selected source index in source selection dialog
}

func NewModel() Model {
	cfg, err := config.LoadConfig()
	var validationErrors []models.ValidationError
	var needsTermixAuth bool

	if err != nil {
		// Check if this is a Termix auth error
		if strings.Contains(err.Error(), "authentication required") {
			needsTermixAuth = true
			cfg = &models.Config{Hosts: []models.Host{}}
		} else {
			// Convert error to validation error for display
			validationErrors = []models.ValidationError{
				{
					Field:   "Config",
					Message: err.Error(),
					Index:   -1,
				},
			}
			cfg = &models.Config{Hosts: []models.Host{}}
		}
	} else {
		// Validate config
		validationErrors = cfg.Validate()
	}

	// Apply saved theme or default to purple
	themeName := cfg.Theme
	if themeName == "" {
		themeName = "purple"
	}
	ApplyTheme(themeName)

	items := []list.Item{}
	for _, h := range cfg.Hosts {
		items = append(items, item{host: h, status: "⚪"})
	}

	// Custom delegate with original styling
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3) // Three lines per item (title + description + tags)
	delegate.SetSpacing(0)

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(primaryColor).
		Bold(true).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 0, 0, 1)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(mutedColor).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 0, 0, 1)

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(primaryColor).
		Padding(0, 0, 0, 2)

	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(dimColor).
		Padding(0, 0, 0, 2)

	l := list.New(items, delegate, 0, 0)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle()
	l.Styles.StatusBar = lipgloss.NewStyle()

	m := Model{
		list:         l,
		form:         NewFormModel(),
		configView:   NewConfigViewModel(),
		termixAuth:   NewTermixAuthModel(),
		state:        stateList,
		config:       cfg,
		pingStatus:   make(map[string]bool),
		pinging:      make(map[string]bool),
		pingTimes:    make(map[string]string),
		editingIndex: -1,
		configErrors: validationErrors,
	}

	// If Termix auth is needed, show auth form
	if needsTermixAuth {
		m.state = stateTermixAuth
	} else if len(validationErrors) > 0 {
		// If there are validation errors, show error state
		m.state = stateConfigError
	}

	return m
}

func (m Model) Init() tea.Cmd {
	// Mark all hosts as pinging on startup
	for _, h := range m.config.Hosts {
		key := GetHostKey(h)
		m.pinging[key] = true
	}
	return StartPingAll(m.config.Hosts)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.state == stateList {
			// Check if we're in search/filter mode - if so, only allow escape and let list handle other keys
			filterState := m.list.FilterState()
			isSearching := filterState == list.Filtering

			// Handle Enter key in both search and normal mode
			if msg.String() == "enter" {
				// Connect to selected host
				if selectedItem, ok := m.list.SelectedItem().(item); ok {
					// Check if host is available in multiple sources
					if len(selectedItem.host.AvailableIn) > 1 {
						// Show source selection dialog
						m.pendingConnectHost = &selectedItem.host
						m.selectedSourceIdx = 0
						m.state = stateSourceSelect
						return m, nil
					}
					// Single source - connect directly
					return m, func() tea.Msg {
						return ConnectMsg{Host: selectedItem.host}
					}
				}
			}

			// Only process shortcuts when NOT in search mode
			if !isSearching {
				switch msg.String() {
				case "s":
					// Open settings/configuration
					m.state = stateConfig
					m.configView = NewConfigViewModel()
					m.configView.width = m.width
					m.configView.height = m.height
					return m, m.configView.Init()
				case "n":
					m.state = stateForm
					m.form = NewFormModel() // Reset form
					m.form.width = m.width
					m.form.height = m.height
					m.editingIndex = -1 // -1 means adding new
					return m, m.form.Init()
				case "p":
					// Ping all servers - mark all as pinging
					for _, h := range m.config.Hosts {
						key := GetHostKey(h)
						m.pinging[key] = true
					}
					m.refreshList()
					return m, StartPingAll(m.config.Hosts)
				case "up", "k":
					// Move up in 2-column layout (go back 2 items)
					currentIdx := m.list.Index()
					if currentIdx >= 2 {
						m.list.Select(currentIdx - 2)
					}
					return m, nil
				case "down", "j":
					// Move down in 2-column layout (go forward 2 items)
					currentIdx := m.list.Index()
					totalItems := len(m.list.Items())
					if currentIdx+2 < totalItems {
						m.list.Select(currentIdx + 2)
					}
					return m, nil
				case "left", "h":
					// Move left in row-wise layout (decrement by 1 if on odd index)
					currentIdx := m.list.Index()
					if currentIdx%2 == 1 { // If on right column
						m.list.Select(currentIdx - 1)
					}
					return m, nil
				case "right", "l":
					// Move right in row-wise layout (increment by 1 if on even index)
					currentIdx := m.list.Index()
					totalItems := len(m.list.Items())
					if currentIdx%2 == 0 && currentIdx+1 < totalItems { // If on left column
						m.list.Select(currentIdx + 1)
					}
					return m, nil
				case "e":
					// Edit selected host (only if from manual/sshbuddy source)
					if selectedItem, ok := m.list.SelectedItem().(item); ok {
						if selectedItem.host.Source == "ssh-config" || selectedItem.host.Source == "termix" {
							// Cannot edit SSH config or Termix hosts
							return m, nil
						}
						m.state = stateForm
						m.form = NewFormModelWithHost(selectedItem.host)
						m.form.width = m.width
						m.form.height = m.height
						m.editingIndex = m.list.Index()
						return m, m.form.Init()
					}
				case "c":
					// Duplicate selected host
					if selectedItem, ok := m.list.SelectedItem().(item); ok {
						m.state = stateForm
						duplicatedHost := selectedItem.host
						// Append " (copy)" to the alias to avoid duplicates
						duplicatedHost.Alias = duplicatedHost.Alias + " (copy)"
						m.form = NewFormModelWithHost(duplicatedHost)
						m.form.width = m.width
						m.form.height = m.height
						m.editingIndex = -1 // -1 means adding new (not editing)
						return m, m.form.Init()
					}
				case "d", "delete":
					// Show delete confirmation (only if from manual/sshbuddy source)
					if selectedItem, ok := m.list.SelectedItem().(item); ok {
						if selectedItem.host.Source == "ssh-config" || selectedItem.host.Source == "termix" {
							// Cannot delete SSH config or Termix hosts
							return m, nil
						}
						currentIdx := m.list.Index()
						if currentIdx >= 0 && currentIdx < len(m.config.Hosts) {
							m.deleteConfirmHost = &selectedItem.host
							m.deleteConfirmIdx = currentIdx
							m.state = stateConfirmDelete
						}
					}
					return m, nil
				case "f":
					// Toggle favorite status
					if _, ok := m.list.SelectedItem().(item); ok {
						return m, func() tea.Msg {
							return ToggleFavoriteMsg{}
						}
					}
				}
			}
		} else if m.state == stateForm {
			if msg.String() == "esc" {
				m.state = stateList
				return m, nil
			}
		} else if m.state == stateConfig {
			if msg.String() == "esc" {
				// Reload config in case it was changed
				cfg, err := config.LoadConfig()
				if err == nil {
					m.config = cfg
					m.refreshList()
				}
				m.state = stateList
				return m, nil
			}
		} else if m.state == stateTermixAuth {
			if msg.String() == "esc" {
				// Cancel auth and return to list (without Termix hosts)
				m.state = stateList
				return m, nil
			}
		} else if m.state == stateConfirmDelete {
			switch msg.String() {
			case "y", "Y":
				// Confirm deletion
				if m.deleteConfirmIdx >= 0 && m.deleteConfirmIdx < len(m.config.Hosts) {
					m.config.Hosts = append(m.config.Hosts[:m.deleteConfirmIdx], m.config.Hosts[m.deleteConfirmIdx+1:]...)
					config.SaveConfig(m.config)
					m.refreshList()
					// Adjust selection if needed
					if m.deleteConfirmIdx >= len(m.config.Hosts) && len(m.config.Hosts) > 0 {
						m.list.Select(len(m.config.Hosts) - 1)
					}
				}
				m.deleteConfirmHost = nil
				m.state = stateList
				return m, nil
			case "n", "N", "esc":
				// Cancel deletion
				m.deleteConfirmHost = nil
				m.state = stateList
				return m, nil
			}
		} else if m.state == stateConfigError {
			switch msg.String() {
			case "e", "E":
				// Open config file for editing
				m.state = stateList
				return m, nil
			case "i", "I":
				// Ignore errors and continue
				m.configErrors = nil
				m.state = stateList
				return m, nil
			case "q", "Q":
				return m, tea.Quit
			}
		} else if m.state == stateSourceSelect {
			switch msg.String() {
			case "up", "k":
				// Move up in source list
				if m.selectedSourceIdx > 0 {
					m.selectedSourceIdx--
				}
				return m, nil
			case "down", "j":
				// Move down in source list
				if m.pendingConnectHost != nil && m.selectedSourceIdx < len(m.pendingConnectHost.AvailableIn)-1 {
					m.selectedSourceIdx++
				}
				return m, nil
			case "enter":
				// Connect with selected source
				if m.pendingConnectHost != nil {
					// Create a copy of the host with the selected source as primary
					selectedSource := m.pendingConnectHost.AvailableIn[m.selectedSourceIdx]
					hostCopy := *m.pendingConnectHost
					hostCopy.Source = selectedSource
					m.pendingConnectHost = nil
					m.state = stateList
					return m, func() tea.Msg {
						return ConnectMsg{Host: hostCopy}
					}
				}
				m.state = stateList
				return m, nil
			case "esc":
				// Cancel source selection
				m.pendingConnectHost = nil
				m.state = stateList
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Fixed width box for 2-column layout
		const boxWidth = 80
		listWidth := boxWidth - 8 // Account for padding and borders
		listHeight := 20          // Height for scrollable list
		m.list.SetSize(listWidth, listHeight)

		// Update config view size
		m.configView.width = msg.Width
		m.configView.height = msg.Height

		// Update form size
		m.form.width = msg.Width
		m.form.height = msg.Height

		// Update termix auth size
		m.termixAuth.width = msg.Width
		m.termixAuth.height = msg.Height

	case PingResultMsg:
		// Update ping status, time, and clear pinging state
		key := GetHostKey(msg.Host)
		m.pingStatus[key] = msg.Status
		m.pingTimes[key] = msg.PingTime
		m.pinging[key] = false
		m.refreshList()
		return m, nil

	case FormSubmittedMsg:
		if m.editingIndex >= 0 && m.editingIndex < len(m.config.Hosts) {
			// Editing existing host
			m.config.Hosts[m.editingIndex] = msg.Host
		} else {
			// Adding new host
			m.config.Hosts = append(m.config.Hosts, msg.Host)
		}
		config.SaveConfig(m.config)
		m.state = stateList
		m.editingIndex = -1
		m.refreshList()
		// Ping the host
		return m, PingHost(msg.Host)

	case ConnectMsg:
		// Store the host and quit the TUI
		m.selectedHost = &msg.Host
		return m, tea.Quit

	case TermixAuthSuccessMsg:
		// Reload config after successful auth
		cfg, err := config.LoadConfig()
		if err != nil {
			// If still auth error, stay in auth state
			if strings.Contains(err.Error(), "authentication required") {
				return m, nil
			}
			// Other errors - show error state
			m.configErrors = []models.ValidationError{
				{
					Field:   "Config",
					Message: err.Error(),
					Index:   -1,
				},
			}
			m.state = stateConfigError
			return m, nil
		}
		m.config = cfg
		m.refreshList()
		m.state = stateList
		// Start pinging all hosts
		for _, h := range m.config.Hosts {
			key := GetHostKey(h)
			m.pinging[key] = true
		}
		return m, StartPingAll(m.config.Hosts)

	case ToggleFavoriteMsg:
		// Toggle favorite status for selected host
		currentIdx := m.list.Index()
		if currentIdx >= 0 && currentIdx < len(m.config.Hosts) {
			// Toggle favorite
			m.config.Hosts[currentIdx].Favorite = !m.config.Hosts[currentIdx].Favorite

			// Update favorites map
			if m.config.Favorites == nil {
				m.config.Favorites = make(map[string]bool)
			}
			alias := m.config.Hosts[currentIdx].Alias
			if m.config.Hosts[currentIdx].Favorite {
				m.config.Favorites[alias] = true
			} else {
				delete(m.config.Favorites, alias)
			}

			// Save config
			config.SaveConfig(m.config)

			// Reload config to re-sort hosts
			cfg, err := config.LoadConfig()
			if err == nil {
				m.config = cfg
				// Find the host we just toggled and select it
				for i, h := range m.config.Hosts {
					if h.Alias == alias {
						m.list.Select(i)
						break
					}
				}
			}
			m.refreshList()
		}
		return m, nil
	}

	if m.state == stateList {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateForm {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateConfig {
		m.configView, cmd = m.configView.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateTermixAuth {
		m.termixAuth, cmd = m.termixAuth.Update(msg)
		cmds = append(cmds, cmd)
	}
	// No update needed for stateConfirmDelete

	return m, tea.Batch(cmds...)
}

// View is implemented in view.go

func (m *Model) refreshList() {
	items := []list.Item{}
	for _, h := range m.config.Hosts {
		key := GetHostKey(h)
		status := "⚪" // Default - unknown
		if pingStatus, exists := m.pingStatus[key]; exists {
			status = GetHostStatus(pingStatus)
		}
		isPinging := m.pinging[key]
		pingTime := m.pingTimes[key]
		items = append(items, item{host: h, status: status, pinging: isPinging, pingTime: pingTime})
	}
	m.list.SetItems(items)
}

// GetSelectedHost returns the host selected for SSH connection
func (m Model) GetSelectedHost() *models.Host {
	return m.selectedHost
}
