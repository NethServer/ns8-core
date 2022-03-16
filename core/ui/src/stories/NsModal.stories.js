import { NsModal } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsModal",
  component: NsModal,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsModal },
  template:
    '<NsModal v-bind="$props">\
        <template slot="title">Modal title</template>\
        <template slot="content">\
          <div class="mg-bottom-md">See <a href="https://vue.carbondesignsystem.com/?path=/story/components-cvmodal--default" target="_blank">CvModal Storybook</a> for full documentation.</div>\
          <div>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\
          Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</div>\
        </template>\
        <template slot="secondary-button">Close</template>\
        <template slot="primary-button">Do something</template>\
    </NsModal>',
});

export const Default = Template.bind({});
Default.args = {
  hideOnClickOutside: false,
  visible: true,
};
