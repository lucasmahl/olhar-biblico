<script>
  import BemVindo from "$lib/components/home/BemVindo.svelte";
  import Contato from "$lib/components/home/Contato.svelte";
  import Header from "$lib/components/home/Header.svelte";
  import Igrejas from "$lib/components/home/Igrejas.svelte";

  export let data;
</script>

<Header />
<Contato {data} />
<Igrejas />
<BemVindo />
