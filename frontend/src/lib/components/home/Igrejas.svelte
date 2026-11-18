<script>
  import { onMount } from "svelte";

  const igrejas = [
    {
      endereco: "R. Santo Amaro, 196 - Centro, Esteio",
      lat: -29.8508830676243,
      long: -51.17488859804776,
    },
    {
      endereco: "R. Pernambuco, 92-178 - Parque Amador, Esteio",
      lat: -29.84620406333103,
      long: -51.15123536055298,
    },
    {
      endereco: "Av. Rio Branco, 336 - Bairro do Parque, Esteio",
      lat: -29.846011583753388,
      long: -51.18457126688511,
    },
  ];

  /**
   * @type {any[]}
   */
  const markers = [];
  /**
   * @type {{ setView: (arg0: number[], arg1: number) => void; }}
   */
  let map;
  onMount(async () => {
    map = await L.map("myMap");

    await L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution:
        '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      maxZoom: 20,
      minZoom: 5,
      accessToken:
        "",
      id: "mapbox/satellite-streets-v11",
    }).addTo(map);

    for await (const igreja of igrejas) {
      const marker = L.marker([igreja.lat, igreja.long])
        .addTo(map)
        .bindPopup(igreja.endereco)
        .openPopup();

      markers.push(marker);

      map.setView([igreja.lat, igreja.long], 15);
    }
  });
</script>

<svelte:head>
  <link
    rel="stylesheet"
    href="/leaflet/leaflet.css"
    integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
    crossorigin=""
  />
  <script
    src="/leaflet/leaflet.js"
    integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
    crossorigin=""
  ></script>
</svelte:head>

<section id="igrejas" class="bg-white">
  <div class="container px-6 py-10 mx-auto">
    <h1
      class="text-2xl font-semibold text-gray-800 capitalize lg:text-3xl dark:text-white"
    >
      ENCONTRE UMA IGREJA
    </h1>

    <div class="mt-8 lg:-mx-6 lg:flex lg:items-center">
      <div
        id="myMap"
        class="object-cover w-full lg:mx-6 lg:w-1/2 rounded-xl h-72 lg:h-96"
      ></div>

      <div class="mt-6 lg:w-1/2 lg:mt-0 lg:mx-6">
        <p
          class="block mt-1 text-4xl font-semibold text-gray-800 hover:underline dark:text-white"
        >
          Nossos endereços
        </p>

        {#if igrejas?.length}
          {#each igrejas as igreja, index}
            <p class="mt-3 text-2xl text-gray-500 dark:text-gray-300">
              <button
                type="button"
                on:click={() => {
                  markers[index].openPopup();
                  map.setView([igreja.lat, igreja.long], 15);
                }}>📍 - {igreja?.endereco}</button
              >
            </p>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</section>

<style>
  #igrejas {
    background-color: #09326c;
  }
</style>
