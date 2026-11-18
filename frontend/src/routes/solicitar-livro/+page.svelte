<script>
  import Header from "$lib/components/home/Header.svelte";
  import { onMount } from "svelte";
  import * as constants from "../../stores/constants";
  import Captcha from "$lib/components/home/Captcha.svelte";
  import { recaptcha } from "svelte-recaptcha-v2";
  import { handlePhone } from "$lib/utils";

  export let data;
  /**
   * @type {any}
   */
  let token;

  onMount(() => {
    alertError.style.display = "none";
    alertSuccess.style.display = "none";
  });

  const save = async () => {
    if (submitButton) submitButton.disabled = true;

    const body = {
      name: thisName?.value || "",
      email: emailAddress?.value || "",
      phone: phone?.value || "",
      address: address?.value || "",
      token: token || "",
    };

    alertError.style.display = "none";
    alertSuccess.style.display = "none";

    await fetch("/api/solicitacao/livro", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })
      .then(async (resp) => {
        if (!resp.ok) {
          const err = await resp.json();
          messageText = err.error || err.message || "Erro desconhecido!";
          alertError.style.display = "block";
          return;
        }

        alertSuccess.style.display = "block";
        messageText =
          "Solicitação enviada com sucesso! Aguarde, que em breve entraremos em contato.";
      })
      .catch((err) => {
        messageText = `Erro ao executar requisição: ${err}`;
        alertError.style.display = "block";
        return;
      });
    if (submitButton) submitButton.disabled = false;
  };

  let messageText = "";

  const cancelFunction = () => {
    if (submitButton) submitButton.disabled = false;
  };
</script>

<Header />
<section
  on:submit|preventDefault={async () => await recaptcha.execute()}
  id="thisForm"
  class="p-6 mx-auto bg-white shadow-md dark:bg-gray-800"
>
  <h2 class="text-lg font-semibold text-gray-700 capitalize dark:text-white">
    Solicitar livro O Grande Conflito
  </h2>

  <form>
    <div class="grid grid-cols-1 gap-6 mt-4 sm:grid-cols-2">
      <div>
        <label class="text-gray-700 dark:text-gray-200" for="thisName"
          >Nome</label
        >
        <input
          id="thisName"
          type="text"
          class="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
          required
        />
      </div>

      <div>
        <label class="text-gray-700 dark:text-gray-200" for="emailAddress"
          >E-mail</label
        >
        <input
          id="emailAddress"
          type="email"
          class="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
          required
        />
      </div>

      <div>
        <label class="text-gray-700 dark:text-gray-200" for="phone"
          >Telefone</label
        >
        <input
          id="phone"
          type="tel"
          maxlength="16"
          on:keyup={(event) => handlePhone(event)}
          pattern={constants.phonePattern}
          class="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
          required
        />
        <small class="text-gray-700 dark:text-gray-200"
          >Formato: {constants.phoneFormat}</small
        ><br /><br />
      </div>

      <div>
        <label class="text-gray-700 dark:text-gray-200" for="address"
          >Endereço</label
        >
        <textarea
          id="address"
          cols="30"
          rows="3"
          maxlength="500"
          class="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
          required
        ></textarea>
      </div>
    </div>

    <Captcha
      data={{ siteKey: data.siteKey, successFunction: save, cancelFunction }}
      bind:token
    />

    <div class="flex justify-end mt-6">
      <button
        id="submitButton"
        class="px-8 py-2.5 leading-5 text-white transition-colors duration-300 transform bg-gray-700 rounded-md hover:bg-gray-600 focus:outline-none focus:bg-gray-600"
        >Solicitar</button
      >
    </div>

    <!-- ALERT SUCCESS -->
    <div
      id="alertSuccess"
      class="flex w-full max-w-sm overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800"
    >
      <div class="flex items-center justify-center w-12 bg-emerald-500">
        <svg
          class="w-6 h-6 text-white fill-current"
          viewBox="0 0 40 40"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M20 3.33331C10.8 3.33331 3.33337 10.8 3.33337 20C3.33337 29.2 10.8 36.6666 20 36.6666C29.2 36.6666 36.6667 29.2 36.6667 20C36.6667 10.8 29.2 3.33331 20 3.33331ZM16.6667 28.3333L8.33337 20L10.6834 17.65L16.6667 23.6166L29.3167 10.9666L31.6667 13.3333L16.6667 28.3333Z"
          />
        </svg>
      </div>

      <div class="px-4 py-2 -mx-3">
        <div class="mx-3">
          <span class="font-semibold text-emerald-500 dark:text-emerald-400"
            >Sucesso</span
          >
          <p class="text-sm text-gray-600 dark:text-gray-200">
            {messageText}
          </p>
        </div>
      </div>
    </div>

    <!-- ALERT Erro -->
    <div
      id="alertError"
      class="flex w-full max-w-sm overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800"
    >
      <div class="flex items-center justify-center w-12 bg-red-500">
        <svg
          class="w-6 h-6 text-white fill-current"
          viewBox="0 0 40 40"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M20 3.36667C10.8167 3.36667 3.3667 10.8167 3.3667 20C3.3667 29.1833 10.8167 36.6333 20 36.6333C29.1834 36.6333 36.6334 29.1833 36.6334 20C36.6334 10.8167 29.1834 3.36667 20 3.36667ZM19.1334 33.3333V22.9H13.3334L21.6667 6.66667V17.1H27.25L19.1334 33.3333Z"
          />
        </svg>
      </div>

      <div class="px-4 py-2 -mx-3">
        <div class="mx-3">
          <span class="font-semibold text-red-500 dark:text-red-400">Erro</span>
          <p class="text-sm text-gray-600 dark:text-gray-200">
            {messageText}
          </p>
        </div>
      </div>
    </div>
    <!-- FIM DO ALERT -->
  </form>
</section>
