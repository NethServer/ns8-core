import { NsDropdownAction } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsDropdownAction",
  component: NsDropdownAction,
  argTypes: {
    kind: {
      options: [
        "default",
        "primary",
        "secondary",
        "tertiary",
        "ghost",
        "danger",
        "danger--ghost",
        "danger--tertiary",
      ],
      control: { type: "select" },
    },
    buttonClick: { action: "clickAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsDropdownAction },
  template:
    '<NsDropdownAction v-bind="$props">\
      <template v-slot:trigger>Select</template>\
      <cv-overflow-menu-item @click="buttonClick(\'all\')">All</cv-overflow-menu-item>\
      <cv-overflow-menu-item @click="buttonClick(\'none\')">None</cv-overflow-menu-item>\
      <cv-overflow-menu-item @click="buttonClick(\'mostBeautiful\')">Most beautiful</cv-overflow-menu-item>\
    </NsDropdownAction>',
});

export const Default = Template.bind({});
Default.args = {
  kind: "secondary",
};
