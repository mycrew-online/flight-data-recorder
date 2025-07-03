<script lang="ts">
  import './../../app.css';
  import SimStatus from '$lib/components/SimStatus.svelte';
  import { simStatus } from '$lib/stores/simStatus';
  import { page } from '$app/stores';
  const { children } = $props();
  let sidebarOpen = $state(false);

</script>

<div class="h-full bg-gray-100 min-h-screen">
    <main class="h-min-screen">
      <div>
        {@render children()}
      </div>
    </main>
</div>