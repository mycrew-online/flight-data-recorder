<script lang="ts">
  import './../app.css';
  import SimStatus from '$lib/components/SimStatus.svelte';
  import { simStatus } from '$lib/stores/simStatus';
  import SimulatorStatePanel from '$lib/components/SimulatorStatePanel.svelte';
  import { simulatorState } from '$lib/stores/simulatorState';
  const { children } = $props();
  let sidebarOpen = $state(false);

</script>

<div class="h-full bg-gray-100 min-h-screen">
  <!-- Off-canvas menu for mobile -->
  {#if $simStatus}
  <div class={sidebarOpen ? "relative z-50 lg:hidden" : "hidden"} role="dialog" aria-modal="true">
    <div class="fixed inset-0 bg-gray-900/80" aria-hidden="true" onclick={() => sidebarOpen = false}></div>
    <div class="fixed inset-0 flex">
      <div class="relative mr-16 flex w-full max-w-xs flex-1">
        <div class="absolute top-0 left-full flex w-16 justify-center pt-5">
          <button type="button" class="-m-2.5 p-2.5" onclick={() => sidebarOpen = false}>
            <span class="sr-only">Close sidebar</span>
            <svg class="size-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-4 ring-1 ring-white/10">
          <div class="flex h-16 shrink-0 items-center">
            <img class="h-8 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500" alt="Your Company" />
          </div>
          <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-7">
              <li>
                <ul role="list" class="-mx-2 space-y-1">
                  <li>
                    <a href="/" class="group flex gap-x-3 rounded-md bg-gray-800 p-2 text-sm font-semibold text-white">
                      Dashboard
                    </a>
                  </li>
                  <li>
                    <a href="/logs" class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold text-gray-400 hover:bg-gray-800 hover:text-white">
                      Logs
                    </a>
                  </li>
                  <li>
                    <a href="/settings" class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold text-gray-400 hover:bg-gray-800 hover:text-white">
                      Settings
                    </a>
                  </li>
                </ul>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>
  </div>
  {/if}

  <!-- Static sidebar for desktop, closeable on small screens -->
  {#if $simStatus}
  <div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col" class:!hidden={sidebarOpen === false && window.innerWidth < 1024}>
    <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-4">
      <div class="flex h-16 shrink-0 items-center justify-between">
        <img class="h-8 w-auto" src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500" alt="Your Company" />
        <button type="button" class="lg:hidden ml-2 p-2 text-gray-400 hover:text-white" onclick={() => sidebarOpen = false} aria-label="Close sidebar">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <nav class="flex flex-1 flex-col">
        <ul role="list" class="flex flex-1 flex-col gap-y-7">
          <li>
            <ul role="list" class="-mx-2 space-y-1">
              <li>
                <a href="/" class="group flex gap-x-3 rounded-md bg-gray-800 p-2 text-sm font-semibold text-white">
                  Dashboard
                </a>
              </li>
              <li>
                <a href="/logs" class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold text-gray-400 hover:bg-gray-800 hover:text-white">
                  Logs
                </a>
              </li>
              <li>
                <a href="/settings" class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold text-gray-400 hover:bg-gray-800 hover:text-white">
                  Settings
                </a>
              </li>
            </ul>
          </li>
        </ul>
      </nav>
    </div>
  </div>
  {/if}

  <div class={$simStatus ? "lg:pl-72" : ""}>
    <div class="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-xs sm:gap-x-6 sm:px-6 lg:px-8">
      {#if $simStatus}
      <button type="button" class="-m-2.5 p-2.5 text-gray-700 lg:hidden" onclick={() => sidebarOpen = true}>
        <span class="sr-only">Open sidebar</span>
        <svg class="size-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
        </svg>
      </button>
      {/if}
      <div class="flex flex-1 items-center justify-end gap-4">
        <SimStatus />
      </div>
    </div>
    <main class="h-min-full">
      <div>
        {@render children()}
      </div>
    </main>
  </div>
</div>