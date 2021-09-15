window.CONFIG = {
  PRODUCT_NAME: "NethServer 8",
  DOCS_URL: "https://www.nethserver.org/documentation/",
  COMPANY_NAME: "Nethesis",
  // API_ENDPOINT: window.location.hostname + ":8081/api", ////
  API_ENDPOINT: "172.25.5.229/cluster-admin/api", //// leader node
  // API_ENDPOINT: "172.25.5.3/cluster-admin/api", //// worker node
  // API_SCHEME: window.location.protocol + "//", ////
  API_SCHEME: "https://",
  // WS_ENDPOINT: "ws://" + window.location.hostname + ":8081/ws", ////
  WS_ENDPOINT: "wss://172.25.5.229/cluster-admin/ws", //// leader node
  // WS_ENDPOINT: "wss://172.25.5.3/cluster-admin/ws", //// worker node
};
