<template>
  <cv-modal size="default" :visible="isShown" @modal-hidden="$emit('hide')">
    <template slot="title">{{ $t("backup.backup_details") }}</template>
    <template slot="content">
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.name") }}</span>
        <span class="setting-value">{{ backup.name }}</span>
      </div>
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.repository") }}</span>
        <span class="setting-value">{{ backup.repoName }}</span>
      </div>
      <div class="mg-bottom-sm">
        <!-- //// format schedule/frequency -->
        <span class="setting-label">{{ $t("backup.frequency") }}</span>
        <span class="setting-value">{{ backup.schedule }}</span>
      </div>
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.retention") }}</span>
        <!-- //// format retention -->
        <span class="setting-value">{{ backup.retention }}</span>
      </div>
      <div class="mg-bottom-sm instances">
        <span class="setting-label">{{
          $tc("backup.instances", instances.length)
        }}</span>
        <!-- one instance -->
        <span v-if="instances.length == 1" class="setting-value">
          {{
            instances[0].ui_name
              ? instances[0].ui_name + " (" + instances[0].module_id + ")"
              : instances[0].module_id
          }}
        </span>
        <!-- multiple instances -->
        <div class="instance-list">
          <div v-for="instance in instances" :key="instance.module_id">
            {{
              instance.ui_name
                ? instance.ui_name + " (" + instance.module_id + ")"
                : instance.module_id
            }}
          </div>
        </div>
      </div>
    </template>
    <template slot="primary-button">{{ $t("common.close") }}</template>
  </cv-modal>
</template>

<script>
import { UtilService } from "@nethserver/ns8-ui-lib";

export default {
  name: "BackupDetailsModal",
  mixins: [UtilService],
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

.setting-label {
  display: inline-block;
  margin-right: $spacing-03;
  margin-bottom: $spacing-02;
  font-weight: bold;
}

.setting-value {
  word-wrap: break-word;
}

.instances {
  display: flex;
}

.instance-list {
  display: inline-block;
}
</style>
