export default {
  name: "UtilService",
  methods: {
    getErrorMessage(error) {
      if (
        error.message === "Network Error" ||
        /^timeout of .+ exceeded$/.test(error.message)
      ) {
        // network issues or axios timeout reached
        return this.$t("error.network_error");
      }

      switch (error.response.status) {
        case 401:
          return this.$t("error.401");
        case 403:
          return this.$t("error.403");
        case 404:
          return this.$t("error.404");
      }

      return this.$t("error.generic_error");
    },
    clearErrors(context) {
      console.log("errors before clear", context.error); ////

      for (const key of Object.keys(context.error)) {
        context.error[key] = "";
      }
      console.log("errors after clear", context.error); ////
    },
  },
};
