<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <div class="bx--grid bx--grid--full-width login-bg">
      <div class="bx--row">
        <div
          class="bx--offset-md-1 bx--col-md-6 bx--offset-lg-4 bx--col-lg-8 bx--offset-xlg-5 bx--col-xlg-6 bx--offset-max-6 bx--col-max-4"
        >
          <div class="test">
            <cv-tile :light="true" class="login-tile">
              <h2 class="login-title">{{ $t("login.title") }}</h2>
              <NsInlineNotification
                v-if="
                  logoutInfo.title && ['username', 'password'].includes(step)
                "
                kind="info"
                :title="logoutInfo.title"
                :description="logoutInfo.description"
                :showCloseButton="true"
                @close="clearLogoutInfo"
              />
              <div v-if="step === 'username'">
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
                      >{{ $t("login.continue") }}</cv-button
                    >
                  </div>
                </cv-form>
              </div>
              <div v-else-if="step === 'password'">
                <NsInlineNotification
                  v-if="error.login"
                  kind="error"
                  :title="$t('login.cannot_login')"
                  :description="error.login"
                  :showCloseButton="false"
                />
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
              <div v-else-if="step === '2fa'">
                <NsInlineNotification
                  v-if="error.verify2FaCode"
                  kind="error"
                  :title="$t('login.cannot_login')"
                  :description="error.verify2FaCode"
                  :showCloseButton="false"
                />
                <h4>
                  {{ $t("login.2fa_title") }}
                </h4>
                <div class="title-description mg-bottom-lg">
                  {{ $t("login.2fa_description") }}
                </div>
                <cv-form @submit.prevent="check2FaCode" class="login-form">
                  <NsTextInput
                    :label="$t('login.6_digit_code')"
                    v-model.trim="twoFaCode"
                    :invalid-message="error.twoFaCode"
                    :disabled="loading.verify2FaCode"
                    maxlength="6"
                    ref="twoFaCode"
                    class="mg-bottom-xlg"
                  />
                  <!-- setup key instructions -->
                  <cv-accordion ref="accordion">
                    <cv-accordion-item :open="toggleAccordion[0]">
                      <template slot="title">{{
                        $t("login.need_to_reconfigure_2fa")
                      }}</template>
                      <template slot="content">
                        <p class="mg-top-sm">
                          {{ $t("login.setup_key_instructions") }}
                        </p>
                      </template>
                    </cv-accordion-item>
                  </cv-accordion>
                  <div class="login-footer">
                    <NsButton
                      kind="primary"
                      :icon="ArrowRight20"
                      class="login-button"
                      :loading="loading.verify2FaCode"
                      :disabled="loading.verify2FaCode"
                      >{{ $t("login.verify_code") }}</NsButton
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
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import TwoFaService from "@/mixins/2fa";

export default {
  name: "Login",
  mixins: [
    IconService,
    LoginService,
    StorageService,
    WebSocketService,
    UtilService,
    QueryParamService,
    PageTitleService,
    TwoFaService,
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
      jwtToken: "",
      twoFaCode: "",
      error: {
        username: "",
        password: "",
        login: "",
        twoFaCode: "",
      },
      loading: {
        login: false,
        verify2FaCode: false,
      },
    };
  },
  computed: {
    ...mapState(["loggedUser", "logoutInfo"]),
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
    if (this.step === "username") {
      this.focusElement("usernameInput");
    } else {
      this.focusElement("passwordInput");
    }
  },
  methods: {
    ...mapActions(["setLoggedUserInStore", "setLogoutInfoInStore"]),
    checkUsername() {
      this.error.username = "";

      if (!this.username.trim()) {
        this.error.username = this.$t("login.username_required");
        this.focusElement("usernameInput");
      } else {
        this.step = "password";
        this.focusElement("passwordInput");
      }
    },
    async checkPassword() {
      this.clearErrors(this);

      if (!this.password.trim()) {
        this.error.password = this.$t("login.password_required");
        this.focusElement("passwordInput");
      } else {
        this.loading.login = true;

        // invoke login API
        const [loginError, response] = await to(
          this.executeLogin(this.username, this.password, "")
        );

        this.loading.login = false;

        if (loginError) {
          this.handleLoginError(loginError); // also enable OTP prompt, if 2FA is needed
          return;
        }

        if (this.rememberUsername) {
          this.saveToStorage("username", this.username);
        } else {
          this.deleteFromStorage("username");
        }
        // credentials are valid, store JWT and complete login
        this.jwtToken = response.data.token;
        this.loginSuccessful();
      }
    },
    handleLoginError(error) {
      let errorMessage = this.$t("error.generic_error");

      if (error.response?.data?.code == 401) {
        switch (error.response.data.message) {
          case "incorrect Username or Password":
            errorMessage = this.$t("error.invalid_username_or_password");
            break;
          case "missing OTP value":
            this.check2Fa();
            return;
        }
      } else if (error.message === "Network Error") {
        errorMessage = this.$t("error.network_error");
      }

      this.error.login = errorMessage;
      this.password = "";
      this.focusElement("passwordInput");
    },
    goToUsername() {
      this.step = "username";
      this.username = "";
      this.password = "";
      this.focusElement("usernameInput");
      this.error.login = "";
    },
    clearLogoutInfo() {
      const logoutInfo = {
        title: "",
        description: "",
      };
      this.setLogoutInfoInStore(logoutInfo);
    },
    check2Fa() {
      // Show the OTP input to the user
      this.step = "2fa";
      this.$nextTick(() => {
        this.focusElement("twoFaCode");
      });
    },
    loginSuccessful() {
      const loginInfo = {
        username: this.username,
        token: this.jwtToken,
      };
      this.saveToStorage("loginInfo", loginInfo);

      this.setLoggedUserInStore(this.username);
      const queryParams = this.getQueryParamsForCore();

      if (queryParams.redirect) {
        // redirect to initially requested URL
        this.$router.replace(queryParams.redirect);
      } else {
        // go to Cluster status page

        this.$router.replace("/status");
      }

      // emit login event to initialize webapp (connect ws, invoke api...)
      this.$root.$emit("login");
    },
    async check2FaCode() {
      if (!this.twoFaCode.trim()) {
        this.error.twoFaCode = this.$t("login.2fa_code_required");
        this.focusElement("twoFaCode");
        return;
      }
      this.loading.verify2FaCode = true;

      // invoke API
      const [verify2FaCodeError, response] = await to(
        this.executeLogin(this.username, this.password, this.twoFaCode)
      );
      this.loading.verify2FaCode = false;

      if (verify2FaCodeError) {
        this.handleVerify2FaError(verify2FaCodeError);
        return;
      }
      // OTP is valid, store JWT and complete login
      this.jwtToken = response.data.token;
      this.loginSuccessful();
    },
    handleVerify2FaError(error) {
      let errorMessage = this.$t("error.generic_error");

      if (
        error.response.data?.code == 401 &&
        error.response.data?.message == "OTP check failed"
      ) {
        errorMessage = this.$t("login.incorrect_code");
      } else if (error.message === "Network Error") {
        errorMessage = this.$t("error.network_error");
      }

      this.error.verify2FaCode = errorMessage;
      this.twoFaCode = "";
      this.focusElement("twoFaCode");
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

.login-title {
  margin-bottom: $spacing-06;
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

.not-you {
  margin-left: $spacing-03;
}

.hidden {
  display: none;
}
</style>
