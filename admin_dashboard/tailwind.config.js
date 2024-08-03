/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: false,
  theme: {
    extend: {
      colors:{
        accent: '#0A84FF',
        half_light: '#E6F3FB',
        light: '#F4FBFF',
      }
    },
  },
  plugins: [],
}

