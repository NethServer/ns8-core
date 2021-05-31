export default {
  name: "ErrorService",
  data() {
    return {
      GENERIC_ERROR_MESSAGE: "Something went wrong", //// i18n
    };
  },
  methods: {
    getErrorMessage(error) {
      if (!error.response && error.message === "Network Error") {
        return "Network error"; //// use i18n string
      }

      switch (error.response.status) {
        case 401:
          return "Missing or invalid authentication"; //// i18n
        case 403:
          return "You don't have permission to access this resource"; //// i18n
      }

      return this.GENERIC_ERROR_MESSAGE;
    },
  },
};
