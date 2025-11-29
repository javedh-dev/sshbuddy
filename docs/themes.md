# Themes

SSHBuddy offers six carefully crafted color themes to match your terminal aesthetic and personal preferences. Each theme uses consistent colors for status indicators while varying the primary interface colors.

## Available Themes

### Purple Dream (Default)

A modern, professional look with soft purple tones. This theme works well in both light and dark terminals and provides excellent contrast for extended use.

**Primary Color**: Purple (#7C3AED)

![Purple Dream Theme](screenshots/theme-purple.png)

### Ocean Blue

Cool blue tones inspired by the ocean. Perfect for those who prefer a calming, professional appearance.

**Primary Color**: Blue (#2563EB)

![Ocean Blue Theme](screenshots/theme-blue.png)

### Matrix Green

Classic terminal green for that retro hacker aesthetic. Ideal if you want your SSH manager to feel like home in a traditional terminal environment.

**Primary Color**: Green (#059669)

![Matrix Green Theme](screenshots/theme-green.png)

### Bubblegum Pink

Vibrant pink for a fun, energetic interface. This theme adds personality while maintaining readability.

**Primary Color**: Pink (#DB2777)

![Bubblegum Pink Theme](screenshots/theme-pink.png)

### Sunset Amber

Warm amber and orange tones reminiscent of a beautiful sunset. Great for reducing eye strain during long sessions.

**Primary Color**: Amber (#D97706)

![Sunset Amber Theme](screenshots/theme-amber.png)

### Cyber Cyan

Electric cyan for a futuristic cyberpunk aesthetic. This theme stands out while remaining easy on the eyes.

**Primary Color**: Cyan (#0891B2)

![Cyber Cyan Theme](screenshots/theme-cyan.png)

## Changing Themes

1. Press `s` to open settings
2. Navigate to "Theme" (marked with a colored diamond ◆)
3. Press Space or Enter to cycle through available themes
4. Your selection is saved automatically

The diamond icon next to "Theme" displays in the current theme's primary color, giving you a preview before you switch.

## Theme Elements

Each theme affects these interface elements:

**Primary Color** (varies by theme):
- ASCII art header
- Host aliases/names
- Selected item borders
- Active form fields
- Settings menu highlights

**Consistent Colors** (same across all themes):
- Status indicators:
  - Green (#10B981) - Online hosts
  - Red (#EF4444) - Offline hosts
  - Gray (#9CA3AF) - Unknown status
  - Yellow (#F59E0B) - Pinging in progress
- Source icons:
  - All source indicators use a muted gray for consistency

**Accent Colors** (varies by theme):
- Secondary highlights
- Checkmarks in settings
- Complementary UI elements

## Design Philosophy

SSHBuddy's themes follow these principles:

**Readability First**: All themes maintain high contrast between text and background, ensuring comfortable reading during extended sessions.

**Consistent Status**: Status indicators (online/offline/pinging) use the same colors across all themes. This consistency helps you quickly assess host availability regardless of which theme you're using.

**Terminal Agnostic**: Themes work well in both light and dark terminal backgrounds. The colors are carefully chosen to provide good contrast in various environments.

**Professional Yet Personal**: While maintaining a professional appearance suitable for work environments, each theme offers enough personality to make SSHBuddy feel like your own tool.

## Customization

Currently, SSHBuddy offers six pre-defined themes. If you need custom colors, you can modify the theme definitions in the source code and rebuild the application. Future versions may include support for user-defined themes through the configuration file.

## Theme Persistence

Your theme choice is saved in the configuration file and persists across sessions. When you launch SSHBuddy, it automatically applies your last selected theme.
