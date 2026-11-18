import { config } from "$lib/config";

export function load() {
  const { RECAPTCHA_SITE_KEY = "" } = config() ?? {};

  return {
    siteKey: RECAPTCHA_SITE_KEY,
  };
}
