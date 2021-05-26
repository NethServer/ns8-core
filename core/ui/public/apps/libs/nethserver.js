nethserver = {
  initUrlBinding(context, page) {
    console.log("initUrlBinding, page", page); ////

    let queryParams = nethserver.getQueryParams();
    const requestedPage = queryParams.page || "home";

    if (requestedPage != page) {
      console.log("page mismatch, return", requestedPage, page); ////
      return;
    }

    nethserver.watchData(context);
    nethserver.syncQueryParamsAndData(context);
    return setInterval(() => nethserver.checkUrlChange(context, page), 500);
  },
  syncQueryParamsAndData(context) {
    nethserver.queryParamsToData(context);
    nethserver.dataToQueryParams(context);
  },
  watchData(context) {
    Object.keys(context._data.q).forEach((dataItem) => {
      context.$watch("q." + dataItem, function () {
        console.log("watch", dataItem); ////
        nethserver.dataToQueryParams(context);
      });
    });
  },
  checkUrlChange(context, page) {
    const newUrl = window.parent.location.href;

    if (newUrl != context.currentUrl) {
      console.log("url changed"); ////
      context.currentUrl = newUrl;
      const queryParams = nethserver.getQueryParams();
      const requestedPage = queryParams.page || "home";

      if (requestedPage !== page) {
        // page has changed, need to change route
        console.log("detected page change, dispatch event", requestedPage); ////

        // trigger a custom event to change app route
        window.dispatchEvent(
          new CustomEvent("changeRoute", { detail: requestedPage })
        );
      } else {
        nethserver.syncQueryParamsAndData(context); //// remove context.?
      }
    }
    context.currentUrl = window.parent.location.href;
  },
  queryParamsToData(context) {
    let queryParams = nethserver.getQueryParams();

    console.log("queryParamsToData, queryParams", queryParams); ////

    Object.keys(context._data.q).forEach((dataItem) => {
      if (typeof queryParams[dataItem] !== "undefined") {
        context.q[dataItem] = nethserver.getTypedValue(queryParams[dataItem]);
      }
    });
  },
  getTypedValue(value) {
    if (value === "true") {
      return true;
    }

    if (value === "false") {
      return false;
    }

    return value;
  },
  dataToQueryParams(context) {
    console.log("dataToQueryParams, q", context._data.q); ////

    let queryParams = [];

    for (const [key, value] of Object.entries(context._data.q)) {
      queryParams.push(key + "=" + value);
    }

    const baseUrl = window.parent.location.hash.split("?").shift();
    const urlWithParams = baseUrl + "?" + queryParams.join("&");
    window.parent.history.replaceState(null, "", urlWithParams);
  },
  getQueryParams() {
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
