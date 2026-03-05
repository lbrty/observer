import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import de from "@/locales/de.json";
import en from "@/locales/en.json";
import ky from "@/locales/ky.json";
import ru from "@/locales/ru.json";
import tr from "@/locales/tr.json";
import uk from "@/locales/uk.json";

i18n.use(initReactI18next).init({
  resources: {
    ky: { translation: ky },
    en: { translation: en },
    ru: { translation: ru },
    uk: { translation: uk },
    de: { translation: de },
    tr: { translation: tr },
  },
  lng: localStorage.getItem("observer-lang") || "ky",
  fallbackLng: "ky",
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
