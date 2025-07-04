<script lang="ts">
  import { goto } from '$app/navigation';

  import GeneralTab from './GeneralTab.svelte';
  import ThirdPartyTab from './ThirdPartyTab.svelte';

  let tabs = [
    { name: 'General', component: GeneralTab },
    { name: '3rd party', component: ThirdPartyTab },
  ];
  let selectedTab: number = 0; // Default to Interview
  function selectTab(idx: number) {
    selectedTab = idx;
  }
  function selectTabFromDropdown(e: Event) {
    const idx = Array.from((e.target as HTMLSelectElement).options).findIndex(option => option.selected);
    selectedTab = idx;
  }

  function onSave() {
    // TODO: Implement save logic
    alert('Settings saved!');
  }
  function onClose() {
    goto('/');
  }

</script>

<div class="bg-gray-900 min-h-screen text-white">
  <div class="p-6 sm:p-8 md:p-10 lg:p-12">
    <div>
      <nav class="sm:hidden" aria-label="Back">
        <a href="/" class="flex items-center text-sm font-medium text-gray-400 hover:text-gray-200">
          <svg class="mr-1 -ml-1 size-5 shrink-0 text-gray-500" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" data-slot="icon">
            <path fill-rule="evenodd" d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z" clip-rule="evenodd" />
          </svg>
          Back
        </a>
      </nav>
      <nav class="hidden sm:flex" aria-label="Breadcrumb">
        <ol role="list" class="flex items-center space-x-4">
          <li>
            <div class="flex">
              <a href="/" class="text-sm font-medium text-gray-400 hover:text-gray-200">Home</a>
            </div>
          </li>
          <li>
            <div class="flex items-center">
              <svg class="size-5 shrink-0 text-gray-500" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" data-slot="icon">
                <path fill-rule="evenodd" d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
              </svg>
              <a href="/settings" aria-current="page" class="ml-4 text-sm font-medium text-gray-400 hover:text-gray-200">Settings</a>
            </div>
          </li>
        </ol>
      </nav>
    </div>
    <div class="mt-6 md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl/7 font-bold text-white sm:truncate sm:text-3xl sm:tracking-tight">Settings</h2>
      </div>
      <div class="mt-4 flex shrink-0 md:mt-0 md:ml-4">
        <a class="inline-flex items-center rounded-full bg-white/10 p-2 mr-2 text-white hover:bg-white/20 focus:outline-none" href="/logs" aria-label="Close settings" title="View Logs">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
            <title>View Logs</title>
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 12.75c1.148 0 2.278.08 3.383.237 1.037.146 1.866.966 1.866 2.013 0 3.728-2.35 6.75-5.25 6.75S6.75 18.728 6.75 15c0-1.046.83-1.867 1.866-2.013A24.204 24.204 0 0 1 12 12.75Zm0 0c2.883 0 5.647.508 8.207 1.44a23.91 23.91 0 0 1-1.152 6.06M12 12.75c-2.883 0-5.647.508-8.208 1.44.125 2.104.52 4.136 1.153 6.06M12 12.75a2.25 2.25 0 0 0 2.248-2.354M12 12.75a2.25 2.25 0 0 1-2.248-2.354M12 8.25c.995 0 1.971-.08 2.922-.236.403-.066.74-.358.795-.762a3.778 3.778 0 0 0-.399-2.25M12 8.25c-.995 0-1.97-.08-2.922-.236-.402-.066-.74-.358-.795-.762a3.734 3.734 0 0 1 .4-2.253M12 8.25a2.25 2.25 0 0 0-2.248 2.146M12 8.25a2.25 2.25 0 0 1 2.248 2.146M8.683 5a6.032 6.032 0 0 1-1.155-1.002c.07-.63.27-1.222.574-1.747m.581 2.749A3.75 3.75 0 0 1 15.318 5m0 0c.427-.283.815-.62 1.155-.999a4.471 4.471 0 0 0-.575-1.752M4.921 6a24.048 24.048 0 0 0-.392 3.314c1.668.546 3.416.914 5.223 1.082M19.08 6c.205 1.08.337 2.187.392 3.314a23.882 23.882 0 0 1-5.223 1.082" />
          </svg>
        </a>
        <button type="button" class="inline-flex items-center rounded-full bg-white/10 p-2 text-white hover:bg-white/20 focus:outline-none" onclick={onClose} aria-label="Close settings" title="Close Settings">
          <svg class="size-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <title>Close Settings</title>
            <path fill-rule="evenodd" d="M10 8.586l4.95-4.95a1 1 0 1 1 1.414 1.414L11.414 10l4.95 4.95a1 1 0 0 1-1.414 1.414L10 11.414l-4.95 4.95a1 1 0 0 1-1.414-1.414L8.586 10l-4.95-4.95A1 1 0 1 1 5.05 3.636L10 8.586z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Settings form/content goes here -->



    <div class="relative border-b border-gray-200 pb-5 sm:pb-0 bg-gray-100 min-h-full mt-6">
      <!-- svelte-ignore a11y_missing_content -->
      <div class="mt-8 px-4 pb-2 pt-7">
        <div class="grid grid-cols-1 sm:hidden">
          <select aria-label="Select a tab" class="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-2 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600" onchange={selectTabFromDropdown}>
            {#each tabs as tab, idx}
              <option value={tab.name} selected={selectedTab === idx}>{tab.name}</option>
            {/each}
          </select>
          <svg class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end fill-gray-500" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" data-slot="icon">
            <path fill-rule="evenodd" d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
          </svg>
        </div>
        <!-- Tabs at small breakpoint and up -->
        <div class="hidden sm:block">
          <nav class="-mb-px flex space-x-8">
            {#each tabs as tab, idx}
              <button type="button"
                class="border-b-2 px-1 pb-4 text-sm font-medium whitespace-nowrap focus:outline-none {selectedTab === idx ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'}"
                aria-current={selectedTab === idx ? 'page' : undefined}
                onclick={() => selectTab(idx)}>{tab.name}</button>
            {/each}
          </nav>
        </div>
      </div>
      <div class="bg-white rounded-b-md p-6 mt-2 min-h-[120px] text-gray-900">
        {#if tabs && typeof selectedTab === 'number' && tabs[selectedTab]}
          {#key selectedTab}
            <svelte:component this={tabs[selectedTab].component} />
          {/key}
        {/if}
      </div>
    </div>
  </div>
</div>