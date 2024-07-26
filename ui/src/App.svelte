<script lang="ts">
	import { onMount } from 'svelte';
	import type { IDatasource } from './api/datasource';
	import type { IR } from './api/http';
	import { BASE_URL } from './config/server';
	import { GenerationQueue } from './context/generation-queue';
	import DatasourceModalButton from './lib/DatasourceModalButton.svelte';
	import FileExplorer from './lib/FileExplorer.svelte';

	let datasources: IDatasource[] = [];

	let selectedDatasourceID: IDatasource['id'] | undefined;
	let datasource: IDatasource | undefined;

	function handleDatasourcesChanged(e: CustomEvent<IDatasource[]>) {
		datasources = e.detail;
	}

	$: {
		if (selectedDatasourceID) {
			datasource = datasources.find(d => d.id === selectedDatasourceID) || undefined;
		} else {
			datasource = undefined;
		}
	}

	onMount(() => {
		const sse = new EventSource(`${BASE_URL}/preview/task/count/sse`);
		sse.addEventListener('EVENT_PREVIEW_TASK_COUNT', e => {
			const data: IR<string> = JSON.parse(e.data);
			try {
				GenerationQueue.set(JSON.parse(data.d) || []);
			} catch (e) {
				console.error('[SSE] EVENT_PREVIEW_TASK_COUNT:', e);
			}
		});
		sse.addEventListener('error', e => {
			console.error('[SSE] error:', e);
		});
		sse.addEventListener('open', e => {
			console.log('[SSE] open:', e);
		});
		return () => {
			sse.close();
		};
	});
</script>

<style lang="scss">
  .wrapper {
    max-width: 1200px;
    margin: auto;
    padding: 10px;
    height: calc(100% - 20px);
    display: flex;
    flex-direction: column;
    justify-content: stretch;
    align-items: stretch;

    .buttons {
      padding: 0 0 10px 0;
      display: flex;
      justify-content: space-between;
			gap: 10px;
			.full {
				flex: 1;
				select {
					width: 100%;
				}
			}
			.taskCount {
				min-width: 150px;
				text-align: left;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}
    }

		.table {
			flex: 1;
			overflow-y: auto;
			overflow-x: hidden;
			border-top: 1px solid lightgray;
			padding: 10px 0;
		}
  }
</style>

<div class="wrapper">
	<div class="buttons">
		<div>
			<div class="taskCount">Remained: {$GenerationQueue.length.toLocaleString()}</div>
		</div>
		<div class="full">
			<select bind:value={selectedDatasourceID}>
				<option>-</option>
				{#each datasources as datasource}
					<option value={datasource.id}>{datasource.name}</option>
				{/each}
			</select>
		</div>
		<div>
			<DatasourceModalButton on:change={handleDatasourcesChanged} />
		</div>
	</div>
	<div class="table">
		<FileExplorer bind:datasource={datasource} />
	</div>
</div>


