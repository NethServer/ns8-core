let nethserver = window.nethserver;

export default {
  name: "QueryParamService",
  methods: {
    queryParamsToData(context, queryParams) {
      Object.keys(context.q).forEach((dataItem) => {
        if (typeof queryParams[dataItem] !== "undefined") {
          context.q[dataItem] = nethserver.getTypedValue(queryParams[dataItem]);
        }
      });
    },
    getQueryParamsForApp() {
      if (
        !window.location.hash.includes("?") ||
        window.location.hash.split("?").length < 2
      ) {
        return {};
      }

      const params = new URLSearchParams(window.location.hash.split("?").pop());
      let queryParams = {};

      params.forEach((value, key) => {
        if (key) {
          queryParams[key] = value;
        }
      });
      return queryParams;
    },
  },
};
