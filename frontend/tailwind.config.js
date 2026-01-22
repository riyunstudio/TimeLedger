/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          500: '#6366F1',
          600: '#4F46E5',
          700: '#4338CA',
        },
        secondary: {
          500: '#A855F7',
          600: '#9333EA',
          700: '#7E22CE',
        },
        success: {
          500: '#10B981',
          600: '#059669',
          700: '#047857',
        },
        critical: {
          500: '#F43F5E',
          600: '#E11D48',
          700: '#BE123C',
        },
        warning: {
          500: '#F59E0B',
          600: '#D97706',
          700: '#B45309',
        },
        glass: {
          dark: 'rgba(30, 41, 59, 0.7)',
          light: 'rgba(255, 255, 255, 0.7)',
        },
      },
      fontFamily: {
        sans: ['Outfit', 'sans-serif'],
        heading: ['Inter', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-mesh': "linear-gradient(135deg, #1e3a8a 0%, #6366F1 25%, #A855F7 50%, #1e3a8a 75%, #0f172a 100%)",
        'gradient-mesh-light': "linear-gradient(135deg, #dbeafe 0%, #bfdbfe 25%, #e9d5ff 50%, #dbeafe 75%, #eff6ff 100%)",
      },
      backdropBlur: {
        'glass': '12px',
      },
    },
  },
  plugins: [],
}
