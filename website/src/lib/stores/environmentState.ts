import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetEnvironmentState } from '$lib/wailsjs/go/internal/App';

export interface EnvironmentState {
  ambient_temperature: number;
  ambient_visibility: number;
  ambient_wind_direction: number;
  ambient_wind_velocity: number;
  local_day: number;
  local_day_of_week: number;
  local_month: number;
  local_time: number;
  local_year: number;
  sea_level_pressure: number;
  sim_time: number;
  zulu_day: number;
  zulu_day_of_week: number;
  zulu_month: number;
  zulu_time: number;
  zulu_year: number;
  // New simvars
  time_zone_offset: number;
  zulu_sunrise_time: number;
  zulu_sunset_time: number;
  time_of_day: number; // 0=dawn, 1=day, 2=dusk, 3=night
}

export const environmentState = writable<EnvironmentState>();

// Initialize with backend state
GetEnvironmentState().then(environmentState.set);

EventsOn('environment::state', (state: EnvironmentState) => {
  environmentState.set(state);
});
