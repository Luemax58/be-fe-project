/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",   // ให้ Tailwind ตรวจทุกไฟล์ใน src
  ],
  theme: {
    extend: {
      colors: {
        viridian: {
          50: "#e9f7f5",
          100: "#c9ede7",
          200: "#a0e1d5",
          300: "#77d5c3",
          400: "#4ecab2",
          500: "#25bea1",
          600: "#1c9a82",
          700: "#147563",
          800: "#0c5145",
          900: "#052d27",
        },
      },
    },
  },
  plugins: [],
};
