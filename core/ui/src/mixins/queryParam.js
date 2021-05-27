let nethserver = window.nethserver;

export default {
  name: "QueryParamService",
  methods: {
    queryParamsToData(context, queryParams) {
      // let queryParams = this.getQueryParams(); ////

      console.log("queryParamsToData, queryParams", queryParams); ////

      Object.keys(context.q).forEach((dataItem) => {
        if (typeof queryParams[dataItem] !== "undefined") {
          context.q[dataItem] = nethserver.getTypedValue(queryParams[dataItem]);
        }
      });

      // if (queryParams.testInput) { //// remove
      //   context.testInput = queryParams.testInput;
      // } else {
      //   context.testInput = "";
      // }

      // if (queryParams.testToggle) {
      //   context.testToggle = queryParams.testToggle === "true";
      // } else {
      //   context.testToggle = false;
      // }
    },
    getQueryParams() {
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
    // getTypedValue(stringValue) { ////
    //   if (stringValue === "true") {
    //     return true;
    //   }

    //   if (stringValue === "false") {
    //     return false;
    //   }

    //   return stringValue;
    // },
  },
};
