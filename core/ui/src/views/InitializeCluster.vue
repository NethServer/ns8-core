<template>
  <div class="bx--grid">
    <cv-loading :active="isCreatingCluster" :overlay="true"></cv-loading>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <div class="logo">
          <img
            :src="require('@/assets/logo.png')"
            :alt="this.$root.config.PRODUCT_NAME + ' logo'"
          />
        </div>
      </div>
    </div>
    <!-- //// loader while waiting api -->
    <template v-if="isPasswordChangeNeeded">
      <div class="bx--row">
        <!-- password change needed -->
        <div class="bx--col-lg-16 page-title">
          <h2>
            {{
              $t("init.welcome", { product: this.$root.config.PRODUCT_NAME })
            }}
          </h2>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <NsInlineNotification
            kind="info"
            :title="$t('init.change_admin_password')"
            :description="$t('init.change_admin_password_description')"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile light class="content-tile">
            <cv-form @submit.prevent="changePassword">
              <cv-text-input
                :label="$t('init.current_password')"
                v-model="currentPassword"
                :invalid-message="$t(error.currentPassword)"
                type="password"
                ref="currentPassword"
              >
              </cv-text-input>
              <NsPasswordInput
                :newPasswordLabel="$t('init.new_password')"
                :confirmPasswordLabel="$t('init.new_password_confirm')"
                v-model="newPassword"
                @passwordValidation="onPasswordValidation"
                :newPasswordInvalidMessage="$t(error.newPassword)"
                :confirmPasswordInvalidMessage="$t(error.confirmPassword)"
                :passwordHideLabel="$t('password.hide_password')"
                :passwordShowLabel="$t('password.show_password')"
                :lengthLabel="$t('password.long_enough')"
                :lowercaseLabel="$t('password.lowercase_letter')"
                :uppercaseLabel="$t('password.uppercase_letter')"
                :numberLabel="$t('password.number')"
                :symbolLabel="$t('password.symbol')"
                :equalLabel="$t('password.equal')"
                :focus="focusPasswordField"
              />
              <NsButton kind="primary" :icon="Password20">{{
                $t("init.change_password")
              }}</NsButton>
            </cv-form>
          </cv-tile>
        </div>
      </div>
    </template>
    <div v-else-if="q.page === 'welcome'">
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>
            {{
              $t("init.welcome", { product: this.$root.config.PRODUCT_NAME })
            }}
          </h2>
        </div>
      </div>
      <!-- create / join cluster -->
      <div class="bx--row">
        <div class="bx--col-md-4">
          <NsTile
            :light="true"
            kind="clickable"
            :icon="EdgeCluster32"
            @click="selectCreateCluster"
            large
          >
            <h6>{{ $t("init.create_cluster") }}</h6>
            <div class="tile-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </NsTile>
        </div>
        <div class="bx--col-md-4">
          <NsTile
            :light="true"
            kind="clickable"
            :icon="Connect32"
            @click="selectJoinCluster"
            large
          >
            <h6>{{ $t("init.join_cluster") }}</h6>
            <div class="tile-description">
              {{ $t("init.join_cluster_description") }}
            </div>
          </NsTile>
        </div>
      </div>
    </div>
    <div v-else-if="q.page === 'create'">
      <!-- create cluster form -->
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("init.create_cluster") }}</h2>
          <div class="title-description">
            {{ $t("init.create_cluster_description") }}
          </div>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile light class="content-tile">
            <cv-form @submit.prevent="createCluster">
              <cv-text-input
                :label="$t('init.vpn_endpoint_address')"
                v-model.trim="vpnEndpointAddress"
                :invalid-message="$t(error.vpnEndpointAddress)"
                ref="vpnEndpointAddress"
              >
              </cv-text-input>
              <cv-text-input
                :label="$t('init.vpn_endpoint_port')"
                v-model.trim="vpnEndpointPort"
                :invalid-message="$t(error.vpnEndpointPort)"
                ref="vpnEndpointPort"
              >
              </cv-text-input>
              <cv-text-input
                :label="$t('init.vpn_cidr')"
                v-model.trim="vpnCidr"
                :invalid-message="$t(error.vpnCidr)"
                ref="vpnCidr"
              >
              </cv-text-input>
              <NsButton
                kind="primary"
                :icon="EdgeCluster20"
                :disabled="isCreatingCluster"
                >{{ $t("init.create_cluster") }}</NsButton
              >
              <div>
                <cv-link @click="selectJoinCluster" class="mg-top">{{
                  $t("init.join_cluster_instead")
                }}</cv-link>
              </div>
            </cv-form>
          </cv-tile>
        </div>
      </div>
    </div>
    <div v-else-if="q.page === 'join'">
      <!-- join cluster form -->
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("init.join_cluster") }}</h2>
          <div class="title-description">
            {{ $t("init.join_cluster_description") }}
          </div>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile light class="content-tile">
            <cv-form @submit.prevent="joinCluster">
              <cv-text-area
                :label="$t('common.join_code')"
                v-model.trim="joinCode"
                :invalid-message="$t(error.joinCode)"
                :helper-text="$t('init.join_code_helper_text')"
                class="join-code"
                ref="joinCode"
              >
              </cv-text-area>
              <cv-checkbox
                :label="$t('init.tls_verify')"
                v-model="tlsVerify"
                value="checkTlsVerify"
              />
              <NsButton kind="primary" :icon="Connect20">{{
                $t("init.join_cluster")
              }}</NsButton>
              <div>
                <cv-link @click="selectCreateCluster" class="mg-top">{{
                  $t("init.create_cluster_instead")
                }}</cv-link>
              </div>
            </cv-form>
          </cv-tile>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  StorageService,
  NsPasswordInput,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";
