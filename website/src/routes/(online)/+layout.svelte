<script lang="ts">
  import './../../app.css';
  import { simStatus } from '$lib/stores/simStatus';
  import { page } from '$app/stores';
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
                    <a href="/" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                      <!-- Heroicons Home -->
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
                      </svg>
                      <span class="align-middle">Dashboard</span>
                    </a>
                  </li>
                  <li>
                    <a href="/flights" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/flights' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                      <!-- Heroicons Paper Airplane -->
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 21l16.5-9-16.5-9v7.5l13.5 1.5-13.5 1.5V21z" />
                      </svg>
                      <span class="align-middle">My Flights</span>
                    </a>
                  </li>
                  <li>
                    <a href="/logs" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/logs' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                      <!-- Heroicons Document Text -->
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z" />
                      </svg>
                      <span class="align-middle">Logs</span>
                    </a>
                  </li>
                </ul>
              </li>
            </ul>
            <div class="mt-auto">
              <a href="/settings" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/settings' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                <!-- Heroicons Cog-6-Tooth -->
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93.398.164.855.142 1.205-.108l.737-.527a1.125 1.125 0 0 1 1.45.12l.773.774c.39.389.44 1.002.12 1.45l-.527.737c-.25.35-.272.806-.107 1.204.165.397.505.71.93.78l.893.15c.543.09.94.559.94 1.109v1.094c0 .55-.397 1.02-.94 1.11l-.894.149c-.424.07-.764.383-.929.78-.165.398-.143.854.107 1.204l.527.738c.32.447.269 1.06-.12 1.45l-.774.773a1.125 1.125 0 0 1-1.449.12l-.738-.527c-.35-.25-.806-.272-1.203-.107-.398.165-.71.505-.781.929l-.149.894c-.09.542-.56.94-1.11.94h-1.094c-.55 0-1.019-.398-1.11-.94l-.148-.894c-.071-.424-.384-.764-.781-.93-.398-.164-.854-.142-1.204.108l-.738.527c-.447.32-1.06.269-1.45-.12l-.773-.774a1.125 1.125 0 0 1-.12-1.45l.527-.737c.25-.35.272-.806.108-1.204-.165-.397-.506-.71-.93-.78l-.894-.15c-.542-.09-.94-.56-.94-1.109v-1.094c0-.55.398-1.02.94-1.11l.894-.149c.424-.07.765-.383.93-.78.165-.398.143-.854-.108-1.204l-.526-.738a1.125 1.125 0 0 1 .12-1.45l.773-.773a1.125 1.125 0 0 1 1.45-.12l.737.527c.35.25.807.272 1.204.107.397-.165.71-.505.78-.929l.15-.894Z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                </svg>
                <span class="align-middle">Settings</span>
              </a>
            </div>
          </nav>
        </div>
      </div>
    </div>
  </div>
  {/if}

  <!-- Static sidebar for desktop, closeable on small screens -->
  {#if $simStatus}
  <div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
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
                <a href="/" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
                  </svg>
                  <span class="align-middle">Dashboard</span>
                </a>
              </li>
              <li>
                <a href="/flights" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/flights' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 21l16.5-9-16.5-9v7.5l13.5 1.5-13.5 1.5V21z" />
                  </svg>
                  <span class="align-middle">My Flights</span>
                </a>
              </li>
              <li>
                <a href="/logs" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/logs' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z" />
                  </svg>
                  <span class="align-middle">Logs</span>
                </a>
              </li>
            </ul>
          </li>
        </ul>
        <div class="mt-auto">
          <a href="/settings" class={`group flex items-center gap-x-3 rounded-md p-2 text-sm font-semibold ${$page.url.pathname === '/settings' ? 'bg-gray-800 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'}`} onclick={() => sidebarOpen = false}>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 align-middle">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93.398.164.855.142 1.205-.108l.737-.527a1.125 1.125 0 0 1 1.45.12l.773.774c.39.389.44 1.002.12 1.45l-.527.737c-.25.35-.272.806-.107 1.204.165.397.505.71.93.78l.893.15c.543.09.94.559.94 1.109v1.094c0 .55-.397 1.02-.94 1.11l-.894.149c-.424.07-.764.383-.929.78-.165.398-.143.854.107 1.204l.527.738c.32.447.269 1.06-.12 1.45l-.774.773a1.125 1.125 0 0 1-1.449.12l-.738-.527c-.35-.25-.806-.272-1.203-.107-.398.165-.71.505-.781.929l-.149.894c-.09.542-.56.94-1.11.94h-1.094c-.55 0-1.019-.398-1.11-.94l-.148-.894c-.071-.424-.384-.764-.781-.93-.398-.164-.854-.142-1.204.108l-.738.527c-.447.32-1.06.269-1.45-.12l-.773-.774a1.125 1.125 0 0 1-.12-1.45l.527-.737c.25-.35.272-.806.108-1.204-.165-.397-.506-.71-.93-.78l-.894-.15c-.542-.09-.94-.56-.94-1.109v-1.094c0-.55.398-1.02.94-1.11l.894-.149c.424-.07.765-.383.93-.78.165-.398.143-.854-.108-1.204l-.526-.738a1.125 1.125 0 0 1 .12-1.45l.773-.773a1.125 1.125 0 0 1 1.45-.12l.737.527c.35.25.807.272 1.204.107.397-.165.71-.505.78-.929l.15-.894Z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
            </svg>
            <span class="align-middle">Settings</span>
          </a>
        </div>
      </nav>
    </div>
  </div>
  {/if}

  <div class={$simStatus ? "lg:pl-72" : ""}>
    {#if $simStatus}
      <div class="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-xs sm:gap-x-6 sm:px-6 lg:px-8">
        <button type="button" class="-m-2.5 p-2.5 text-gray-700 lg:hidden" onclick={() => sidebarOpen = true}>
          <span class="sr-only">Open sidebar</span>
          <svg class="size-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
          </svg>
        </button>
        <div class="flex flex-1 items-center justify-end gap-4">
          {#if $simulatorState}
            <span class="ml-2 flex items-center">
              {#if $simulatorState.pause === 1}
                <svg
                  class="size-6 text-yellow-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                >
                  <title>
                    Paused{typeof $simulatorState.realism === 'number' ? `\nRealism: ${($simulatorState.realism * 100).toFixed(0)}%` : ''}{typeof $simulatorState.simulation_rate === 'number' ? `\nRate: ${$simulatorState.simulation_rate}` : ''}
                  </title>
                  <rect x="6" y="4" width="3" height="16" rx="1" fill="currentColor" />
                  <rect x="15" y="4" width="3" height="16" rx="1" fill="currentColor" />
                </svg>
              {:else}
                <svg
                  class="size-6 text-emerald-500"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                >
                  <title>
                    Running{typeof $simulatorState.realism === 'number' ? `\nRealism: ${($simulatorState.realism * 100).toFixed(0)}%` : ''}{typeof $simulatorState.simulation_rate === 'number' ? `\nRate: ${$simulatorState.simulation_rate}` : ''}
                  </title>
                  <polygon points="6,4 20,12 6,20" fill="currentColor" />
                </svg>
              {/if}
            </span>
          {/if}
        </div>
      </div>
    {/if}
    <main class="h-min-full">
      <div>
        {@render children()}
      </div>
    </main>
  </div>
</div>