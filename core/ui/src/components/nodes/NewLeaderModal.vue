<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="goToClusterAdmin()"
  >
    <template slot="title">{{ $t("nodes.promote_node_to_leader") }}</template>
    <template slot="content">
      <NsInlineNotification
        kind="success"
        :title="$t('nodes.promotion_completed')"
        :description="
          $t('nodes.new_cluster_admin_message', {
            endpointAddress,
          })
        "
        :showCloseButton="false"
      />
    </template>
    <template slot="primary-button">{{
      $t("nodes.go_to_cluster_admin")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "NewLeaderModal",
  mixins: [UtilService, IconService],
  props: {
    isShown: Boolean,
    endpointAddress: {
      type: String,
      required: true,
    },
  },
  methods: {
    goToClusterAdmin() {
      window.location.href = `https://${this.endpointAddress}/cluster-admin/`;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
