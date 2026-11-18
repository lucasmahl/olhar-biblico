<script>
  import { page } from "$app/stores";
  import Microfone from "./Microfone.svelte";

  /**
   * @type {{ href: string; title: string; }[]}
   */
  let navItems = [];
  if ($page.url.pathname === "/") {
    navItems = [
      { href: "#contact", title: "Contato" },
      { href: "#igrejas", title: "Encontre uma igreja" },
      { href: "#bemVindos", title: "Quem somos" },
    ];
  }
</script>

<nav class="bg-white shadow nav">
  <div
    class="header container flex items-center justify-center mx-auto text-gray-600 dark:text-gray-300"
  >
    <div>
      <button on:click={() => window.location.replace("/")}>
        <Microfone />
      </button>
    </div>

    {#each navItems || [] as { href, title }}
      <a
        {href}
        class="border-b-2 border-transparent hover:text-gray-800 transition-colors duration-300 transform dark:hover:text-gray-200 hover:border-blue-500 mx-1.5 sm:mx-6"
        >{title}</a
      >
    {/each}
  </div>
</nav>

<style>
  .header {
    height: 8rem;
  }

  .nav {
    background-color: #09326c;
  }
</style>
