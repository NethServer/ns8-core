<template>
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
            kind="ghost"
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
    <!-- users -->
    <cv-row>
      <cv-column>
        <h4 class="mg-bottom-md">
          {{ $t("domain_users.users") }}
        </h4>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <DomainUsers
          :domain="domain"
          :groups="groups"
          @usersLoaded="onUsersLoaded"
        />
      </cv-column>
    </cv-row>
    <!-- groups -->
    <cv-row class="mg-top-lg">
      <cv-column>
        <h4
          ref="groups"
          :class="[
            'mg-bottom-md',
            { 'highlight-anchor': highlightAnchor == 'groups' },
          ]"
        >
          {{ $t("domain_users.groups") }}
        </h4>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <DomainGroups
          :domain="domain"
          :users="users"
          @groupsLoaded="onGroupsLoaded"
        />
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
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
  mixins: [TaskService, UtilService, QueryParamService, IconService],
  pageTitle() {
    return this.$t("domain_users.title");
  },
  data() {
    return {
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

      //// remove mock
      // this.domain.location = "external";

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
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
