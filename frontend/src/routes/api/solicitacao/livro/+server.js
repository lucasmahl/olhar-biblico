import { config } from "$lib/config.js";
import { validateCaptchaToken } from "$lib/utils.js";

export async function POST({ request, getClientAddress }) {
  const req = await request?.json();

  if (
    String(req.name).trim() == "" ||
    String(req.email).trim() == "" ||
    String(req.phone).trim() == "" ||
    String(req.address).trim() == "" ||
    String(req.token).trim() == ""
  ) {
    return new Response(
      JSON.stringify({ error: "Formulário está incompleto!" }),
      { status: 400 }
    );
  }

  const headers = {
    "Content-Type": "application/json",
  };

  const { API_HOST = "", API_PORT = "" } = config() ?? {};

  try {
    await validateCaptchaToken(req.token);
  } catch (/** @type {*} */ err) {
    return err;
  }

  const body = {
    name: req.name || "",
    email: req.email || "",
    phone: req.phone || "",
    address: req.address || "",
    remote_address_front_end: getClientAddress(),
  };

  return await fetch(`${API_HOST}:${API_PORT}/book-requests/create`, {
    method: "POST",
    headers: headers,
    body: JSON.stringify(body),
  })
    .then((res) => res)
    .catch(
      (_) =>
        new Response(
          JSON.stringify({
            error: "Ocorreu um erro ao enviar sua solicitação.",
          }),
          { status: 500 }
        )
    );
}
