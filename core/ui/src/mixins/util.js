export default {
  name: "UtilService", //// rename to ErrorService?
  data() {
    return {
      GENERIC_ERROR_MESSAGE: "Something went wrong",
    };
  },
  methods: {
    getErrorMessage(error) {
      if (!error.response && error.message === "Network Error") {
        return "Network error"; //// use i18n string
      }

      switch (error.response.status) {
        case 403:
          return "You don't have permission to access this resource";
      }

      return this.GENERIC_ERROR_MESSAGE;
    },
  },
};
