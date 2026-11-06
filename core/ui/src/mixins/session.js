import LoginService from "@/mixins/login";
import { StorageService } from "@nethserver/ns8-ui-lib";

// refresh the access token this time before it expires
const REFRESH_ADVANCE_MS = 60 * 1000;
// retry a refresh that failed for a network or server error
const REFRESH_RETRY_MS = 30 * 1000;

// shared by all components of this browser tab
let refreshTimer = null;
let refreshPromise = null;

export default {
  name: "SessionService",
  mixins: [LoginService, StorageService],
  methods: {
    /**
     * Schedule the refresh of the access token, before it expires
     */
    scheduleTokenRefresh(delay) {
      this.cancelTokenRefresh();
      const loginInfo = this.getFromStorage("loginInfo");

      if (!loginInfo?.refreshToken || !loginInfo?.expire) {
        return;
      }

      if (delay === undefined) {
        delay = Math.max(
          new Date(loginInfo.expire).getTime() -
            Date.now() -
            REFRESH_ADVANCE_MS,
          0
        );
      }
      refreshTimer = setTimeout(this.onTokenRefreshTimer, delay);
    },
    cancelTokenRefresh() {
      if (refreshTimer) {
        clearTimeout(refreshTimer);
        refreshTimer = null;
      }
    },
    async onTokenRefreshTimer() {
      refreshTimer = null;
      const loginInfo = this.getFromStorage("loginInfo");
      const result = await this.refreshSession(loginInfo?.token);

      if (result === "ok") {
        this.$root.$emit("sessionTokenChanged");
      } else if (result === "invalid") {
        this.$root.$emit(
          "logout",
          this.$t("login.session_expired_title"),
          this.$t("login.session_expired_description")
        );
      } else {
        this.scheduleTokenRefresh(REFRESH_RETRY_MS);
      }
    },
    /**
     * Get a new access token with the refresh token. staleToken is the
     * access token the caller found expired: if another browser tab
     * already replaced it, no request is needed. Resolves to "ok",
     * "invalid" (the session ended: log in again) or "error".
     */
    refreshSession(staleToken) {
      // concurrent callers of this tab share the same request
      if (!refreshPromise) {
        refreshPromise = this.withRefreshLock(() =>
          this.doRefreshSession(staleToken)
        ).finally(() => {
          refreshPromise = null;
        });
      }
      return refreshPromise;
    },
    withRefreshLock(callback) {
      // browser tabs share the refresh token: one tab at a time can
      // replace it
      if (navigator.locks?.request) {
        return navigator.locks.request("ns8-session-refresh", callback);
      }
      return callback();
    },
    async doRefreshSession(staleToken) {
      const loginInfo = this.getFromStorage("loginInfo");

      if (!loginInfo?.refreshToken) {
        return "invalid";
      }

      if (
        staleToken &&
        loginInfo.token !== staleToken &&
        new Date(loginInfo.expire).getTime() - Date.now() > REFRESH_ADVANCE_MS
      ) {
        // another tab has already refreshed the token
        return "ok";
      }

      try {
        const response = await this.executeRefresh(loginInfo.refreshToken);
        this.saveToStorage("loginInfo", {
          ...loginInfo,
          token: response.data.token,
          refreshToken: response.data.refresh_token,
          expire: response.data.expire,
        });
        return "ok";
      } catch (error) {
        if (error.response?.status == 401) {
          return "invalid";
        }
        console.error("session refresh failed", error);
        return "error";
      }
    },
  },
};
