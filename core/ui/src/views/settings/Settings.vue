<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column class="page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/cluster')"
          :icon="EdgeCluster32"
        >
          <h6>{{ $t("settings_cluster.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/account')"
          :icon="UserAvatar32"
        >
          <h6>{{ $t("settings_account.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/subscription')"
          :icon="Badge32"
        >
          <h6>{{ $t("settings_subscription.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/software-repository')"
          :icon="Application32"
        >
          <h6>{{ $t("settings_sw_repositories.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/http-routes')"
          :icon="Router32"
        >
          <h6>{{ $t("settings_http_routes.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/tls-certificates')"
          :icon="Certificate32"
        >
          <h6>{{ $t("settings_tls_certificates.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/smarthost')"
          :icon="Email32"
        >
          <h6>{{ $t("smarthost.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/cluster-admins')"
          :icon="UserMilitary32"
        >
          <h6>{{ $t("settings_cluster_admins.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/firewall')"
          :icon="Firewall32"
        >
          <h6>{{ $t("firewall.title") }}</h6>
        </NsTile>
      </cv-column>
      <cv-column :md="4" :xlg="4">
        <NsTile
          :light="true"
          kind="clickable"
          @click="goTo('/settings/system-logs')"
          :icon="Catalog32"
        >
          <h6>{{ $t("system_logs.title") }}</h6>
        </NsTile>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Settings",
  mixins: [UtilService, IconService, QueryParamService, PageTitleService],
  pageTitle() {
    return this.$t("settings.title");
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
  methods: {
    goTo(path) {
      this.$router.push(path);
    },
  },
};
</script>
