import { writable } from 'svelte/store';
import type { FileKey } from '../api/datasource';

export const GenerationQueue = writable<FileKey[]>([]);
