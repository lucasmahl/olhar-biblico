import { config } from "./config";

/**
 * @param {string} token
 */
export async function validateCaptchaToken(token) {
  const { RECAPTCHA_SECRET_KEY = "" } = config() ?? {};
  const response = await fetch(
    "https://www.google.com/recaptcha/api/siteverify",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `secret=${RECAPTCHA_SECRET_KEY}&response=${token}`,
    }
  )
    .then((res) => res)
    .catch((_) => {
      throw new Response(
        JSON.stringify({
          error: "Erro ao enviar solicitação para validar token!",
        }),
        { status: 404 }
      );
    });

  const responseData = await response.json();
  if (responseData?.success !== true) {
    throw new Response(
      JSON.stringify({ error: "Token do captcha é inválido!" }),
      { status: 400 }
    );
  }
}

/** @param {KeyboardEvent & { currentTarget: EventTarget & HTMLInputElement; }} event*/
export const handlePhone = (event) => {
  let input = event?.target;
  if (!input?.value) return "";
  input.value = input?.value?.replace(/\D/g, "");
  if (input?.value?.length === 10) {
    input.value = input.value.replace(/(\d{2})(\d{4})(\d{4})/, "($1) $2-$3");
  } else if (input?.value?.length === 11) {
    input.value = input?.value.replace(/(\d{2})(\d{5})(\d{4})/, "($1) $2-$3");
  } else if (input?.value?.length === 12) {
    input.value = input?.value.replace(/(\d{3})(\d{5})(\d{4})/, "($1) $2-$3");
  }

  return input.value;
};
