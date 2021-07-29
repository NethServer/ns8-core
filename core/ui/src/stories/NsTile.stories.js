import { NsTile } from "@nethserver/ns8-ui-lib";
import Notification20 from "@carbon/icons-vue/es/notification/20";

export default {
  title: "Components/NsTile",
  component: NsTile,
  argTypes: {
    kind: {
      options: ["standard", "clickable", "selectable"],
      control: { type: "radio" },
    },
    tileClick: { action: "clickAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsTile },
  template:
    '<NsTile v-bind="$props" @click="tileClick(\'testArg\')">{{ slotContent }}</NsTile>',
});

export const Default = Template.bind({});
Default.args = {
  kind: "standard",
  slotContent: `Tile title`,
  // required for selectable kind
  value: "selectedValue",
  light: false,
  selected: false,
};

export const Icon = Template.bind({});
Icon.args = {
  slotContent: `Tile title`,
  kind: "standard",
  // required for selectable kind
  value: "selectedValue",
  light: false,
  selected: false,
  icon: Notification20,
};
