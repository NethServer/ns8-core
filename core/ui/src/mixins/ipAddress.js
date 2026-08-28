const IPV4_OCTET = "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])";
const IPV4 = `${IPV4_OCTET}(\\.${IPV4_OCTET}){3}`;
const IPV6 = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}${IPV4}|([0-9a-fA-F]{1,4}:){1,4}:${IPV4})`;
const IPV4_PREFIX = "(\\/([0-9]|[1-2][0-9]|3[0-2]))?";
const IPV6_PREFIX = "(\\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?";

const IPV4_PATTERN = new RegExp(`^${IPV4}$`);
const IPV6_PATTERN = new RegExp(`^${IPV6}$`);
const IPV4_CIDR_PATTERN = new RegExp(`^${IPV4}${IPV4_PREFIX}$`);
const IPV6_CIDR_PATTERN = new RegExp(`^${IPV6}${IPV6_PREFIX}$`);

export default {
  name: "IpAddressService",
  methods: {
    isIpv4Address(value) {
      return IPV4_PATTERN.test(value);
    },
    isIpv6Address(value) {
      return IPV6_PATTERN.test(value);
    },
    // no prefix: set-trusted-proxies validates with ipaddress.ip_address()
    isIpAddress(value) {
      return this.isIpv4Address(value) || this.isIpv6Address(value);
    },
    isIpAddressOrCidr(value) {
      return IPV4_CIDR_PATTERN.test(value) || IPV6_CIDR_PATTERN.test(value);
    },
    // picks the error message: a rejected value with a colon was meant as IPv6
    looksLikeIpv6(value) {
      return value.includes(":");
    },
  },
};
