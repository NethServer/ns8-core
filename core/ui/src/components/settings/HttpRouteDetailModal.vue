<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal size="default" :visible="isShown" @modal-hidden="$emit('hide')">
    <template slot="title">{{
      $t("settings_http_routes.http_route_details")
    }}</template>
    <template v-if="route" slot="content">
      <div class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.name") }}</span>
        <span class="value">{{ route.name }}</span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.url") }}</span>
        <span class="value">{{ route.url }}</span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{
          $t("settings_http_routes.skip_cert_verify")
        }}</span>
        <span class="value">
          <cv-tag
            :kind="route.skip_cert_verify ? 'green' : 'high-contrast'"
            size="sm"
            :label="
              route.skip_cert_verify
                ? $t('common.enabled')
                : $t('common.disabled')
            "
            class="no-margin"
          ></cv-tag>
        </span>
      </div>
      <div v-if="route.host" class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.host") }}</span>
        <span class="value">{{ route.host }}</span>
      </div>
      <div v-if="route.path" class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.path") }}</span>
        <span class="value">{{ route.path }}</span>
      </div>
      <div v-if="route.path" class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.strip_prefix") }}</span>
        <span class="value">
          <cv-tag
            :kind="route.strip_prefix ? 'green' : 'high-contrast'"
            size="sm"
            :label="
              route.strip_prefix ? $t('common.enabled') : $t('common.disabled')
            "
            class="no-margin"
          ></cv-tag>
        </span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{
          $t("settings_http_routes.request_lets_encrypt_certificate")
        }}</span>
        <span class="value">
          <cv-tag
            :kind="route.lets_encrypt ? 'green' : 'high-contrast'"
            size="sm"
            :label="
              route.lets_encrypt ? $t('common.enabled') : $t('common.disabled')
            "
            class="no-margin"
          ></cv-tag>
        </span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.http2https") }}</span>
        <span class="value">
          <cv-tag
            :kind="route.http2https ? 'green' : 'high-contrast'"
            size="sm"
            :label="
              route.http2https ? $t('common.enabled') : $t('common.disabled')
            "
            class="no-margin"
          ></cv-tag>
        </span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{ $t("settings_http_routes.node") }}</span>
        <span class="value">{{ route.longNodeLabel }}</span>
      </div>
      <div class="key-value-setting">
        <span class="label">{{
          $t("settings_http_routes.traefik_instance")
        }}</span>
        <span class="value">{{ route.traefikInstance }}</span>
      </div>
    </template>
    <template slot="primary-button">{{ $t("common.close") }}</template>
  </NsModal>
</template>

<script>
import { UtilService } from "@nethserver/ns8-ui-lib";

export default {
  name: "HttpRouteDetailModal",
  mixins: [UtilService],
  props: {
    isShown: Boolean,
    route: Object,
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>
