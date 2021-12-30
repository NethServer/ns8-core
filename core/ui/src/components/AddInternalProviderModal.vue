<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    class="wizard-modal"
  >
    <template slot="title">{{ $t("domain_detail.add_provider") }}</template>
    <template slot="content">
      <template v-if="step == 'node'">
        <!-- //// disable unavailable nodes -->
        <div class="mg-bottom-md">
          {{ $t("domains.choose_node_for_account_provider_installation") }}
        </div>
        <div class="bx--grid">
          <div class="bx--row">
            <div
              v-for="(node, index) in nodes"
              :key="index"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsTile
                :light="true"
                kind="selectable"
                v-model="node.selected"
                value="nodeValue"
                :footerIcon="Chip20"
                @click="deselectOtherNodes(node)"
                :class="[
                  'min-height-card',
                  { 'disabled-node': node.unavailable },
                ]"
                :disabled="node.unavailable"
              >
                <h6>
                  {{
                    node.ui_name
                      ? node.ui_name
                      : $t("common.node") + " " + node.id
                  }}
                </h6>
                <div v-if="node.unavailable" class="mg-top">
                  {{ $t("domain_detail.used_by_another_provider") }}
                </div>
                <div v-else-if="node.ui_name" class="mg-top">
                  {{ $t("common.node") }} {{ node.id }}
                </div>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'installingProvider'">
        <NsInlineNotification
          v-if="error.addInternalProvider"
          kind="error"
          :title="$t('action.add-internal-provider')"
          :description="error.addInternalProvider"
          :showCloseButton="false"
        />
        <NsEmptyState
          :title="$t('domains.installing_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="installProviderProgress"
          :indeterminate="!installProviderProgress"
          class="mg-bottom-md"
        />
      </template>
      <template v-if="step == 'internalConfig'">
        <!-- openldap -->
        <cv-form v-if="isOpenLdap"> //// openldap config </cv-form>
        <!-- samba -->
        <template v-if="isSamba">
          <NsInlineNotification
            v-if="error.samba.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.samba.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <cv-text-input
              :label="$t('samba.adminuser')"
              v-model.trim="samba.adminuser"
              :helper-text="$t('samba.enter_samba_admin_username')"
              :invalid-message="$t(error.samba.adminuser)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="adminuser"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('samba.adminpass')"
              type="password"
              v-model="samba.adminpass"
              :helper-text="$t('samba.enter_samba_admin_password')"
              :invalid-message="$t(error.samba.adminpass)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="adminpass"
              name="adminpass"
            ></cv-text-input>
            <cv-combo-box
              v-model="samba.ipaddress"
              :options="samba.ipAddressOptions"
              auto-highlight
              :title="$t('samba.ipaddress')"
              :invalid-message="$t(error.samba.ipaddress)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              light
              ref="ipaddress"
            >
            </cv-combo-box>
            <cv-text-input
              :label="$t('samba.hostname')"
              v-model.trim="samba.hostname"
              :invalid-message="$t(error.samba.hostname)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="hostname"
            >
            </cv-text-input>
          </cv-form>
          <NsInlineNotification
            v-if="error.samba.configureModule"
            kind="error"
            :title="$t('action.configure-module')"
            :description="error.samba.configureModule"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.samba.cannotReachDc"
            kind="error"
            :title="$t('domains.cannot_connect_to_domain_controller')"
            :description="
              $t('domains.cannot_connect_to_domain_controller_description')
            "
            :showCloseButton="false"
          />
        </template>
      </template>
      <div class="wizard-buttons">
        <NsButton
          kind="secondary"
          :icon="Close20"
          @click="$emit('hide')"
          class="wizard-button"
          >{{ $t("common.cancel") }}
        </NsButton>
        <!-- <NsButton
          kind="secondary"
          :icon="ChevronLeft20"
          @click="previousStep"
          :disabled="
            isResumeConfiguration ||
            ['location', 'installingProvider', 'internalConfig'].includes(step)
          "
          class="wizard-button"
          >{{ $t("common.previous") }}
        </NsButton> -->
        <NsButton
          kind="primary"
          :icon="ChevronRight20"
          @click="nextStep"
          :disabled="isNextStepButtonDisabled()"
          :loading="loading.samba.configureModule"
          class="wizard-button"
          ref="wizardNext"
          >{{ nextStepLabel }}
        </NsButton>
      </div>
    </template>
  </cv-modal>
