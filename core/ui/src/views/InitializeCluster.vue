<template>
  <div class="bx--grid">
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
                :label="$t('init.new_password')"
                v-model="newPassword"
                @passwordValidation="onPasswordValidation"
                :invalid-message="$t(error.newPassword)"
                :passwordHideLabel="$t('password.hide_password')"
                :passwordShowLabel="$t('password.show_password')"
                :lengthLabel="$t('password.long_enough')"
                :lowercaseLabel="$t('password.lowercase_letter')"
                :uppercaseLabel="$t('password.uppercase_letter')"
                :numberLabel="$t('password.number')"
                :symbolLabel="$t('password.symbol')"
                ref="newPassword"
              />
              <cv-text-input
                :label="$t('init.new_password_confirm')"
                v-model="newPasswordConfirm"
                :invalid-message="$t(error.newPasswordConfirm)"
                type="password"
                ref="newPasswordConfirm"
              >
              </cv-text-input>
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
              <NsButton kind="primary" :icon="EdgeCluster20">{{
                $t("init.create_cluster")
              }}</NsButton>
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
              <cv-text-input
                :label="$t('init.cluster_url')"
                v-model.trim="clusterUrl"
                :invalid-message="$t(error.clusterUrl)"
                ref="clusterUrl"
              >
              </cv-text-input>
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
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";

export default {
  name: "InitializeCluster",
  components: { NsPasswordInput },
  mixins: [UtilService, IconService, QueryParamService, StorageService],
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
      newPasswordConfirm: "",
      isPasswordValid: false,
      vpnEndpointAddress: "",
      vpnEndpointPort: "",
      vpnCidr: "",
      clusterUrl: "",
      error: {
        currentPassword: "",
        newPassword: "",
        newPasswordConfirm: "",
        vpnEndpointAddress: "",
        vpnEndpointPort: "",
        vpnCidr: "",
        clusterUrl: "",
      },
    };
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      console.log("beforeRouteEnter", to, from); ////
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    console.log("beforeRouteUpdate", to, from); ////
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.checkPasswordChange();
  },
  methods: {
    ...mapActions(["setClusterInitializedInStore"]),
    selectCreateCluster() {
      this.$router.push("/init?page=create");
      this.focusElement("vpnEndpointAddress");
    },
    selectJoinCluster() {
      this.$router.push("/init?page=join");
      this.focusElement("clusterUrl");
    },
    async checkPasswordChange() {
      const loginInfo = this.getFromStorage("loginInfo");

      console.log("loginInfo", loginInfo); ////

      //// TEST
      // const [checkPasswordChangeError, response] = await to(
      //   this.verifyPasswordChange(loginInfo.username, loginInfo.token)
      // );

      // if (checkPasswordChangeError) {
      //   this.createTaskErrorNotification(
      //     checkPasswordChangeError,
      //     this.$t("init.???")
      //   );
      //   return;
      // }

      //// use response
      this.isPasswordChangeNeeded = true; //// remove
    },
    onPasswordValidation(isPasswordValid) {
      this.isPasswordValid = isPasswordValid;
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
          this.focusNewPassword();
          isValidationOk = false;
        }
      } else {
        if (this.newPasswordConfirm !== this.newPassword) {
          this.error.newPassword = "password.passwords_do_not_match";
          this.error.newPasswordConfirm = "password.passwords_do_not_match";

          if (isValidationOk) {
            this.focusNewPassword();
            isValidationOk = false;
          }
        }

        if (this.currentPassword === this.newPassword) {
          this.error.newPassword =
            "password.old_new_passwords_must_be_different";

          if (isValidationOk) {
            this.focusNewPassword();
            isValidationOk = false;
          }
        }

        if (!this.isPasswordValid) {
          this.error.newPassword = "password.password_not_secure";

          if (isValidationOk) {
            this.focusNewPassword();
            isValidationOk = false;
          }
        }
      }

      if (!this.newPasswordConfirm) {
        this.error.newPasswordConfirm = "common.required";

        if (isValidationOk) {
          this.focusElement("newPasswordConfirm");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    focusNewPassword() {
      // NsPasswordInput has a <div> as root element, need to focus nested <input> element
      this.$nextTick(() => {
        const element = this.$refs.newPassword.$el
          .getElementsByClassName("password-input")[0]
          .getElementsByClassName("bx--password-input")[0];
        element.focus();
      });
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
      //   this.createTaskErrorNotification(
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
    createCluster() {
      console.log("createCluster"); ////

      if (!this.validateCreateCluster()) {
        return;
      }

      //// todo

      this.setClusterInitializedInStore(true);
      this.$router.replace("/status");
      console.log("done"); ////
    },
    validateJoinCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.clusterUrl) {
        this.error.clusterUrl = "common.required";

        if (isValidationOk) {
          this.focusElement("clusterUrl");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    joinCluster() {
      console.log("joinCluster"); ////

      if (!this.validateJoinCluster()) {
        return;
      }

      //// todo

      this.setClusterInitializedInStore(true);
      this.$router.replace("/status");
      console.log("done"); ////
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
