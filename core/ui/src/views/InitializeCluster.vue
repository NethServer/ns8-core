<template>
  <div class="bx--grid">
    <cv-loading
      :active="isCreatingCluster || isJoiningCluster"
      :overlay="true"
    ></cv-loading>
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
    <div v-if="q.page === 'welcome'">
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
      <template v-if="isPasswordChangeNeeded">
        <div class="bx--row">
          <!-- password change needed -->
          <div class="bx--col-lg-16 page-title">
            <h2>{{ $t("init.create_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.create_cluster_description") }}
            </div>
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
                <NsButton
                  kind="primary"
                  :icon="Password20"
                  :disabled="isChangingPassword"
                  >{{ $t("init.change_password") }}</NsButton
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
      </template>
      <!-- admin password was changed -->
      <template v-else>
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
      </template>
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
    <div v-else-if="q.page === 'redirect'">
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("init.redirect_cluster") }}</h2>
          <div class="title-description">
            {{ $t("init.redirect_cluster_description") }}
          </div>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile light class="content-tile">
            <cv-form>
              <div>
                <a
                  target="_blank"
                  :href="this.joinEndpoint + '/cluster-admin/'"
                  >{{ $t("init.redirect_cluster_link") }}</a
                >
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
        endpoint: null,
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
      joinEndpoint: this.$route.query.endpoint
        ? "https://" + this.$route.query.endpoint
        : "",
      joinPort: "",
      joinToken: "",
      isCreatingCluster: false,
      isJoiningCluster: false,
      isChangingPassword: false,
      isMaster: true,
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
    this.retrieveClusterStatus();
    this.getDefaults();
  },
  methods: {
    ...mapActions(["setClusterInitializedInStore"]),
    async getDefaults() {
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
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          return;
        }
      }
    },
    getDefaultsCompleted(taskContext, taskResult) {
      console.log("getDefaultsCompleted", taskResult.output); ////

      this.$root.$off("get-defaults-completed");
      const defaults = taskResult.output;
      this.vpnEndpointAddress = defaults.vpn.host;
      this.vpnEndpointPort = defaults.vpn.port.toString();
      this.vpnCidr = defaults.vpn.network;
    },
    selectCreateCluster() {
      this.$router.push("/init?page=create");

      if (this.isPasswordChangeNeeded) {
        this.focusElement("currentPassword");
      } else {
        this.focusElement("vpnEndpointAddress");
      }
    },
    selectJoinCluster() {
      this.$router.push("/init?page=join");
      this.focusElement("joinCode");
    },
    async retrieveClusterStatus() {
      const taskAction = "get-cluster-status";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.getClusterStatusCompleted);

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
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          return;
        }
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      console.log("getClusterStatusCompleted", taskResult.output); ////

      this.$root.$off("get-cluster-status-completed");
      const clusterStatus = taskResult.output;

      if (clusterStatus.initialized && this.isMaster) {
        // redirect to status page
        this.setClusterInitializedInStore(true);
        this.$root.$emit("clusterInitialized");
        this.$router.replace("/status");
        return;
      }
      this.isPasswordChangeNeeded = clusterStatus.default_password;

      //// remove mock
      // this.isPasswordChangeNeeded = false; ////
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
    async changePassword() {
      this.isChangingPassword = true;

      if (!this.validatePasswordChange()) {
        this.isChangingPassword = false;
        return;
      }

      const loginInfo = this.getFromStorage("loginInfo");
      const taskAction = "change-user-password";

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$on(
        taskAction + "-completed",
        this.changeUserPasswordCompleted
      );

      // register to task validation
      this.$root.$off(taskAction + "-validation-ok");
      this.$root.$on(
        taskAction + "-validation-ok",
        this.changeUserPasswordValidationOk
      );
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$on(
        taskAction + "-validation-failed",
        this.changeUserPasswordValidationFailed
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: loginInfo.username,
            current_password: this.currentPassword,
            new_password: this.newPassword,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          return;
        }
      }
    },
    changeUserPasswordCompleted(taskContext, taskResult) {
      console.log("changeUserPasswordCompleted", taskResult.output); ////

      this.$root.$off("change-user-password-completed");
      this.isPasswordChangeNeeded = false;
    },
    changeUserPasswordValidationOk() {
      this.$root.$off("change-user-password-validation-ok");
    },
    changeUserPasswordValidationFailed(validationErrors) {
      this.$root.$off("change-user-password-validation-failed");
      this.isChangingPassword = false;

      for (const validationError of validationErrors) {
        console.log("validationError", validationError); ////

        // set i18n error message
        this.error.currentPassword = "password." + validationError.error;
        this.focusElement("currentPassword");
      }
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
      if (!this.validateCreateCluster()) {
        return;
      }

      const taskAction = "create-cluster";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.createClusterCompleted);

      // register to task error
      this.$root.$on(taskAction + "-aborted", this.createClusterAborted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            network: this.vpnCidr,
            endpoint: this.vpnEndpointAddress + ":" + this.vpnEndpointPort,
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
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
        }
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
      this.isCreatingCluster = false;
    },
    createClusterAborted(taskResult) {
      console.log("createClusterAborted", taskResult); ////

      this.$root.$off("create-cluster-aborted");
      this.isCreatingCluster = false;
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
      if (!this.validateJoinCluster()) {
        return;
      }

      const taskAction = "join-cluster";

      // register to task events
      this.$root.$once(taskAction + "-completed", this.joinClusterCompleted);
      this.$root.$once(
        taskAction + "-validation-failed",
        this.joinClusterValidationFailed
      );
      this.$root.$once(taskAction + "-aborted", this.joinClusterAborted);

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
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
          return;
        }
      }

      this.joinClusterCompleted();
    },
    joinClusterCompleted() {
      console.log("joinClusterCompleted"); ////

      // this.$root.$off("join-cluster-completed"); ////
      this.isJoiningCluster = false;
      this.$router.push("/init?page=redirect");
      //// needed?
      //this.setClusterInitializedInStore(true);
      //this.$root.$emit("clusterInitialized");
      //this.$router.replace("/status");
    },
    joinClusterAborted(taskResult) {
      console.log("joinClusterAborted", taskResult); ////

      // this.$root.$off("join-cluster-aborted"); ////
      this.isJoiningCluster = false;
    },
    joinClusterValidationFailed(validationErrors) {
      //// needed?
      console.error("validationErrors", validationErrors); //// asdf
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
  min-height: 5rem;
}
</style>
