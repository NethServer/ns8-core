export function loadLanguage(lang, i18n) {
  return import(
    /* webpackChunkName: "lang-[request]" */ `../../public/i18n/${lang}/translation.json`
  )
    .then((messages) => {
      i18n.setLocaleMessage(lang, messages.default);
      i18n.locale = lang;
    })
    .catch((error) => {
      console.warn(`Cannot import ${lang} language messages, falling back to English.`, error);

      // fallback to english

      if (lang !== "en") {
        loadLanguage("en", i18n);
      }
    });
}
