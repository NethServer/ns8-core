export default {
  name: "ModuleService",
  methods: {
    getAppDescription(app) {
      const langCode = this.$root.$i18n.locale;
      let description = app.description[langCode];

      if (!description) {
        // fallback to english
        description = app.description.en;
      }
      return description;
    },
    getAppCategories(app) {
      let i18nCategories = [];

      for (const category of app.categories) {
        if (category === "unknown") {
          return "-";
        }

        i18nCategories.push(
          this.$t("software_center.app_categories." + category)
        );
      }
      return i18nCategories.join(", ");
    },
  },
};
