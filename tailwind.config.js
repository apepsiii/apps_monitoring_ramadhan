/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    "./web/templates/**/*.html",
    "./web/static/js/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#0D7E5E',
          50: '#E7F5F0',
          100: '#C4EBE0',
          200: '#9FDACC',
          300: '#75C8B6',
          400: '#4DB69E',
          500: '#2DA487',
          600: '#0D7E5E',
          700: '#0A6349',
          800: '#074936',
          900: '#043023',
        },
        accent: {
          DEFAULT: '#D4AF37',
          50: '#FBF6E8',
          100: '#F6EACC',
          200: '#EED999',
          300: '#E5C566',
          400: '#DCB144',
          500: '#D4AF37',
          600: '#B8941F',
          700: '#967717',
          800: '#735B12',
          900: '#503F0D',
        },
        warm: {
          50: '#F8F7F4',
          100: '#E8E6E1',
          200: '#D4D1C9',
          300: '#BFB9AD',
        },
        secondary: {
          50: '#f5f3ff',
          100: '#ede9fe',
          200: '#ddd6fe',
          300: '#c4b5fd',
          400: '#a78bfa',
          500: '#8b5cf6',
          600: '#7c3aed',
          700: '#6d28d9',
          800: '#5b21b6',
          900: '#4c1d95',
        },
      },
      fontFamily: {
        sans: ['Cairo', 'Inter', 'system-ui', '-apple-system', 'sans-serif'],
        arabic: ['LPMQ IsepMisbah', 'Amiri', 'serif'],
        display: ['Cairo', 'Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'soft': '0 2px 15px -3px rgba(0, 0, 0, 0.07), 0 10px 20px -2px rgba(0, 0, 0, 0.04)',
        'soft-lg': '0 10px 40px -10px rgba(0, 0, 0, 0.1)',
        'card': '0 4px 6px rgba(0, 0, 0, 0.1)',
        'card-hover': '0 10px 25px rgba(0, 0, 0, 0.15)',
      },
      borderRadius: {
        'card': '20px',
        'container': '48px',
        '2xl': '1rem',
        '3xl': '1.5rem',
        '4xl': '2rem',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in',
        'slide-up': 'slideUp 0.3s ease-out',
        'bounce-soft': 'bounceSoft 2s infinite',
        'pulse-soft': 'pulseSoft 2s ease-in-out infinite',
        'float': 'float 3s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideUp: {
          '0%': { transform: 'translateY(100%)' },
          '100%': { transform: 'translateY(0)' },
        },
        bounceSoft: {
          '0%, 100%': { transform: 'translateY(-5%)' },
          '50%': { transform: 'translateY(0)' },
        },
        pulseSoft: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.7' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        },
      },
      backgroundImage: {
        'gradient-primary': 'linear-gradient(135deg, #0D7E5E 0%, #0A6349 100%)',
        'gradient-accent': 'linear-gradient(135deg, #D4AF37 0%, #B8941F 100%)',
        'gradient-header': 'linear-gradient(180deg, #0D7E5E 0%, #0A6349 100%)',
        'gradient-warm': 'linear-gradient(135deg, #F4E5C2 0%, #E8D7B0 100%)',
      },
    },
  },
  plugins: [],
}
