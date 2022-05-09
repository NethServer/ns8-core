import { NsWizard } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsWizard",
  component: NsWizard,
  argTypes: {
    size: {
      options: ["default", "small", "large"],
      control: { type: "radio" },
    },
    buttonClick: { action: "clickAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsWizard },
  template:
    '<NsWizard v-bind="$props"\
      @cancel="buttonClick(\'cancel\')"\
      @previousStep="buttonClick(\'previous\')"\
      @nextStep="buttonClick(\'next\')"\
    >\
      <template slot="label">\
        {{ slotLabel }}</template>\
      </template>\
      <template slot="title">\
        {{ slotTitle }}\
      </template>\
      <template slot="content">\
        {{ slotContent }}\
      </template>\
    </NsWizard>',
});

export const Default = Template.bind({});
Default.args = {
  cancelLabel: "Cancel",
  previousLabel: "Previous",
  nextLabel: "Next",
  isPreviousDisabled: false,
  isNextDisabled: false,
  isNextLoading: false,
  isCancelDisabled: false,
  closeAriaLabel: "Close modal",
  autoHideOff: false,
  visible: true,
  size: "default",
  slotLabel: "",
  slotTitle: "Wizard title",
  slotContent: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
  Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
  Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
  Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
};
