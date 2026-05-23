/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  theme: {
    extend: {
      colors: {
        accent: {
          50: '#effaff',
          100: '#dff4ff',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
        },
      },
      borderRadius: {
        fluent: '1.25rem',
      },
      boxShadow: {
        fluent: '0 24px 80px rgba(14, 116, 144, 0.16)',
        subtle: '0 12px 40px rgba(14, 116, 144, 0.10)',
      },
      fontFamily: {
        sans: [
          'Inter',
          'ui-sans-serif',
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'sans-serif',
        ],
      },
    },
  },
}
