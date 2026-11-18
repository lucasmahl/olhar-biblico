<script>
  import * as constants from "../../../stores/constants";
  import Alerta from "./Alerta.svelte";
  import Captcha from "./Captcha.svelte";
  import { recaptcha } from "svelte-recaptcha-v2";
  import { handlePhone } from "$lib/utils.js";
  /**
   * @type {{ siteKey: any; }}
   */
  export let data;

  /**
   * @type {string}
   */
  let token;

  let showAlert = false;
  let message = "";
  let success = true;

  const successFunction = async () => {
    if (submitButton) submitButton.disabled = true;
    showAlert = false;

    const body = {
      name: nameToVisit?.value || "",
      email: emailToVisit?.value || "",
      phone: phoneToVisit?.value || "",
      address: addressToVisit?.value || "",
      token: token || "",
    };

    await fetch("/api/solicitacao/visita", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })
      .then(async (resp) => {
        if (!resp.ok) {
          const err = await resp.json();
          message = err.error || err.message || "Erro desconhecido!";
          showAlert = true;
          success = false;
          return;
        }

        showAlert = true;
        success = true;
        message =
          "Solicitação enviada com sucesso! Aguarde, que em breve entraremos em contato.";
      })
      .catch((err) => {
        message = `Erro ao executar requisição: ${err}`;
        showAlert = true;
        success = false;
        return;
      });

    showAlert = true;
    if (submitButton) submitButton.disabled = false;
  };

  const cancelFunction = () => {
    if (submitButton) submitButton.disabled = false;
  };
</script>

<section id="contact" class="bg-white dark:bg-gray-900">
  <div class="container px-6 py-12 mx-auto">
    <div class="grid grid-cols-1 gap-12 mt-10 lg:grid-cols-2">
      <div class="grid grid-cols-1 gap-12 md:grid-cols-2">
        <div>
          <p class="font-semibold text-2xl text-blue-500 dark:text-yellow-400">
            RECEBA NOSSA VISITA
          </p>

          <h1
            class="mt-2 text-2xl font-semibold text-gray-800 md:text-3xl dark:text-white"
          >
            GOSTARIA DE UMA VISITA OU ESTUDO BÍBLICO?
          </h1>

          <p class="mt-3 text-gray-500 dark:text-gray-400">
            Preencha os dados ao lado para solicitar um estudo bíblico ou apenas
            uma visita.
          </p>
        </div>
      </div>

      <div class="p-4 py-6 rounded-lg bg-gray-50 dark:bg-gray-800 md:p-8">
        <form on:submit|preventDefault={async () => await recaptcha.execute()}>
          <div class="-mx-2 md:items-center md:flex">
            <div class="flex-1 px-2">
              <label
                class="block mb-2 text-sm text-gray-600 dark:text-gray-200"
                for="nameToVisit">Nome</label
              >
              <input
                id="nameToVisit"
                type="text"
                placeholder="John"
                class="block w-full px-5 py-2.5 mt-2 text-gray-700 placeholder-gray-400 bg-white border border-gray-200 rounded-lg dark:placeholder-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:border-gray-700 focus:border-blue-400 dark:focus:border-blue-400 focus:ring-blue-400 focus:outline-none focus:ring focus:ring-opacity-40"
                required
              />
            </div>
          </div>

          <div class="mt-4">
            <label
              class="block mb-2 text-sm text-gray-600 dark:text-gray-200"
              for="emailToVisit">E-mail</label
            >
            <input
              id="emailToVisit"
              type="email"
              placeholder="johndoe@example.com"
              class="block w-full px-5 py-2.5 mt-2 text-gray-700 placeholder-gray-400 bg-white border border-gray-200 rounded-lg dark:placeholder-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:border-gray-700 focus:border-blue-400 dark:focus:border-blue-400 focus:ring-blue-400 focus:outline-none focus:ring focus:ring-opacity-40"
              required
            />
          </div>

          <div class="mt-4">
            <label
              class="block mb-2 text-sm text-gray-600 dark:text-gray-200"
              for="phoneToVisit">Telefone para contato</label
            >
            <input
              id="phoneToVisit"
              type="tel"
              maxlength="16"
              on:keyup={(event) => handlePhone(event)}
              placeholder={constants.phoneFormat}
              pattern={constants.phonePattern}
              class="block w-full px-5 py-2.5 mt-2 text-gray-700 placeholder-gray-400 bg-white border border-gray-200 rounded-lg dark:placeholder-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:border-gray-700 focus:border-blue-400 dark:focus:border-blue-400 focus:ring-blue-400 focus:outline-none focus:ring focus:ring-opacity-40"
              required
            />
            <small class="text-gray-700 dark:text-gray-200"
              >Formato: {constants.phoneFormat}</small
            ><br /><br />
          </div>

          <div class="w-full mt-4">
            <label
              class="block mb-2 text-sm text-gray-600 dark:text-gray-200"
              for="addressToVisit">Endereço</label
            >
            <textarea
              id="addressToVisit"
              class="block w-full h-32 px-5 py-2.5 mt-2 text-gray-700 placeholder-gray-400 bg-white border border-gray-200 rounded-lg dark:placeholder-gray-600 dark:bg-gray-900 dark:text-gray-300 dark:border-gray-700 focus:border-blue-400 dark:focus:border-blue-400 focus:ring-blue-400 focus:outline-none focus:ring focus:ring-opacity-40"
              placeholder="Endereço"
              required
            ></textarea>
          </div>

          {#if showAlert}
            <Alerta {success} {message} />
          {/if}

          <Captcha
            data={{ siteKey: data.siteKey, successFunction, cancelFunction }}
            bind:token
          />

          <button
            id="submitButton"
            class="w-full px-6 py-3 mt-4 text-sm font-medium tracking-wide text-white capitalize transition-colors duration-300 transform bg-blue-500 rounded-lg hover:bg-blue-400 focus:outline-none focus:ring focus:ring-blue-300 focus:ring-opacity-50"
          >
            Solicitar
          </button>
        </form>
      </div>
    </div>
  </div>
</section>
