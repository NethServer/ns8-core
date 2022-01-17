<template>
  <div class="shell-footer">
    <div class="footer-left">
      <div class="item">
        <!-- product name and URL -->
        <cv-link
          v-if="$root.config.PRODUCT_URL"
          :href="$root.config.PRODUCT_URL"
          target="_blank"
        >
          {{ $root.config.PRODUCT_NAME }}
        </cv-link>
        <span v-else>{{ $root.config.PRODUCT_NAME }}</span>
        <!-- company name and URL -->
        <template v-if="$root.config.COMPANY_NAME">
          {{ $t("footer.by_company") }}
          <cv-link
            v-if="$root.config.COMPANY_URL"
            :href="$root.config.COMPANY_URL"
            target="_blank"
          >
            {{ $root.config.COMPANY_NAME }}
          </cv-link>
          <span v-else>
            {{ $root.config.COMPANY_NAME }}
          </span>
        </template>
      </div>
      <div v-if="$root.config.DOCS_URL" class="item">
        <cv-link :href="$root.config.DOCS_URL" target="_blank">
          {{ $t("footer.docs") }}
        </cv-link>
      </div>
      <div v-if="$root.config.HELPDESK_URL" class="item">
        <cv-link :href="$root.config.HELPDESK_URL" target="_blank">
          {{ $t("footer.helpdesk") }}
        </cv-link>
      </div>
    </div>
    <div class="footer-right">
      <div class="item">
        {{ $t("footer.animated_icons_by") }}
        <cv-link href="https://icons8.com/" target="_blank"> icons8 </cv-link>
      </div>
    </div>
  </div>
</template>

<script>
//// delete file?

export default {
  name: "ShellFooter",
  data() {
    return {};
  },
  methods: {},
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.shell-footer {
  position: relative;
  top: 3rem;
  height: 3rem;
  z-index: 8000;
  bottom: 0;
  right: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px solid $inverse-02;
  background-color: $ui-05;
  color: $ui-04;
  font-size: 0.7rem;

  .footer-left {
    display: flex;
  }

  .footer-right {
    display: flex;
    justify-content: flex-end;
  }

  .item {
    padding-left: $spacing-05;
    padding-right: $spacing-05;
  }
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.shell-footer .item {
  .bx--link {
    font-size: 0.7rem;
    color: $ui-04;
    text-decoration: underline;
  }

  .bx--link:active {
    color: $ui-04 !important;
  }
}
</style>
