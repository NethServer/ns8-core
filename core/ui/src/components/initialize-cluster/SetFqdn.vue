<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-row>
      <cv-column class="welcome">
        <div>
          <h2>
            {{
              action === "create"
                ? $t("init.create_cluster")
                : $t("init.join_cluster")
            }}
          </h2>
          <div class="title-description">
            {{
              action === "create"
                ? $t("init.create_cluster_description")
                : $t("init.join_cluster_description")
            }}
          </div>
        </div>
        <WelcomeLogo />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <NsInlineNotification
          kind="info"
          :title="
            $t(
              action === 'create'
                ? 'init.review_leader_fqdn_title'
                : 'init.review_worker_fqdn_title'
            )
          "
          :description="
            $t(
              action === 'create'
                ? 'init.review_leader_fqdn_description'
                : 'init.review_worker_fqdn_description'
            )
          "
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <cv-tile light>
          <NsInlineNotification
            v-if="error.getFqdn"
            kind="error"
            :title="$t('action.get-fqdn')"
            :description="error.getFqdn"
            :showCloseButton="false"
          />
          <cv-form v-else @submit.prevent="setFqdn">
            <cv-skeleton-text
              v-if="loading.getFqdn"
              :paragraph="true"
              :line-count="5"
              heading
            ></cv-skeleton-text>
            <template v-else>
              <cv-text-input
                :label="
                  action === 'create'
                    ? $t('init.leader_node_hostname')
                    : $t('init.worker_node_hostname')
                "
                v-model.trim="hostname"
                :invalid-message="error.hostname"
                :disabled="loading.getFqdn || loading.setFqdn"
                ref="hostname"
              >
              </cv-text-input>
              <cv-text-input
                :label="
                  action === 'create'
                    ? $t('init.leader_node_domain')
                    : $t('init.worker_node_domain')
                "
                v-model.trim="domain"
                placeholder="example.org"
                :invalid-message="error.domain"
                :disabled="loading.getFqdn || loading.setFqdn"
                ref="domain"
              >
              </cv-text-input>
            </template>
            <cv-button-set class="footer-buttons">
              <NsButton
                type="button"
                kind="secondary"
                :icon="ChevronLeft20"
                size="lg"
                @click="goBack"
                >{{ $t("common.go_back") }}
              </NsButton>
              <NsButton
                kind="primary"
                :icon="ChevronRight20"
                :loading="loading.setFqdn"
                :disabled="loading.getFqdn || loading.setFqdn"
                size="lg"
                >{{ $t("init.set_fqdn") }}
              </NsButton>
            </cv-button-set>
          </cv-form>
          <NsInlineNotification
            v-if="error.setFqdn"
            kind="error"
            :title="$t('action.set-fqdn')"
            :description="error.setFqdn"
            :showCloseButton="false"
          />
        </cv-tile>
      </cv-column>
    </cv-row>
  </div>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import WelcomeLogo from "@/components/initialize-cluster/WelcomeLogo";

export default {
  name: "SetFqdn",
  components: { WelcomeLogo },
  mixins: [UtilService, IconService, TaskService],
  props: {
    action: String,
  },
  data() {
    return {
      hostname: "",
      domain: "",
      loading: {
        getFqdn: false,
        setFqdn: false,
      },
      error: {
        getFqdn: "",
        setFqdn: "",
        hostname: "",
        domain: "",
      },
    };
  },
  created() {
    this.getFqdn();
  },
  methods: {
    async getFqdn() {
      this.error.getFqdn = "";
      this.loading.getFqdn = true;
      const taskAction = "get-fqdn";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.getFqdnAborted);

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getFqdnCompleted
      );

      const res = await to(
        this.createNodeTask(1, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
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
    },
    getFqdnCompleted(taskContext, taskResult) {
      this.hostname = taskResult.output.hostname;
      this.domain = taskResult.output.domain;
      this.loading.getFqdn = false;
      this.focusElement("hostname");
    },
    getFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getFqdn = false;
    },
    validateSetFqdn() {
      this.error.hostname = "";
      this.error.domain = "";
      let isValidationOk = true;

      // hostname

      if (!this.hostname) {
        this.error.hostname = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("hostname");
          isValidationOk = false;
        }
      }

      // domain

      if (!this.domain) {
        this.error.domain = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("domain");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async setFqdn() {
      if (!this.validateSetFqdn()) {
        return;
      }
      this.loading.setFqdn = true;
      this.error.setFqdn = "";
      const taskAction = "set-fqdn";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.setFqdnAborted);

      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setFqdnValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setFqdnCompleted
      );

      const res = await to(
        this.createNodeTask(1, {
          action: taskAction,
          data: {
            hostname: this.hostname,
            domain: this.domain,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setFqdn = this.getErrorMessage(err);
        this.loading.setFqdn = false;
        return;
      }
    },
    setFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setFqdn = false;
    },
    setFqdnValidationFailed(validationErrors) {
      this.loading.setFqdn = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "init." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setFqdnCompleted() {
      this.loading.setFqdn = false;
      this.$emit("fqdnSet");
    },
    goBack() {
      this.$emit("goBack");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
