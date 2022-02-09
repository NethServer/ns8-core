<template>
  <cv-grid>
    <cv-loading
      :active="isCreatingCluster || isJoiningCluster"
      :overlay="true"
    ></cv-loading>
    <cv-row>
      <cv-column>
        <div class="logo">
          <img
            :src="require('@/assets/logo.png')"
            :alt="this.$root.config.PRODUCT_NAME + ' logo'"
          />
        </div>
      </cv-column>
    </cv-row>
    <template v-if="q.page === 'welcome'">
      <cv-row>
        <cv-column class="welcome">
          <h2>
            {{
              $t("init.welcome", { product: this.$root.config.PRODUCT_NAME })
            }}
          </h2>
        </cv-column>
      </cv-row>
      <!-- create / join cluster -->
      <cv-row>
        <cv-column :md="4">
          <NsTile
            :light="true"
            kind="clickable"
            :icon="EdgeCluster32"
            @click="selectCreateCluster"
            large
          >
            <h6 class="mg-bottom-sm">{{ $t("init.create_cluster") }}</h6>
            <div class="tile-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </NsTile>
        </cv-column>
        <cv-column :md="4">
          <NsTile
            :light="true"
            kind="clickable"
            :icon="Connect32"
            @click="selectJoinCluster"
            large
          >
            <h6 class="mg-bottom-sm">{{ $t("init.join_cluster") }}</h6>
            <div class="tile-description">
              {{ $t("init.join_cluster_description") }}
            </div>
          </NsTile>
        </cv-column>
      </cv-row>
    </template>
    <template v-else-if="q.page === 'create'">
      <template v-if="isPasswordChangeNeeded">
        <cv-row>
          <!-- password change needed -->
          <cv-column class="welcome">
            <h2>{{ $t("init.create_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsInlineNotification
              kind="info"
              :title="$t('init.change_admin_password')"
              :description="$t('init.change_admin_password_description')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <cv-tile light>
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
                  <cv-link @click="selectJoinCluster" class="mg-top-lg">{{
                    $t("init.join_cluster_instead")
                  }}</cv-link>
                </div>
              </cv-form>
            </cv-tile>
          </cv-column>
        </cv-row>
      </template>
      <!-- admin password was changed -->
      <template v-else>
        <!-- create cluster form -->
        <cv-row>
          <cv-column class="welcome">
            <h2>{{ $t("init.create_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <cv-tile light>
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
                  <cv-link @click="selectJoinCluster" class="mg-top-lg">{{
                    $t("init.join_cluster_instead")
                  }}</cv-link>
                </div>
              </cv-form>
            </cv-tile>
          </cv-column>
        </cv-row>
      </template>
    </template>
    <template v-else-if="q.page === 'join'">
      <!-- join cluster form -->
      <cv-row>
        <cv-column class="welcome">
          <h2>{{ $t("init.join_cluster") }}</h2>
          <div class="title-description">
            {{ $t("init.join_cluster_description") }}
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-form @submit.prevent="joinCluster">
              <cv-text-area
                :label="$t('common.join_code')"
                v-model.trim="joinCode"
                :invalid-message="$t(error.joinCode)"
                :helper-text="
                  $t('init.join_code_helper_text') +
                  ' https://LEADER_NODE_IP/cluster-admin/#/nodes?isShownAddNodeModal=true'
                "
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
                <cv-link @click="selectCreateCluster" class="mg-top-lg">{{
                  $t("init.create_cluster_instead")
                }}</cv-link>
              </div>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
    </template>
    <template v-else-if="q.page === 'redirect'">
      <cv-row>
        <cv-column class="welcome">
          <h2>{{ $t("init.redirect_cluster") }}</h2>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <div class="title-description">
            {{ $t("init.redirect_cluster_description") }}
          </div>
        </cv-column>
      </cv-row>
      <cv-row class="mg-top-lg">
        <cv-column>
          <a
            :href="this.joinEndpoint + '/cluster-admin/'"
            class="external-link-button"
          >
            <NsButton kind="primary" :icon="Launch20">
              {{ $t("init.redirect_cluster_link") }}
            </NsButton>
          </a>
        </cv-column>
      </cv-row>
    </template>
  </cv-grid>
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
import NotificationService from "@/mixins/notification";

export default {
  name: "InitializeCluster",
  components: { NsPasswordInput },
  mixins: [
    UtilService,
    IconService,
    QueryParamService,
    StorageService,
    TaskService,
    NotificationService,
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
      passwordValidation: null,
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
      this.$root.$once(taskAction + "-completed", this.getDefaultsCompleted);

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
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    getDefaultsCompleted(taskContext, taskResult) {
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
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
      );

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
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;

      //// remove mock
      // clusterStatus.initialized = false; ////

      if (clusterStatus.initialized && this.isMaster) {
        // redirect to status page
        this.setClusterInitializedInStore(true);
        this.$root.$emit("clusterInitialized");
        this.$router.replace("/status");
        return;
      }
      this.isPasswordChangeNeeded = clusterStatus.default_password;
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
      this.$root.$once(
        taskAction + "-completed",
        this.changeUserPasswordCompleted
      );

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
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
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    changeUserPasswordCompleted() {
      this.isPasswordChangeNeeded = false;
    },
    changeUserPasswordValidationFailed(validationErrors) {
      this.isChangingPassword = false;

      for (const validationError of validationErrors) {
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
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.createClusterCompleted);

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.createClusterAborted);

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
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
        }
        return;
      }
      this.isCreatingCluster = true;
    },
    createClusterCompleted() {
      this.setClusterInitializedInStore(true);
      this.$root.$emit("clusterInitialized");
      this.$router.replace("/status");
      this.isCreatingCluster = false;
    },
    createClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
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
          let [endpoint, port, token] = decoded.split("|");

          if (!(endpoint && port && token)) {
            this.error.joinCode = "init.invalid_join_code";

            if (isValidationOk) {
              this.focusElement("joinCode");
              isValidationOk = false;
            }
          } else {
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
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.joinClusterCompleted);

      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.joinClusterValidationFailed
      );

      this.$root.$off(taskAction + "-aborted");
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
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
      this.isJoiningCluster = true;
    },
    joinClusterCompleted(taskContext, taskResult) {
      console.log("joinClusterCompleted"); ////
      console.log("taskContext", taskContext); ////
      console.log("taskResult", taskResult); ////

      this.isJoiningCluster = false;
      this.$router.push("/init?page=redirect");

      const nodeId = taskResult.output.nodeId;

      console.log("nodeId", nodeId); ////

      const notification = {
        title: this.$t("action.join-cluster"),
        description: this.$t("init.node_added_to_cluster", { nodeId }),
        type: "success",
      };
      this.createNotification(notification);
    },
    joinClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.isJoiningCluster = false;
    },
    joinClusterValidationFailed(validationErrors) {
      console.error("validation failed", validationErrors);
      this.isJoiningCluster = false;
      this.error.joinCode = "init.invalid_join_code";
      this.focusElement("joinCode");
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

.welcome {
  margin-top: 2rem;
  margin-bottom: 4rem;
}

.tile-description {
  color: $text-02;
}

.bx--form .bx--form-item {
  margin-bottom: $spacing-06;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.join-code textarea {
  min-height: 5rem;
}
</style>
