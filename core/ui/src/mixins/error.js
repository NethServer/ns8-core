export default {
  name: "ErrorService",
  data() {
    return {
      GENERIC_ERROR_MESSAGE: "Something went wrong", //// i18n
    };
  },
  methods: {
    getErrorMessage(error) {
      // network issues or axios timeout reached
      if (
        error.message === "Network Error" ||
        /^timeout of .+ exceeded$/.test(error.message)
      ) {
        return "Network error"; //// i18n
      }

      switch (error.response.status) {
        case 401:
          return "Missing or invalid authentication"; //// i18n
        case 403:
          return "You don't have permission to access the resource"; //// i18n
        case 404:
          return "The resource cannot be found"; //// i18n
      }

      return this.GENERIC_ERROR_MESSAGE;
    },
  },
};
