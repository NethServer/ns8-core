<template>
  <NsModal
    size="default"
    :visible="visible"
    @modal-hidden="$emit('hide')"
    @primary-click="$emit('primary-click')"
  >
    <template slot="title">
      {{
        $t("applications.app_restart", {
          app: instanceToRestart ? instanceToRestart.id : "",
        })
      }}
    </template>
    <template slot="content">
      <NsInlineNotification
        kind="warning"
        :title="$t('common.please_read_carefully')"
        :description="$t('applications.restart_module_warning')"
        :showCloseButton="false"
      />
      <div>
        {{
          $t("applications.restart_app", {
            name: instanceToRestart ? instanceToRestart.id : "",
          })
        }}
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("applications.restart_instance") }}</template>
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
