nethserver = {
  // used to map a query string parameter value to its typed value
  getTypedValue(stringValue) {
    if (stringValue === "true") {
      return true;
    }

    if (stringValue === "false") {
      return false;
    }

    return stringValue;
  },
  // used only by NS8 core to extract query parameters from URL
  getQueryParamsForCore() {
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
  getPage() {
    let queryParams = nethserver.getQueryParamsForApp();
    const page = queryParams.page || "status";
    return page;
  },
  // used only by external apps to sync UI status with URL query parameters
  initUrlBinding(context, page) {
    console.log("initUrlBinding, page", page); ////

    let queryParams = nethserver.getQueryParamsForApp();
    const requestedPage = queryParams.page || "status";

    if (requestedPage != page) {
      console.log("page mismatch, return", requestedPage, page); ////
      return;
    }

    nethserver.syncQueryParamsAndData(context);
    return setInterval(() => nethserver.checkUrlChange(context, page), 500);
  },
  // used only by external apps to sync UI status with URL query parameters
  syncQueryParamsAndData(context) {
    nethserver.queryParamsToData(context);
    nethserver.dataToQueryParams(context);
  },
  // used by vuejs core and external apps to dynamically register a watch for every object inside data.q
  watchQueryData(context) {
    Object.keys(context.q).forEach((dataItem) => {
      context.$watch("q." + dataItem, function () {
        console.log("watch", dataItem); ////
        nethserver.dataToQueryParams(context);
      });
    });
  },
  // used only by external apps to sync UI status with URL query parameters
  checkUrlChange(context, page) {
    const newUrl = window.parent.location.href;

    if (newUrl != context.currentUrl) {
      console.log("url changed!"); ////
      context.currentUrl = newUrl;
      const queryParams = nethserver.getQueryParamsForApp();
      const requestedPage = queryParams.page || "status";

      // emit app navigation event
      const appInstance = /#\/apps\/(\w+)/.exec(window.parent.location.hash)[1];
      context.$root.$emit("appNavigation", appInstance, requestedPage);

      if (requestedPage !== page) {
        // page has changed, need to change route
        console.log("detected page change, dispatch event", requestedPage); ////

        // trigger a custom event to change app route
        window.dispatchEvent(
          new CustomEvent("changeRoute", { detail: requestedPage })
        );
      } else {
        nethserver.syncQueryParamsAndData(context);
      }
    }
    context.currentUrl = window.parent.location.href;
  },
  // used only by external apps to sync UI status with URL query parameters
  queryParamsToData(context) {
    let queryParams = nethserver.getQueryParamsForApp();

    console.log("queryParamsToData, queryParams", queryParams); ////

    Object.keys(context.q).forEach((dataItem) => {
      if (typeof queryParams[dataItem] !== "undefined") {
        context.q[dataItem] = nethserver.getTypedValue(queryParams[dataItem]);
      }
    });
  },
  // used only by external apps to sync UI status with URL query parameters
  dataToQueryParams(context) {
    console.log("dataToQueryParams, q", context.q); ////

    let queryParams = [];

    for (const [key, value] of Object.entries(context.q)) {
      queryParams.push(key + "=" + encodeURIComponent(value));
    }

    const urlWithParams =
      window.parent.ns8.$route.path + "?" + queryParams.join("&");

    console.log(
      "urlWithParams, route",
      urlWithParams,
      window.parent.ns8.$route.fullPath
    ); ////

    // window.parent.history.replaceState(null, "", urlWithParams); ////
    if (window.parent.ns8.$route.fullPath != urlWithParams) {
      window.parent.ns8.$router.replace(urlWithParams);
    }
  },
  // used only by external apps to extract query parameters from URL
  getQueryParamsForApp() {
    if (
      !window.parent.location.hash.includes("?") ||
      window.parent.location.hash.split("?").length < 2
    ) {
      return {};
    }

    const params = new URLSearchParams(
      window.parent.location.hash.split("?").pop()
    );
    let queryParams = {};

    params.forEach((value, key) => {
      if (key) {
        queryParams[key] = value;
      }
    });
    return queryParams;
  },
};
