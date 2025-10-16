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
      control: { type: "inline-radio" },
    },
    size: {
      options: ["default", "field", "small", "sm", "lg", "xl"],
      control: { type: "inline-radio" },
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
  kind: "secondary",
  size: "default",
  loading: false,
  value: "selected",
  slotContent: `Button text`,
};

export const Icon = Template.bind({});
Icon.args = {
  kind: "secondary",
  size: "default",
  loading: false,
  value: "selected",
  icon: Save,
  slotContent: `Button text`,
};
