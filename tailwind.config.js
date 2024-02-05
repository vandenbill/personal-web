/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./public/**/*.{html,js}"],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Inter'],
        'mono': ['Jetbrains Mono'],
      }
    },
  },
  plugins: [],
}