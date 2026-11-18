<script>
  import { Recaptcha } from "svelte-recaptcha-v2";

  export let token;

  /**
   * @type {{ successFunction: () => any; cancelFunction: () => any; siteKey: any; }}
   */
  export let data;

  const onCaptchaReady = () => {
    console.log("recaptcha init has completed.");
    captchaDiv.style.display = "none";
  };

  /** @param {{ detail: { token: any; }; }} event */
  const onCaptchaSuccess = async (event) => {
    let userTracker = {
      resolve: async function () {
        token = event.detail.token;
        await data.successFunction();
      },
    };

    await userTracker.resolve();
  };

  const onCaptchaError = () => {
    console.log("recaptcha init has failed.");
  };

  const onCaptchaExpire = () => {
    console.log("recaptcha api has expired");
    /*
      │Normally, you wouldn't need to do anything.
      │Recaptcha should reinit itself automatically.
      */
  };

  const onCaptchaOpen = () => {
    console.log("google decided to challange the user");
    captchaDiv.style.display = "block";
  };

  const onCaptchaClose = async () => {
    console.log("User closed the challenge");
    await data.cancelFunction();
    captchaDiv.style.display = "none";
  };
</script>

<div id="captchaDiv" hidden>
  <Recaptcha
    sitekey={data.siteKey}
    badge={"top"}
    size={"invisible"}
    on:success={onCaptchaSuccess}
    on:error={onCaptchaError}
    on:expired={onCaptchaExpire}
    on:close={onCaptchaClose}
    on:ready={onCaptchaReady}
    on:open={onCaptchaOpen}
  />
</div>
