<!--
  Copyright (C) 2024 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <div v-if="instance.module_ui_name">
      {{ `${instance.module_ui_name} (${instance.module_id})` }}
    </div>
    <div v-else-if="instance.module_id">
      {{ instance.module_id }}
    </div>
    <div v-else>
      {{ instance.path }}
    </div>
    <!-- timestamp -->
    <div class="secondary-row">
      {{ formatDate(new Date(instance.timestamp * 1000), "PPpp") }}
    </div>
    <!-- backup repository -->
    <div class="secondary-row">
      {{ `${instance.repository_name} (${instance.repository_url})` }}
    </div>
    <!-- node and/or cluster -->
    <div
      v-if="instance.node_fqdn || instance.is_generated_locally !== null"
      class="secondary-row"
    >
      <template
        v-if="instance.node_fqdn && instance.is_generated_locally !== null"
      >
        <!-- node and cluster -->
        <template v-if="instance.is_generated_locally">
          {{
            $t("backup.from_node_name_of_this_cluster", {
              name: instance.node_fqdn,
            })
          }}
        </template>
        <template v-else>
          {{
            $t("backup.from_node_name_of_different_cluster", {
              name: instance.node_fqdn,
            })
          }}
        </template>
      </template>
      <template v-else-if="instance.node_fqdn">
        <!-- node only -->
        {{
          $t("backup.from_node_name", {
            name: instance.node_fqdn,
          })
        }}
      </template>
      <template v-else>
        <!-- cluster only -->
        <template v-if="instance.is_generated_locally">
          {{ $t("backup.from_this_cluster") }}
        </template>
        <template v-else>
          {{ $t("backup.from_different_cluster") }}
        </template>
      </template>
    </div>
  </div>
</template>

<script>
import { DateTimeService } from "@nethserver/ns8-ui-lib";

export default {
  name: "InstanceToRestoreInfo",
  mixins: [DateTimeService],
  props: {
    instance: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {};
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.secondary-row {
  margin-top: $spacing-02;
  color: $ui-04;
}
</style>
