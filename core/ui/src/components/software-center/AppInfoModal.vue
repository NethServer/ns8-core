<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div v-if="app">
    <NsModal
      size="default"
      :visible="isShown"
      @modal-hidden="$emit('close')"
      class="no-pad-modal"
    >
      <template slot="title">{{ $t("software_center.app_info") }}</template>
      <!-- v-if="isShown" is needed to reset modal scroll to top -->
      <template v-if="isShown" slot="content">
        <div class="mg-bottom-md">
          <div class="logo-and-name">
            <div class="app-logo">
              <img
                :src="
                  app.logo
                    ? app.logo
                    : require('@/assets/module_default_logo.png')
                "
                :alt="app.name + ' logo'"
              />
            </div>
            <div class="app-name">
              <h3 data-modal-primary-focus>{{ app.name || app.id }}</h3>
            </div>
          </div>
        </div>
        <div v-if="app.screenshots.length" class="mg-bottom-md">
          <div class="screenshots">
            <ImageGallery :fullScreen="true" :images="app.screenshots" />
          </div>
        </div>
        <div class="description">
          {{ getApplicationDescription(app) }}
        </div>
        <!-- warning for rootfull app -->
        <NsInlineNotification
          v-if="app.rootfull && app.certification_level < 3"
          kind="warning"
          :title="$t('software_center.rootfull_app_warning_title')"
          :description="
            $t('software_center.rootfull_app_warning_description', {
              appName: app.name,
            })
          "
          :showCloseButton="false"
        />
        <div class="key-value-setting">
          <span class="label">
            {{ $t("software_center.certification") }}
          </span>
          <span class="value">
            <CertificationLevelBadge :level="app.certification_level" />
          </span>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">{{
              $tc("software_center.categories", app.categories.length)
            }}</span>
            <span class="value">
              {{ getApplicationCategories(app) }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">
              {{ $t("software_center.latest_version") }}
            </span>
            <span class="value">
              {{
                app.versions.length
                  ? app.versions[0].tag
                  : $t("common.not_available")
              }}
              <span v-if="app.upstream_name"> ({{ app.upstream_name }}) </span>
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">{{ $t("software_center.repository") }}</span>
            <span class="value">
              {{ app.repository || "-" }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <span class="label"
            >{{ $tc("software_center.authors", app.authors.length) }}
          </span>
          <span class="value">
            <span v-if="app.authors.length == 1">
              <cv-link :href="'mailto:' + app.authors[0].email" target="_blank">
                {{ app.authors[0].name }}
              </cv-link>
            </span>
            <ul v-else>
              <li v-for="(author, index) in app.authors" :key="index">
                <cv-link :href="'mailto:' + author.email" target="_blank">
                  {{ author.name }}
                </cv-link>
              </li>
            </ul>
          </span>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">
              {{ $t("software_center.images_label") }}
            </span>
            <span class="value">
              <ul class="image-list unordered-list">
                <li>
                  {{ app.source }}:{{
                    app.versions.length ? app.versions[0].tag : "-"
                  }}
                </li>
                <li v-for="(image, index) in appImages" :key="index">
                  {{ image }}
                </li>
              </ul>
            </span>
          </div>
        </div>
        <div v-if="appPermissions.length" class="key-value-setting">
          <div>
            <span class="label">
              {{ $t("software_center.permissions_label") }}
            </span>
          </div>
          <NsDataTable
            class="permission-table"
            :allRows="appPermissions"
            :columns="permissionColumnLabels"
            :rawColumns="permissionColumns"
            :sortable="true"
            :pageSizes="[10, 25, 50, 100]"
            :itemsPerPageLabel="$t('pagination.items_per_page')"
            :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
            :ofTotalPagesLabel="$t('pagination.of_total_pages')"
            :backwardText="$t('pagination.previous_page')"
            :forwardText="$t('pagination.next_page')"
            :pageNumberLabel="$t('pagination.page_number')"
            @updatePage="permissionTablePage = $event"
          >
            <template slot="data">
              <cv-data-table-row
                v-for="perm in permissionTablePage"
                :key="perm.raw"
                :value="perm.raw"
              >
                <cv-data-table-cell>{{ perm.description }}</cv-data-table-cell>
                <cv-data-table-cell>{{ perm.target }}</cv-data-table-cell>
                <cv-data-table-cell>{{ perm.raw }}</cv-data-table-cell>
              </cv-data-table-row>
            </template>
          </NsDataTable>
        </div>
        <div class="info-divider"></div>
        <div class="key-value-setting">
          <div>
            <span>
              <cv-link :href="app.docs.documentation_url" target="_blank">
                {{ $t("software_center.documentation") }}
              </cv-link>
            </span>
            &bull;
            <span>
              <cv-link :href="app.docs.bug_url" target="_blank">
                {{ $t("software_center.bugs") }}
              </cv-link>
            </span>
            &bull;
            <span>
              <cv-link :href="app.docs.code_url" target="_blank">
                {{ $t("software_center.source_code") }}
              </cv-link>
            </span>
            <template v-if="app.docs.terms_url">
              &bull;
              <span>
                <cv-link :href="app.docs.terms_url" target="_blank">
                  {{ $t("common.terms_and_conditions") }}
                </cv-link>
              </span>
            </template>
            <template v-if="app.docs.relnotes_url">
              &bull;
              <span>
                <cv-link :href="app.docs.relnotes_url" target="_blank">
                  {{ $t("common.release_notes") }}
                </cv-link>
              </span>
            </template>
          </div>
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </NsModal>
  </div>
</template>

<script>
import { UtilService } from "@nethserver/ns8-ui-lib";
import ImageGallery from "./ImageGallery";
import CertificationLevelBadge from "./CertificationLevelBadge.vue";

export default {
  name: "AppInfoModal",
  components: { ImageGallery, CertificationLevelBadge },
  mixins: [UtilService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  data() {
    return {
      permissionColumns: ["description", "target", "raw"],
      permissionTablePage: [],
    };
  },
  computed: {
    permissionColumnLabels() {
      return [
        this.$t("software_center.permission_description"),
        this.$t("software_center.permission_target"),
        this.$t("software_center.permission_label"),
      ];
    },
    appImages() {
      if (
        this.app &&
        this.app.versions.length &&
        this.app.versions[0]["labels"]["org.nethserver.images"]
      ) {
        return this.app.versions[0]["labels"]["org.nethserver.images"]
          .trim()
          .replace(/\s+/g, " ")
          .split(" ");
      } else {
        return [];
      }
    },
    appPermissions() {
      if (!this.app) {
        return [];
      }
      const permissions = [];

      // a rootfull app is not limited by its roles, so this goes first
      if (this.app.rootfull) {
        permissions.push({
          raw: "org.nethserver.rootfull",
          description: this.$t("software_center.app_permission_rootfull"),
          target: this.$t("software_center.app_permission_targets.node"),
        });
      }
      if (!this.app.versions.length) {
        return permissions;
      }
      const label =
        this.app.versions[0]["labels"]["org.nethserver.authorizations"] || "";
      const seen = new Set();

      for (const entry of label.trim().split(/\s+/)) {
        const sep = entry.indexOf(":");
        if (sep < 1) {
          continue;
        }
        const selector = entry.slice(0, sep);
        const at = selector.indexOf("@");
        const module = at < 0 ? selector : selector.slice(0, at);
        const scope = at < 0 ? "" : selector.slice(at + 1);

        // core grants self:selfadm to every app, it gives nothing on other components
        if (module === "self") {
          continue;
        }
        for (const role of entry.slice(sep + 1).split(",")) {
          const raw = `${selector}:${role}`;
          if (!role || seen.has(raw)) {
            continue;
          }
          seen.add(raw);
          permissions.push(this.describePermission(module, scope, role, raw));
        }
      }
      return permissions;
    },
  },
  methods: {
    getApplicationDescription(app) {
      return this.getAppDescription(app, this);
    },
    getApplicationCategories(app) {
      return this.getAppCategories(app, this);
    },
    describePermission(module, scope, role, raw) {
      // a bare number selects a node by its ID, like "1:fwadm"
      const nodeId = !scope && /^\d+$/.test(module) ? module : "";
      const agentType = nodeId ? "node" : module;
      // without a scope, only node and cluster resolve to an agent
      const resolved = scope || agentType === "node" || agentType === "cluster";
      const moduleName = module.charAt(0).toUpperCase() + module.slice(1);
      const pairKey = `software_center.app_permissions.${agentType}_${role}`;
      const coreRoleKey = `software_center.app_permission_core_roles.${role}`;
      let description;

      if (resolved && this.$te(pairKey)) {
        description = this.$t(pairKey);
      } else if (resolved && this.$te(coreRoleKey)) {
        // roles like "reader" are created by core on every module
        description = this.$t(coreRoleKey, { module: moduleName });
      } else {
        description = this.$t("software_center.app_permission_generic", {
          role,
          module: moduleName,
        });
      }

      let target = moduleName;
      if (nodeId) {
        target = this.$t("software_center.app_permission_targets.node_id", {
          id: nodeId,
        });
      } else {
        const targetKey = scope
          ? `software_center.app_permission_scopes.${scope}`
          : `software_center.app_permission_targets.${module}`;
        if (this.$te(targetKey)) {
          target = this.$t(targetKey, { module: moduleName });
        }
      }
      return { raw, description, target };
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.logo-and-name {
  display: flex;
  align-items: center;
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  margin-right: $spacing-05;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.description {
  margin-bottom: $spacing-06;
}

.info-divider {
  margin-top: $spacing-02;
  margin-bottom: $spacing-03;
  border-bottom: 1px solid $ui-04;
}

.image-list {
  margin-left: $spacing-03;
}

.permission-table {
  margin-top: $spacing-03;
}
</style>
