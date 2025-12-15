Output demo.gif

Set FontSize 16
Set Width 1200
Set Height 800
Set Padding 20
Set FontFamily "JetBrains Mono"
Set Theme { "name": "FedRAMP-TUI", "background": "#1E1F29", "foreground": "#FFFFFF", "black": "#282A36", "red": "#FF5F56", "green": "#04B575", "yellow": "#FFCC00", "blue": "#45B7D1", "magenta": "#7D56F4", "cyan": "#4ECDC4", "white": "#F8F8F2", "brightBlack": "#626262", "brightRed": "#FF6B6B", "brightGreen": "#58D68D", "brightYellow": "#F7DC6F", "brightBlue": "#85C1E9", "brightMagenta": "#BB8FCE", "brightCyan": "#98D8C8", "brightWhite": "#FFFFFF", "cursor": "#7D56F4", "selection": "#44475A" }
Set TypingSpeed 50ms

# Launch app (hidden)
Hide
Type "./fedramp"
Enter
Sleep 3s
Show

# ══════════════════════════════════════════════════════════════
# 1. DOCUMENTS TAB - Browse FedRAMP documents with colored badges
# ══════════════════════════════════════════════════════════════
Sleep 1.5s

# Scroll through documents to show different colored badges (VDR, FRD, KSI, etc.)
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 1.5s

# Show help overlay
Type "?"
Sleep 2.5s
Type "?"
Sleep 500ms

# ══════════════════════════════════════════════════════════════
# 2. DOCUMENT DETAIL - Open a document and view its details
# ══════════════════════════════════════════════════════════════
Sleep 500ms

# Go back to top and select VDR document
Type "g"
Sleep 300ms
Type "g"
Sleep 500ms
Type "j"
Sleep 300ms
Type "j"
Sleep 500ms
Enter
Sleep 2.5s

# Scroll in detail view to show content
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 1.5s

# Go to requirements for this document
Enter
Sleep 2s

# ══════════════════════════════════════════════════════════════
# 3. REQUIREMENTS TAB - Filters and search
# ══════════════════════════════════════════════════════════════

# Clear document filter to show all requirements
Type "f"
Sleep 1s

# Toggle affects filter (cycle through: all → Providers → 3PAO → ...)
Type "x"
Sleep 1.2s
Type "x"
Sleep 1.2s
Type "x"
Sleep 1s

# Show MUST requirements only
Type "m"
Sleep 1.5s

# Switch to SHOULD requirements
Type "s"
Sleep 1.5s

# Clear filters
Type "f"
Sleep 1s

# Search for specific requirement
Type "/"
Type "vulnerability"
Enter
Sleep 1.5s

# Open requirement detail
Enter
Sleep 2.5s

# Scroll in detail view
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 1.5s
Escape
Sleep 500ms

# ══════════════════════════════════════════════════════════════
# 4. DEFINITIONS TAB - Search FedRAMP terminology
# ══════════════════════════════════════════════════════════════
Type "3"
Sleep 1.2s

# Search for a definition
Type "/"
Type "cloud"
Enter
Sleep 1.2s

# View definition detail
Enter
Sleep 2.5s
Escape
Sleep 500ms

# ══════════════════════════════════════════════════════════════
# 5. KEY SECURITY INDICATORS - Control mappings
# ══════════════════════════════════════════════════════════════
Type "4"
Sleep 1.2s

# Search for an indicator
Type "/"
Type "KSI-AFR"
Enter
Sleep 1.2s

# Open indicator to show control mappings
Enter
Sleep 2.5s

# Scroll to show control mappings
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 1.5s
Escape
Sleep 500ms

# ══════════════════════════════════════════════════════════════
# 6. QUICK TAB NAVIGATION
# ══════════════════════════════════════════════════════════════
Type "1"
Sleep 800ms
Type "2"
Sleep 800ms
Type "3"
Sleep 800ms
Type "4"
Sleep 800ms
Type "1"
Sleep 1s

# Quit
Type "q"
