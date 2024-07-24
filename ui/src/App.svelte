<script lang="ts">
	import type { IDatasource } from './api/datasource';
	import Button from './lib/Button.svelte';
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
			<Button onClick={() => new Promise(r => setTimeout(r, 3000))}>Tasks</Button>
		</div>
	</div>
	<div class="table">
		<FileExplorer bind:datasource={datasource} />
	</div>
</div>


