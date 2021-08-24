<template>
  <div class="bx--grid">
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
            :title="$t('init.change_password_description')"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-form @submit.prevent="changePassword">
            <cv-text-input
              :label="$t('init.current_password')"
              v-model="currentPassword"
              :invalid-message="$t(error.currentPassword)"
              type="password"
              light
              ref="currentPassword"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('init.new_password')"
              v-model="newPassword"
              :invalid-message="$t(error.newPassword)"
              type="password"
              light
              ref="newPassword"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('init.new_password_confirm')"
              v-model="newPasswordConfirm"
              :invalid-message="$t(error.newPasswordConfirm)"
              type="password"
              light
              ref="newPasswordConfirm"
            >
            </cv-text-input>
            <NsButton kind="primary" :icon="Edit20">{{
              $t("init.change_password")
            }}</NsButton>
          </cv-form>
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
            <h6>{{ $t("init.create_new_cluster") }}</h6>
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
            <h6>{{ $t("init.join_cluster_title") }}</h6>
          </NsTile>
        </div>
      </div>
    </div>
    <div v-else-if="q.page === 'create'">
      <!-- create cluster form -->
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("init.create_new_cluster") }}</h2>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-form @submit.prevent="createCluster">
            <cv-text-input
              :label="$t('init.vpn_endpoint_address')"
              v-model.trim="vpnEndpointAddress"
              :invalid-message="$t(error.vpnEndpointAddress)"
              light
              ref="vpnEndpointAddress"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('init.vpn_endpoint_port')"
              v-model.trim="vpnEndpointPort"
              :invalid-message="$t(error.vpnEndpointPort)"
              light
              ref="vpnEndpointPort"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('init.vpn_cidr')"
              v-model.trim="vpnCidr"
              :invalid-message="$t(error.vpnCidr)"
              light
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
        </div>
      </div>
    </div>
    <div v-else-if="q.page === 'join'">
      <!-- join cluster form -->
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("init.join_cluster_title") }}</h2>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-form @submit.prevent="joinCluster">
            <cv-text-input
              :label="$t('init.cluster_url')"
              v-model.trim="clusterUrl"
              :invalid-message="$t(error.clusterUrl)"
              light
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
        </div>
      </div>
    </div>
  </div>
</template>

<script>
//// remove useless imports

import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "InitializeCluster",
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("init.welcome", { product: this.$root.config.PRODUCT_NAME });
  },
  data() {
    return {
      q: {
        page: "welcome",
      },
      // createClusterSelected: false, ////
      // joinClusterSelected: false, ////
      isPasswordChangeNeeded: false,
      currentPassword: "",
      newPassword: "",
      newPasswordConfirm: "",
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
    selectCreateCluster() {
      console.log("selectCreateCluster"); ////
      // this.createClusterSelected = true; ////
      this.$router.push("/init?page=create"); ////
      this.focusElement("vpnEndpointAddress");
    },
    selectJoinCluster() {
      console.log("selectJoinCluster"); ////
      // this.joinClusterSelected = true; ////
      this.$router.push("/init?page=join"); ////
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
      this.isPasswordChangeNeeded = false; //// remove
    },
    validatePasswordChange() {
      this.clearErrors(this);
      let isValidationOk = true;

      //// todo force password strength?

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
          this.focusElement("newPassword");
          isValidationOk = false;
        }
      } else {
        if (this.newPasswordConfirm !== this.newPassword) {
          this.error.newPassword = "error.passwords_do_not_match";
          this.error.newPasswordConfirm = "error.passwords_do_not_match";

          if (isValidationOk) {
            this.focusElement("newPasswordConfirm");
            isValidationOk = false;
          }
        }

        if (this.currentPassword === this.newPassword) {
          this.error.newPassword = "error.old_new_passwords_must_be_different";

          if (isValidationOk) {
            this.focusElement("newPassword");
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
      }

      if (!this.vpnCidr) {
        this.error.vpnCidr = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnCidr");
          isValidationOk = false;
        }
      }
    },
    createCluster() {
      console.log("createCluster"); ////

      if (!this.validateCreateCluster()) {
        return;
      }

      //// todo
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
    },
    joinCluster() {
      console.log("joinCluster"); ////

      if (!this.validateJoinCluster()) {
        return;
      }

      //// todo
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.mg-top {
  margin-top: $spacing-07;
}
</style>
