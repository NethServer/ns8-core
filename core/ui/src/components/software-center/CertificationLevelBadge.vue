<!--
  Copyright (C) 2024 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-interactive-tooltip alignment="center" direction="top">
    <template slot="trigger">
      <NsTag
        :label="
          $t('software_center.level_level', {
            level: `${level}/5`,
          })
        "
        :kind="getLevelColor()"
        :icon="getLevelIcon()"
        class="level-tooltip"
      />
    </template>
    <template slot="content">
      <h6 class="mg-bottom-sm">
        {{ $t("software_center.level_level", { level }) }}
      </h6>
      <p class="mg-bottom-sm">
        {{ $t(`software_center.level_${level}_description`) }}
      </p>
      <cv-link
        href="https://docs.nethserver.org/projects/ns8/en/latest/software_center.html#certification-levels"
        target="_blank"
        rel="noreferrer"
        class="inline"
      >
        {{ $t("common.more_info") }}
      </cv-link>
    </template>
  </cv-interactive-tooltip>
</template>

<script>
import Security16 from "@carbon/icons-vue/es/security/16";
import WarningAlt16 from "@carbon/icons-vue/es/warning--alt/16";
import Error16 from "@carbon/icons-vue/es/error/16";
import Rule16 from "@carbon/icons-vue/es/rule/16";

export default {
  name: "CertificationLevelBadge",
  props: {
    level: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      Security16,
      WarningAlt16,
      Rule16,
      Error16,
    };
  },
  methods: {
    getLevelColor() {
      if (this.level > 0 && this.level < 5) {
        return "blue";
      } else if (this.level == 5) {
        return "green";
      } else {
        return "gray";
      }
    },
    getLevelIcon() {
      switch (this.level) {
        case 1:
          return this.WarningAlt16;
        case 2:
        case 3:
        case 4:
          return this.Rule16;
        case 5:
          return this.Security16;
        default:
          return this.Error16;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.level-tooltip {
  cursor: pointer;
}
</style>
