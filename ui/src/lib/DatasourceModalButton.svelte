<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import { getAll, type IDatasource } from '../api/datasource';
  import DatasourceEditModal from './DatasourceEditModal.svelte';
  import Modal from './Modal.svelte';

  const dispatch = createEventDispatcher();

  let open: boolean = false;
  let datasources: IDatasource[] = [];

  let saveOpen: boolean = false;
  let saveModel: IDatasource = {} as IDatasource;

  async function getList() {
    datasources = await getAll();
    dispatch('change', datasources);
  }

  onMount(() => {
    handleChange();
  });

  function handleChange() {
    getList().then();
  }

  function saveDatasource(ds?: IDatasource) {
    saveModel = ds || ({} as IDatasource);
    saveOpen = true;
  }
</script>

<style lang="scss">
  @import "../style/common";

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
    }

    .tableWrapper {
      flex: 1;
      width: 100%;
      overflow: auto;

      .table {
        @include ThemedBackground;

        & {
          width: 100%;
          border-collapse: collapse;
        }

        td, th {
          border: 1px solid lightgray;
        }

        td {
          padding: 0 10px;
        }
      }
    }
  }
</style>

<button on:click={() => open = true}>Manage Datasources</button>
<Modal bind:open>
  <div class="wrapper">
    <div class="buttons">
      <div>
        <button on:click={handleChange}>Reload</button>
        <button on:click={() => saveDatasource()}>Add</button>
      </div>
      <div>
        <button on:click={() => open = false}>Close</button>
      </div>
    </div>
    <div class="tableWrapper">
      <table class="table">
        <thead>
        <tr>
          <th>ID</th>
          <th>Created At</th>
          <th>Name</th>
          <th>Type</th>
          <th>CWD</th>
          <th>Ops.</th>
        </tr>
        </thead>
        <tbody>
        {#each datasources as ds (ds.id)}
          <tr>
            <td>{ds.id}</td>
            <td>{new Date(ds.createdAt).toLocaleString()}</td>
            <td>{ds.name}</td>
            <td>{ds.type}</td>
            <td>{ds.cwd}</td>
            <td>
              <button on:click={() => saveDatasource(ds)}>Edit</button>
            </td>
          </tr>
        {/each}
        </tbody>
      </table>
    </div>
  </div>
</Modal>

<DatasourceEditModal bind:open={saveOpen} bind:ds={saveModel} on:done={handleChange} />
