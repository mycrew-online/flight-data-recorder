import { writable } from 'svelte/store';

import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetSimStatus } from '$lib/wailsjs/go/internal/App';

export const simStatus = writable(false);

// Initialize with backend status
GetSimStatus().then(simStatus.set);

EventsOn('global::sim-status', (status: boolean) => {
  simStatus.set(status);
});
