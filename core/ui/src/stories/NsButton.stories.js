import { NsButton } from "@nethserver/ns8-ui-lib";
import Save from "@carbon/icons-vue/es/save/20";

export default {
  title: "Components/NsButton",
  component: NsButton,
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
    size: {
      options: ["default", "field", "small", "sm", "lg", "xl"],
      control: { type: "select" },
    },
    buttonClick: { action: "clickAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsButton },
  template:
    '<NsButton v-bind="$props" @click="buttonClick(\'testArg\')">{{ slotContent }}</NsButton>',
});

export const Default = Template.bind({});
Default.args = {
  slotContent: `Button text`,
  kind: "secondary",
  loading: false,
  value: "selected",
};

export const Icon = Template.bind({});
Icon.args = {
  slotContent: `Button text`,
  kind: "secondary",
  loading: false,
  value: "selected",
  icon: Save,
};
