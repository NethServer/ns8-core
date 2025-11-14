<template>
  <NsModal
    size="default"
    :visible="visible"
    @modal-hidden="$emit('hide')"
    @primary-click="$emit('primary-click')"
  >
    <template slot="title">
      {{
        $t("software_center.app_restart", {
          app: instanceToRestart ? instanceToRestart.id : "",
        })
      }}
    </template>
    <template slot="content">
      <NsInlineNotification
        kind="warning"
        :title="$t('common.please_read_carefully')"
        :description="$t('software_center.restart_module_warning')"
        :showCloseButton="false"
      />
      <div>
        {{
          $t("software_center.restart_app", {
            name: instanceToRestart ? instanceToRestart.id : "",
          })
        }}
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("software_center.restart_instance") }}</template>
  </NsModal>
</template>

<script>
export default {
  name: "RestartModuleModal",
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    instanceToRestart: {
      type: Object,
      default: null
    }
  }
};
</script>
