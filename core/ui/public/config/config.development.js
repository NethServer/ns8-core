window.CONFIG = {
  PRODUCT_NAME: "NethServer 8",
  PRODUCT_URL: "https://www.nethserver.org/",
  DOCS_URL: "https://www.nethserver.org/documentation/",
  HELPDESK_URL: "https://www.nethserver.org/documentation/", ////
  COMPANY_NAME: "Nethesis", ////
  COMPANY_URL: "https://www.nethesis.it/", ////
  // API_ENDPOINT: window.location.hostname + ":8081/api", ////
  API_ENDPOINT: "192.168.122.183/cluster-admin/api", //// leader node
  // API_SCHEME: window.location.protocol + "//", ////
  API_SCHEME: "https://",
  // WS_ENDPOINT: "ws://" + window.location.hostname + ":8081/ws", ////
  WS_ENDPOINT: "wss://192.168.122.183/cluster-admin/ws", //// leader node
};