</template>

<script>
import {
  UtilService,
  TaskService,
  IconService,
  LottieService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

//// review all (copy/paste)

export default {
  name: "AddInternalProviderModal",
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    nodes: {
      type: Array,
      required: true,
    },
    domain: {
      type: Object,
      required: true,
    },
    isResumeConfiguration: {
      type: Boolean,
      default: false,
    },
    providerId: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      step: "node",
      installProviderProgress: 0,
      samba: {
        adminuser: "",
        adminpass: "",
        ipaddress: "",
        ipAddressOptions: [],
        hostname: "",
        nbdomain: "",
      },
      loading: {
        samba: {
          getDefaults: false,
          configureModule: false,
        },
      },
      error: {
        addInternalProvider: "",
        samba: {
          adminuser: "",
          adminpass: "",
          confirmPassword: "",
          ipaddress: "",
          hostname: "",
          nbdomain: "",
          configureModule: "",
          getDefaults: "",
          cannotReachDc: false,
        },
      },
    };
  },
  computed: {
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
    isOpenLdap() {
      return this.domain.schema == "rfc2307";
    },
    isSamba() {
      return this.domain.schema == "ad";
    },
    nextStepLabel() {
      if (this.step == "node") {
        return this.$t("domains.install_provider");
      } else if (this.step == "internalConfig") {
        return this.$t("domains.configure_provider");
      } else {
        return this.$t("common.next");
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // this.clearOpenLdapErrors(); ////
        this.clearSambaErrors();

        if (!this.isResumeConfiguration) {
          // start wizard from first step
          this.step = "node";
        } else {
          // resume configuration
          this.step = "internalConfig";

          if (this.isOpenLdap) {
            // this.getOpenLdapDefaults(); ////
          } else if (this.isSamba) {
            this.getSambaDefaults();
          }
        }

        if (this.step !== "internalConfig") {
          // set focus to next button
          setTimeout(() => {
            const element = this.$refs["wizardNext"].$el;
            element.focus();
          }, 300);
        } else {
          if (this.isOpenLdap) {
            //// focus first input field
          } else if (this.isSamba) {
            setTimeout(() => {
              this.focusElement("adminuser");
            }, 300);
          }
        }
      }
    },
    step: function () {
      if (this.step == "internalConfig") {
        if (this.isOpenLdap) {
          //// focus first input field
        } else if (this.isSamba) {
          setTimeout(() => {
            this.focusElement("adminuser");
          }, 300);
        }
      }
    },
    providerId: function () {
      this.newProviderId = this.providerId;
    },
  },
  created() {
    this.newProviderId = this.providerId;
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "node":
          this.step = "installingProvider";
          this.installProvider();
          break;
        case "internalConfig":
          if (this.isSamba) {
            this.configureSambaModule();
          } else if (this.isOpenLdap) {
            this.configureOpenLdapModule();
          }
          break;
      }
    },
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
    },
    async installProvider() {
      this.error.addInternalProvider = "";

      //// todo select version
      let version = "latest";

      const taskAction = "add-internal-provider";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.addInternalProviderCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        taskAction + "-progress",
        this.addInternalProviderProgress
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.addInternalProviderAborted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: "ghcr.io/nethserver/samba:" + version, ////
            node: parseInt(this.selectedNode.id),
            domain: this.domain.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            node: this.selectedNode.id,
            isNotificationHidden: true,
            isProgressNotified: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addInternalProvider = this.getErrorMessage(err);
        return;
      }
    },
    addInternalProviderAborted(taskResult) {
      console.log("add internal provider aborted", taskResult);

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addInternalProviderCompleted(taskContext, taskResult) {
      // unregister to task progress
      this.$root.$off("add-internal-provider-progress");

      this.step = "internalConfig";

      if (this.isSamba) {
        this.newProviderId = taskResult.output.module_id;
        this.getSambaDefaults();
      } //// else openldap

      // reload domains
      this.$emit("reloadDomains");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    addInternalProviderProgress(progress) {
      this.installProviderProgress = progress;
    },
    async getSambaDefaults() {
      this.loading.samba.getDefaults = true;
      this.error.samba.getDefaults = "";
      const taskAction = "get-defaults";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSambaDefaultsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            provision: "join-domain",
            realm: this.domain.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.samba.getDefaults = this.getErrorMessage(err);
        return;
      }
    },
    getSambaDefaultsCompleted(taskContext, taskResult) {
      console.log("getSambaDefaultsCompleted", taskResult.output); ////

      this.loading.samba.getDefaults = false;
      const defaults = taskResult.output;

      this.samba.adminuser = defaults.adminuser;
      this.samba.hostname = defaults.hostname;

      // ip address combo box
      let index = 0;
      for (const ipObject of defaults.ipaddress_list) {
        const option = {
          name: ipObject.ipaddress.replace(/\W/g, "_"),
          label: ipObject.ipaddress + " - " + ipObject.label,
          value: ipObject.ipaddress,
        };
        this.$set(this.samba.ipAddressOptions, index, option);
        index++;
      }

      // focus samba admin user
      setTimeout(() => {
        this.focusElement("adminuser");
      }, 300);
    },
    clearSambaErrors() {
      for (const key of Object.keys(this.error.samba)) {
        this.error.samba[key] = "";
      }
    },
    validateConfigureSambaModule() {
      this.clearSambaErrors();
      let isValidationOk = true;

      // samba admin user

      if (!this.samba.adminuser) {
        this.error.samba.adminuser = "common.required";

        if (isValidationOk) {
          this.focusElement("adminuser");
          isValidationOk = false;
        }
      }

      // samba admin password

      if (!this.samba.adminpass) {
        this.error.samba.adminpass = "common.required";

        if (isValidationOk) {
          this.focusElement("adminpass");
          isValidationOk = false;
        }
      }

      // samba ip address

      if (!this.samba.ipaddress) {
        this.error.samba.ipaddress = "common.required";

        if (isValidationOk) {
          this.focusElement("ipaddress");
          isValidationOk = false;
        }
      }

      // samba hostname

      if (!this.samba.hostname) {
        this.error.samba.hostname = "common.required";

        if (isValidationOk) {
          this.focusElement("hostname");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async configureSambaModule() {
      const isValidationOk = this.validateConfigureSambaModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.samba.configureModule = true;
      const taskAction = "configure-module";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.configureSambaModuleValidationFailed
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.configureSambaModuleCompleted
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.configureSambaModuleAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            adminuser: this.samba.adminuser,
            adminpass: this.samba.adminpass,
            realm: this.domain.name,
            ipaddress: this.samba.ipaddress,
            hostname: this.samba.hostname,
            provision: "join-domain",
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.samba.configureModule = this.getErrorMessage(err);
        this.loading.samba.configureModule = false;
        return;
      }
    },
    configureSambaModuleValidationFailed(validationErrors) {
      this.loading.samba.configureModule = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        if (
          [
            "realm_dc_reachable_check_failed",
            "realm_dc_avail_check_failed",
          ].includes(validationError.error)
        ) {
          this.error.samba.cannotReachDc = true;
        } else {
          const param = validationError.parameter;
          // set i18n error message
          this.error.samba[param] = "domains." + validationError.error;

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    configureSambaModuleCompleted() {
      this.loading.samba.configureModule = false;

      // hide modal
      this.$emit("hide");
    },
    configureSambaModuleAborted(taskResult) {
      console.log("configure samba module aborted", taskResult);
      this.loading.samba.configureModule = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    async configureOpenLdapModule() {
      console.log("configureOpenLdapModule"); ////
    },
    isNextStepButtonDisabled() {
      return (
        this.step == "installingProvider" ||
        (this.step == "node" && !this.selectedNode) ||
        this.loading.samba.configureModule
      );
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
.mg-top {
  margin-top: $spacing-05;
}

.row {
  margin-bottom: $spacing-03;
}

.label {
  padding-right: 0.5rem;
  font-weight: bold;
}

.bx--form {
  .bx--form-item,
  .cv-combo-box {
    margin-bottom: $spacing-06;
  }
}

.disabled-node {
  color: $disabled-02;
}
</style>
