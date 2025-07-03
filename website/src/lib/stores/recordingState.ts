import { writable } from 'svelte/store';

export type RecordingState = 'idle' | 'recording' | 'stopping';

export const recordingState = writable<RecordingState>('idle');

export function startRecording() {
  recordingState.set('recording');
  // TODO: Call backend to start recording
}

export function stopRecording() {
  recordingState.set('stopping');
  // TODO: Call backend to stop recording
  // Simulate stop after delay for UI feedback
  setTimeout(() => recordingState.set('idle'), 1000);
}
