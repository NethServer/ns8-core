export async function loadLanguage(lang) {
  try {
    const messages = await import(
      /* webpackChunkName: "lang-[request]" */ `../../public/i18n/${lang}/translation.json`
    );
    return messages;
  } catch (error) {
    console.warn(
      `Cannot import ${lang} language messages, falling back to English.`,
      error
    );

    // fallback to english
    if (lang !== "en") {
      return loadLanguage("en");
    }
  }
}

export async function getDateFnsLocale() {
  const lang = navigator.language;
  try {
    const mod = await import(
      /* webpackChunkName: "date-fns-locale-[request]" */
      `date-fns/locale/${lang}`
    );
    return mod.default || mod;
  } catch {
    // try base language (e.g. "it" from "it-IT")
    const baseLang = lang.split("-")[0];
    try {
      const mod = await import(
        /* webpackChunkName: "date-fns-locale-[request]" */
        `date-fns/locale/${baseLang}`
      );
      return mod.default || mod;
    } catch {
      const mod = await import("date-fns/locale/en-US");
      return mod.default || mod;
    }
  }
}