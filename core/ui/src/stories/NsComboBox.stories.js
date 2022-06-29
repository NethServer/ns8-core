// import NsComboBox from "../components/NsComboBox"; //// delete
import { NsComboBox } from "@nethserver/ns8-ui-lib"; //// uncomment

export default {
  title: "Components/NsComboBox",
  component: NsComboBox,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsComboBox },
  template: '<NsComboBox v-bind="$props" />',
});

const vegetablesOptions = ["Salad", "Tomato", "Cucumber", "Carrot"].map(
  (item) => {
    const nameVal = item.replace(/\W/, "_").toLowerCase();
    return {
      name: nameVal,
      label: item,
      value: nameVal,
      disabled: false,
      type: "vegetable",
    };
  }
);

const fruitOptions = [
  "Apple",
  "Banana",
  "Cherry",
  "Date",
  "Elderberry",
  "Fig",
  "Grape",
  "Kiwi Fruit",
  "Lemon",
  "Lime",
  "Mango",
  "Orange",
  "Passion Fruit",
  "Raisin",
  "Satsuma",
  "Tangerine",
  "Ugli Fruit",
  "Watermelon",
].map((item) => {
  const nameVal = item.replace(/\W/, "_").toLowerCase();
  return {
    name: nameVal,
    label: item,
    value: nameVal,
    disabled: false,
    type: "fruit",
  };
});

export const Default = Template.bind({});
Default.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: false,
  options: fruitOptions,
};

export const ShowItemType = Template.bind({});
ShowItemType.args = {
  autoFilter: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: false,
  showItemType: true,
  options: vegetablesOptions.concat(fruitOptions),
};

export const AcceptUserInput = Template.bind({});
AcceptUserInput.args = {
  autoFilter: true,
  autoHighlight: false,
  disabled: false,
  invalidMessage: undefined,
  helperText: undefined,
  title: "Title",
  label: "Choose an option or input any value",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear filter",
  userInputLabel: "user input",
  maxDisplayOptions: 100,
  acceptUserInput: true,
  showItemType: true,
  options: fruitOptions,
};
