/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#eef8ff",
          100: "#d8efff",
          200: "#b7e1ff",
          300: "#86ceff",
          400: "#4fb6ff",
          500: "#199cff",
          600: "#0a7fe0",
          700: "#0b67b0",
          800: "#0d558f",
          900: "#0f4878"
        }
      },
      boxShadow: {
        soft: "0 20px 60px rgba(15, 72, 120, 0.12)",
        card: "0 12px 30px rgba(15, 72, 120, 0.12)",
        glow: "0 25px 70px rgba(25, 156, 255, 0.35)"
      }
    }
  },
  plugins: []
};