import to from "await-to-js";

export default {
  name: "InitializeCluster",
  components: { NsPasswordInput },
  mixins: [
    UtilService,
    IconService,
    QueryParamService,
    StorageService,
    TaskService,
  ],
  pageTitle() {
    return this.$t("init.welcome", { product: this.$root.config.PRODUCT_NAME });
  },
  data() {
    return {
      q: {
        page: "welcome",
      },
      isPasswordChangeNeeded: false,
      currentPassword: "",
      newPassword: "",
      passwordValidation: false,
      focusPasswordField: { element: "" },
      vpnEndpointAddress: "",
      vpnEndpointPort: "",
      vpnCidr: "",
      joinCode: "",
      tlsVerify: true,
      joinEndpoint: "",
      joinPort: "",
      joinToken: "",
      isCreatingCluster: false,
      error: {
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
        vpnEndpointAddress: "",
        vpnEndpointPort: "",
        vpnCidr: "",
        joinCode: "",
      },
    };
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.checkPasswordChange();
    this.getDefaults();
  },
  methods: {
    ...mapActions(["setClusterInitializedInStore"]),
    async getDefaults() {
      console.log("getDefaults"); ////

      const taskAction = "get-defaults";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.getDefaultsCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getDefaultsCompleted(taskContext, taskResult) {
      console.log("getDefaultsCompleted"); ////
      console.log("defaults", taskResult.output); ////

      this.$root.$off("get-defaults-completed");
      const defaults = taskResult.output;
      this.vpnEndpointAddress = defaults.vpn.host;
      this.vpnEndpointPort = defaults.vpn.port.toString();
      this.vpnCidr = defaults.vpn.network;
    },
    selectCreateCluster() {
      this.$router.push("/init?page=create");
      this.focusElement("vpnEndpointAddress");
    },
    selectJoinCluster() {
      this.$router.push("/init?page=join");
      this.focusElement("joinCode");
    },
    async checkPasswordChange() {
      const loginInfo = this.getFromStorage("loginInfo");

      console.log("loginInfo", loginInfo); ////

      //// TEST
      // const [checkPasswordChangeError, response] = await to(
      //   this.verifyPasswordChange(loginInfo.username, loginInfo.token)
      // );

      // if (checkPasswordChangeError) {
      //   this.createErrorNotification(
      //     checkPasswordChangeError,
      //     this.$t("init.???")
      //   );
      //   return;
      // }

      //// use response
      this.isPasswordChangeNeeded = false; //// remove
    },
    onPasswordValidation(passwordValidation) {
      this.passwordValidation = passwordValidation;
    },
    validatePasswordChange() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.currentPassword) {
        this.error.currentPassword = "common.required";

        if (isValidationOk) {
          this.focusElement("currentPassword");
          isValidationOk = false;
        }
      }

      if (!this.newPassword) {
        this.error.newPassword = "common.required";

        if (isValidationOk) {
          this.focusPasswordField = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (this.currentPassword === this.newPassword) {
          if (!this.error.newPassword) {
            this.error.newPassword =
              "password.old_new_passwords_must_be_different";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (
          !this.passwordValidation.isLengthOk ||
          !this.passwordValidation.isLowercaseOk ||
          !this.passwordValidation.isUppercaseOk ||
          !this.passwordValidation.isNumberOk ||
          !this.passwordValidation.isSymbolOk
        ) {
          if (!this.error.newPassword) {
            this.error.newPassword = "password.password_not_secure";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!this.passwordValidation.isEqualOk) {
          if (!this.error.newPassword) {
            this.error.newPassword = "password.passwords_do_not_match";
          }

          if (!this.error.confirmPassword) {
            this.error.confirmPassword = "password.passwords_do_not_match";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "confirmPassword" };
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
    },
    changePassword() {
      if (!this.validatePasswordChange()) {
        return;
      }

      //// todo call api

      const loginInfo = this.getFromStorage("loginInfo");

      console.log("loginInfo", loginInfo); ////

      //// TEST
      // const [changePasswordError, response] = await to(
      //   this.changePasswordApi(loginInfo.username, loginInfo.token, this.currentPassword, this.newPassword)
      // );

      // if (changePasswordError) {
      //   this.createErrorNotification(
      //     changePasswordError,
      //     this.$t("init.???")
      //   );
      //   return;
      // }

      //// use response
      this.isPasswordChangeNeeded = false; ////
    },
    validateCreateCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.vpnEndpointAddress) {
        this.error.vpnEndpointAddress = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnEndpointAddress");
          isValidationOk = false;
        }
      }

      if (!this.vpnEndpointPort) {
        this.error.vpnEndpointPort = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnEndpointPort");
          isValidationOk = false;
        }
      } else {
        const vpnEndpointPortNumber = Number(this.vpnEndpointPort);

        if (
          !(
            Number.isInteger(vpnEndpointPortNumber) && vpnEndpointPortNumber > 0
          )
        ) {
          this.error.vpnEndpointPort = "error.invalid_port_number";

          if (isValidationOk) {
            this.focusElement("vpnEndpointPort");
            isValidationOk = false;
          }
        }
      }

      if (!this.vpnCidr) {
        this.error.vpnCidr = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnCidr");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async createCluster() {
      console.log("createCluster"); ////

      if (!this.validateCreateCluster()) {
        return;
      }

      const taskAction = "create-cluster";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.createClusterCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            network: this.vpnCidr,
            endpoint: this.vpnEndpointAddress,
            listen_port: parseInt(this.vpnEndpointPort),
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
      this.isCreatingCluster = true;
    },
    createClusterCompleted() {
      console.log("createClusterCompleted"); ////

      this.$root.$off("create-cluster-completed");
      this.setClusterInitializedInStore(true);
      this.$root.$emit("clusterInitialized");
      this.$router.replace("/status");
    },
    validateJoinCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.joinCode) {
        this.error.joinCode = "common.required";

        if (isValidationOk) {
          this.focusElement("joinCode");
          isValidationOk = false;
        }
      } else {
        let decoded;

        try {
          decoded = atob(this.joinCode);
        } catch (DOMException) {
          this.error.joinCode = "init.invalid_join_code";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        }

        if (!decoded) {
          this.error.joinCode = "init.invalid_join_code";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        } else {
          console.log("decoded", decoded); ////

          let [endpoint, port, token] = decoded.split("|");

          if (!(endpoint && port && token)) {
            this.error.joinCode = "init.invalid_join_code";

            if (isValidationOk) {
              this.focusElement("joinCode");
              isValidationOk = false;
            }
          } else {
            console.log("endpoint", endpoint); ////
            console.log("port", port); ////
            console.log("token length", token.length); ////

            this.joinEndpoint = endpoint;
            this.joinPort = port;
            this.joinToken = token;
          }
        }
      }

      return isValidationOk;
    },
    async joinCluster() {
      console.log("joinCluster"); ////

      if (!this.validateJoinCluster()) {
        return;
      }

      console.log("join code ok, tls verify", this.tlsVerify); ////

      const taskAction = "join-cluster";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.joinClusterCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            url: this.joinEndpoint,
            jwt: this.joinToken,
            listen_port: parseInt(this.joinPort),
            tls_verify: this.tlsVerify,
          },
          extra: {
            title: this.$t("//// join"),
            description: this.$t("////"),
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }

      // this.setClusterInitializedInStore(true); ////
      // this.$router.replace("/status"); ////
      // console.log("done"); ////
    },
    joinClusterCompleted() {
      console.log("joinClusterCompleted"); ////

      this.$root.$off("join-cluster-completed");
      this.setClusterInitializedInStore(true);
      this.$root.$emit("clusterInitialized");
      this.$router.replace("/status");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.logo {
  width: 4rem;
  height: 4rem;
  margin-top: $spacing-07;
  flex-shrink: 0;
}

.logo img {
  width: 100%;
  height: 100%;
}

.mg-top {
  margin-top: $spacing-07;
}

.tile-description {
  margin-top: $spacing-03;
  color: $text-02;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.join-code textarea {
  min-height: 7rem;
}
</style>
