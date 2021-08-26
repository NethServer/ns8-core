<template>
  <div>
    <div class="bx--grid bx--grid--full-width login-bg">
      <div class="bx--row">
        <div class="bx--offset-md-1 bx--col-md-6 bx--offset-lg-4 bx--col-lg-8">
          <div class="test">
            <cv-tile :light="true" class="login-tile">
              <h2 class="login-title">{{ $t("login.title") }}</h2>
              <div v-if="step === 'username'">
                <NsInlineNotification
                  v-if="error.login"
                  kind="error"
                  :title="$t('login.cannot_login')"
                  :description="error.login"
                  :showCloseButton="false"
                />
                <cv-form @submit.prevent="checkUsername" class="login-form">
                  <cv-text-input
                    :label="$t('login.username')"
                    v-model="username"
                    class="mg-bottom-md"
                    :placeholder="$t('login.username_placeholder')"
                    :invalid-message="error.username"
                    ref="usernameInput"
                    name="username"
                  ></cv-text-input>
                  <cv-checkbox
                    :label="$t('login.remember_username')"
                    class="mg-bottom-md"
                    v-model="rememberUsername"
                    value="checkRememberUsername"
                  />
                  <div class="login-footer">
                    <cv-button
                      kind="primary"
                      :icon="ArrowRight20"
                      class="login-button"
                      >Continue</cv-button
                    >
                  </div>
                </cv-form>
              </div>
              <div v-else-if="step === 'password'">
                <div class="login-as">
                  {{ $t("login.logging_in_as") }}
                  <strong>{{ username }}</strong>
                  <cv-link @click="goToUsername" class="not-you">{{
                    $t("login.not_you")
                  }}</cv-link>
                </div>
                <cv-form @submit.prevent="checkPassword" class="login-form">
                  <!-- hidden username field (to help browser saving credentials) -->
                  <cv-text-input
                    v-model="username"
                    name="username"
                    class="hidden"
                  ></cv-text-input>
                  <cv-text-input
                    :label="$t('login.password')"
                    type="password"
                    v-model="password"
                    class="bx--password-input mg-bottom-md"
                    :placeholder="$t('login.password_placeholder')"
                    :invalid-message="error.password"
                    :password-hide-label="$t('password.hide_password')"
                    :password-show-label="$t('password.show_password')"
                    ref="passwordInput"
                    name="password"
                  ></cv-text-input>
                  <div class="login-footer">
                    <NsButton
                      kind="primary"
                      :icon="ArrowRight20"
                      class="login-button"
                      :loading="loading.login"
                      :disabled="loading.login"
                      >{{ $t("login.login") }}</NsButton
                    >
                  </div>
                </cv-form>
              </div>
            </cv-tile>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { IconService } from "@nethserver/ns8-ui-lib";
import LoginService from "@/mixins/login";
import { mapState, mapActions } from "vuex";
import to from "await-to-js";
import WebSocketService from "@/mixins/websocket";
import {
  QueryParamService,
  UtilService,
  StorageService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Login",
  mixins: [
    IconService,
    LoginService,
    StorageService,
    WebSocketService,
    UtilService,
    QueryParamService,
  ],
  pageTitle() {
    return this.$t("login.title");
  },
  data() {
    return {
      username: "",
      password: "",
      rememberUsername: false,
      step: "username",
      error: {
        username: "",
        password: "",
        login: "",
      },
      loading: {
        login: false,
      },
    };
  },
  computed: {
    ...mapState(["loggedUser"]),
  },
  created() {
    const rememberedUsername = this.getFromStorage("username");

    if (rememberedUsername) {
      this.rememberUsername = true;
      this.username = rememberedUsername;
      this.step = "password";
    }
  },
  mounted() {
    console.log("mounted login page"); ////

    if (this.step === "username") {
      this.focusElement("usernameInput");
    } else {
      this.focusElement("passwordInput");
    }
  },
  methods: {
    ...mapActions(["setLoggedUserInStore"]),
    checkUsername() {
      this.error.username = "";

      if (!this.username.trim()) {
        this.error.username = "Username is required";
        this.focusElement("usernameInput");
      } else {
        this.step = "password";
        this.focusElement("passwordInput");
      }
    },
    async checkPassword() {
      this.clearErrors(this);

      if (!this.password.trim()) {
        this.error.password = "Password is required";
        this.focusElement("passwordInput");
      } else {
        this.loading.login = true;

        // invoke login API
        const [loginError, response] = await to(
          this.executeLogin(this.username, this.password)
        );

        this.loading.login = false;

        if (loginError) {
          this.handleLoginError(loginError);
          return;
        }

        const loginInfo = {
          username: this.username,
          token: response.data.token,
        };

        this.saveToStorage("loginInfo", loginInfo);
        this.setLoggedUserInStore(this.username);

        if (this.rememberUsername) {
          console.log("saving username to storage", this.username); ////
          this.saveToStorage("username", this.username);
        } else {
          this.deleteFromStorage("username");
        }

        const queryParams = this.getQueryParamsForCore();

        if (queryParams.redirect) {
          // redirect to initially requested URL

          console.log("queryParams.redirect", queryParams.redirect); ////

          this.$router.replace(queryParams.redirect);
        } else {
          // go to NS8 home page

          this.$router.replace("/status");
        }

        // emit login event to initialize webapp (connect ws, invoke api...)
        this.$root.$emit("login");
      }
    },
    handleLoginError(error) {
      let errorMessage = this.$t("error.generic_error");

      if (error.response) {
        switch (error.response.data.code) {
          case 401:
            errorMessage = this.$t("error.invalid_username_or_password");
            break;
        }
      } else if (error.message === "Network Error") {
        errorMessage = this.$t("error.network_error");
      }

      this.error.login = errorMessage;
      this.step = "username";
      this.password = "";
      this.focusElement("usernameInput");
    },
    goToUsername() {
      this.step = "username";
      this.username = "";
      this.password = "";
      this.focusElement("usernameInput");
      this.error.login = "";
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.login-bg {
  height: 100vh;
  background-color: $interactive-02;
}

.login-tile {
  margin-top: 25vh;
  padding-bottom: 0;
  margin-bottom: $spacing-03;
}

.login-form {
  margin-top: $spacing-07;
}

.login-footer {
  display: flex;
  height: 4rem;
  justify-content: flex-end;
  margin-top: $spacing-07;
}

.login-button {
  width: 50%;
  height: 100%;
  margin-right: -1rem;
}

.login-as {
  margin-top: $spacing-05;
}

.not-you {
  margin-left: $spacing-03;
}

.hidden {
  display: none;
}
</style>
