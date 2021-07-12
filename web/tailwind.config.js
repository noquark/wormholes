const colors = require('tailwindcss/colors')

module.exports = {
  mode: 'jit',
  purge: ['./components/**/*.js', './pages/**/*.js'],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        gray: colors.coolGray,
      },
    },
  },
}
