window.CONFIG = {
  PRODUCT_NAME: "NethServer 8",
  DOCS_URL: "https://www.nethserver.org/documentation/",
  COMPANY_NAME: "Nethesis",
  API_ENDPOINT:
    window.location.hostname +
    (window.location.port ? ":" + window.location.port : "") +
    window.location.pathname +
    "api",
  API_SCHEME: window.location.protocol + "//",
  WS_ENDPOINT:
    "wss://" +
    window.location.hostname +
    (window.location.port ? ":" + window.location.port : "") +
    window.location.pathname +
    "ws",
};
