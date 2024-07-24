<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal size="default" :visible="isShown" @modal-hidden="$emit('hide')">
    <template slot="title">{{ $t("backup.backup_details") }}</template>
    <template slot="content">
      <div class="key-value-setting">
        <span class="label">{{ $t("backup.name") }}</span>
        <span class="value">{{ backup.name }}</span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("backup.repository") }}</span>
        <span class="value">{{ backup.repoName }}</span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("backup.schedule") }}</span>
        <span class="value">
          {{ getBackupScheduleDescription(backup.schedule) }}
        </span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("backup.retention") }}</span>
        <span class="value">{{
          $tc("backup.num_backup_snapshots", backup.retention, {
            num: backup.retention,
          })
        }}</span>
      </div>
      <div class="setting instances key-value-setting">
        <span class="label">{{
          $tc("backup.Instances", instances.length)
        }}</span>
        <div class="instance-list">
          <div v-if="instances.length === 0">-</div>
          <div v-else v-for="instance in instances" :key="instance.module_id">
            <cv-interactive-tooltip
              alignment="center"
              direction="top"
              class="tooltip-with-text-trigger"
            >
              <template slot="trigger">
                <span>
                  {{
                    instance.ui_name
                      ? instance.ui_name + " (" + instance.module_id + ")"
                      : instance.module_id
                  }}
                </span>
                <NsSvg
                  v-if="instance.status && instance.status.success == true"
                  :svg="CheckmarkFilled16"
                  class="ns-success backup-status-icon"
                />
                <NsSvg
                  v-else-if="
                    instance.status && instance.status.success == false
                  "
                  :svg="ErrorFilled16"
                  class="ns-error backup-status-icon"
                />
                <NsSvg
                  v-else
                  :svg="Warning16"
                  class="ns-warning backup-status-icon"
                />
              </template>
              <template slot="content">
                <h5 class="last-backup">{{ $t("backup.last_backup") }}</h5>
                <div class="key-value-setting narrow">
                  <span class="label">
                    {{ $t("backup.status") }}
                  </span>
                  <span class="value status">
                    <span
                      v-if="instance.status && instance.status.success == true"
                      class="ns-success-dark-bg"
                    >
                      {{ $t("common.success") }}
                    </span>
                    <span
                      v-else-if="
                        instance.status && instance.status.success == false
                      "
                      class="ns-error-dark-bg"
                    >
                      {{ $t("error.error") }}
                    </span>
                    <span v-else class="ns-warning-dark-bg">
                      {{ $t("backup.backup_has_not_run_yet") }}
                    </span>
                  </span>
                </div>
                <div class="key-value-setting narrow">
                  <span class="label">
                    {{ $t("backup.completed") }}
                  </span>
                  <span class="value">
                    <span v-if="instance.status && instance.status.end">
                      {{
                        (instance.status.end * 1000)
                          | date("yyyy-MM-dd HH:mm:ss")
                      }}
                    </span>
                    <span v-else>-</span>
                  </span>
                </div>
                <div class="key-value-setting narrow">
                  <span class="label">
                    {{ $t("backup.duration") }}
                  </span>
                  <span class="value">
                    <span
                      v-if="
                        instance.status &&
                        instance.status.end &&
                        instance.status.start
                      "
                    >
                      {{
                        (instance.status.end - instance.status.start)
                          | secondsFormat
                      }}
                    </span>
                    <span v-else>-</span>
                  </span>
                </div>
                <div class="key-value-setting narrow">
                  <span class="label">
                    {{ $t("backup.total_size") }}
                  </span>
                  <span class="value">
                    <span v-if="instance.status && instance.status.total_size">
                      {{ instance.status.total_size | byteFormat }}
                    </span>
                    <span v-else>-</span>
                  </span>
                </div>
                <div class="key-value-setting narrow">
                  <span class="label">
                    {{ $t("backup.total_file_count") }}
                  </span>
                  <span class="value">
                    <span
                      v-if="instance.status && instance.status.total_file_count"
                    >
                      {{ instance.status.total_file_count }}
                    </span>
                    <span v-else>-</span>
                  </span>
                </div>
              </template>
            </cv-interactive-tooltip>
          </div>
        </div>
      </div>
    </template>
    <template slot="primary-button">{{ $t("common.close") }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "BackupDetailsModal",
  mixins: [UtilService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    backup: {
      type: Object,
      required: true,
    },
  },
  computed: {
    instances() {
      // sort function mutate the array, need to make a shallow copy. Cannot mutate backup.instances because it's a prop
      const backupInstances = [...this.backup.instances];
      return backupInstances.sort(this.sortByProperty("module_id"));
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.instances {
  display: flex;
}

.instance-list {
  display: flex;
  flex-wrap: wrap;
  column-gap: $spacing-06;
  row-gap: $spacing-04;
}

.backup-status-icon {
  margin-left: $spacing-02;
}

.last-backup {
  margin-bottom: $spacing-05;
}

.status {
  font-weight: bold;
}
</style>
