import { format, formatDistance, subDays } from "date-fns";

export default {
  name: "DateTimeService",
  data() {
    return {
      formatDate: format,
      formatDistance,
      subDays,
    };
  },
  methods: {},
};
