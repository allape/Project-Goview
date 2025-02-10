import { i18n } from "@allape/gocrud-react";
import i18next from "i18next";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.scss";
import { initReactI18next } from "react-i18next";
import App from "./App.tsx";

i18next.use(initReactI18next).init({
  resources: {
    en: {
      translation: {
        ...i18n,
      },
    },
  },
  lng: navigator.language,
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
