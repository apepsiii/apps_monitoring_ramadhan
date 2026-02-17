// Check for saved theme preference or system preference
if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
} else {
    document.documentElement.classList.remove('dark');
}

function toggleTheme() {
    if (document.documentElement.classList.contains('dark')) {
        document.documentElement.classList.remove('dark');
        localStorage.theme = 'light';
    } else {
        document.documentElement.classList.add('dark');
        localStorage.theme = 'dark';
    }
    updateThemeIcon();
}

function updateThemeIcon() {
    const isDark = document.documentElement.classList.contains('dark');
    // Finds all elements with data-toggle-theme-icon attribute
    const icons = document.querySelectorAll('[data-toggle-theme-icon]');
    icons.forEach(icon => {
        if (isDark) {
            // Show moon, hide sun (or vice versa depending on design)
            // Assuming button has specific SVG for current state
            // Simplified: Just toggle a class or swap SVG content if complex
            // For now, let's assume the button html will handle the swap based on 'dark' class on parent
        }
    });
}

// Optional: Listen for system preference changes
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
    if (!('theme' in localStorage)) {
        if (event.matches) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }
});
