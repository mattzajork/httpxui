module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  purge: ['./src/**/*.svelte', './src/**/*.css'],
  darkMode: 'media',
  theme: {
    fontFamily: {
      sans: ["'Inter'", 'sans-serif'],
      mono: ["'Source Code Pro'", 'mono']
    },
    extend: {
      backgroundImage: {
        'hero-pattern': "url('/img/3dgridlines.svg')",
      },
      colors: {
        'pdyellow': '#f9cb28',
        'pdorange': '#ff4d4d',
        'pdgray': '#131314',
        'pdgrayoutline': '#3f3f46',
        'pdcontentmuted': '#A1A1AA',
      },
    }
  },
  variants: {
    extend: {},
  },
  daisyui: {
    themes: [
      {
        mytheme: {
          "primary": "#ffffff",
          "primary-content": "#160202",
          "secondary": "#ff4d4d",
          "secondary-content": "#160202",
          "accent": "#f9cb28",
          "accent-content": "#150f01",
          "neutral": "#18181b",
          "neutral-content": "#ffffff",
          "base-100": "#09090b",
          "base-200": "#18181b",
          "base-300": "#18181b",
          "base-content": "#ffffff",
          "info": "#6367ed",
          "info-content": "#000a13",
          "success": "#1fc45d",
          "success-content": "#0a1301",
          "warning": "#f79c0a",
          "warning-content": "#160d05",
          "error": "#ef433f",
          "error-content": "#150404",
          },
        },
      ],
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('daisyui'),
    require('@tailwindcss/typography'),
  ],
}