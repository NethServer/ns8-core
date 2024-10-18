import { NsTimePicker } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsTimePicker",
  component: NsTimePicker,
};

const Template = (args) => ({
  props: ['value', 'label', 'hideClearButton', 'invalidMessage', 'light'], // Props specifiche per Template
  components: { NsTimePicker },
  template: '<NsTimePicker v-bind="$props" />',
});

const Template2 = (args) => ({
  props: ['value', 'label', 'hideClearButton', 'invalidMessage', 'light', 'dropDirection'], // Props specifiche per Template2
  components: { NsTimePicker },
  template: '<NsTimePicker v-bind="$props" style="margin-top: 400px"/>',
});

// Esportazione per il componente Default senza dropDirection
export const Default = Template.bind({});
Default.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
};

// Esportazione per Drop_Up con controllo radio per dropDirection
export const Drop_Up = Template2.bind({});
Drop_Up.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
  dropDirection: 'up',
};
Drop_Up.argTypes = {
  dropDirection: {
    control: {
      type: 'radio',
      options: ['up', 'down', 'auto'],
    },
  },
};
