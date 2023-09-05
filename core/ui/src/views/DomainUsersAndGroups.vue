<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/domains">{{ $t("domains.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{
                $t("domain_users.domain_name_users_and_groups", {
                  domain: domainName,
                })
              }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="subpage-title">
          <h3>
            {{
              $t("domain_users.domain_name_users_and_groups", {
                domain: domainName,
              })
            }}
          </h3>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="tertiary"
              size="field"
              :icon="Settings20"
              @click="goToDomainConfiguration()"
              class="subpage-toolbar-item"
              >{{ $t("domain_configuration.configuration") }}</NsButton
            >
          </div>
        </cv-column>
      </cv-row>
      <cv-row v-if="error.listUserDomains">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-user-domains')"
            :description="error.listUserDomains"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
            class="landscape-warning"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <NsTabs
              :container="false"
              :aria-label="$t('common.tab_navigation')"
              :noDefaultToFirst="true"
              @tab-selected="tabSelected"
            >
              <cv-tab
                id="tab-1"
                :label="$t('domain_users.users')"
                :selected="q.view === 'users'"
              >
                <DomainUsers
                  :domain="domain"
                  :groups="groups"
                  @usersLoaded="onUsersLoaded"
                />
              </cv-tab>
              <cv-tab
                id="tab-2"
                :label="$t('domain_users.groups')"
                :selected="q.view === 'groups'"
              >
                <DomainGroups
                  :domain="domain"
                  :users="users"
                  @groupsLoaded="onGroupsLoaded"
                />
              </cv-tab>
            </NsTabs>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import DomainUsers from "@/components/domains/DomainUsers";
import DomainGroups from "@/components/domains/DomainGroups";
import to from "await-to-js";

export default {
  name: "DomainUsersAndGroups",
  components: {
    DomainUsers,
    DomainGroups,
  },
  mixins: [
    TaskService,
    UtilService,
    QueryParamService,
    IconService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("domain_users.title");
  },
  data() {
    return {
      q: {
        view: "",
      },
      domainName: "",
      domain: null,
      users: [],
      groups: [],
      loading: {
        listUserDomains: false,
      },
      error: {
        listUserDomains: "",
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
    this.domainName = this.$route.params.domainName;
    this.listUserDomains();
  },
  mounted() {
    if (!this.q.view) {
      this.q.view = "users";
    }
  },
  methods: {
    goToDomainConfiguration() {
      this.$router.push({
        name: "DomainConfiguration",
        params: { domainName: this.domainName },
      });
    },
    async listUserDomains() {
      this.loading.listUserDomains = true;
      this.error.listUserDomains = "";
      const taskAction = "list-user-domains";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listUserDomainsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listUserDomainsCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listUserDomains = this.getErrorMessage(err);
        return;
      }
    },
    listUserDomainsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listUserDomains = this.$t("error.generic_error");
      this.loading.listUserDomains = false;
    },
    listUserDomainsCompleted(taskContext, taskResult) {
      let domains = taskResult.output.domains;
      this.domain = domains.find((domain) => domain.name === this.domainName);
      this.loading.listUserDomains = false;
    },
    onUsersLoaded(users) {
      this.users = users;

      // scroll to anchor if route URL contains a hash (#)
      setTimeout(() => {
        this.checkAndScrollToAnchor();
      }, 100);
    },
    onGroupsLoaded(groups) {
      this.groups = groups;
    },
    tabSelected(tabNum) {
      if (tabNum == 0) {
        this.q.view = "users";
      } else if (tabNum == 1) {
        this.q.view = "groups";
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
