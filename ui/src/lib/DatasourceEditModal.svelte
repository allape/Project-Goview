<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { type IDatasource, save, Types } from '../api/datasource';
  import { stringify } from '../util/error';
  import Modal from './Modal.svelte';

  const dispatch = createEventDispatcher();

  export let open: boolean = false;
  export let ds: IDatasource = {} as IDatasource;

  async function handleSubmit() {
    try {
      await save(ds);
      dispatch('done');
      open = false;
    } catch(e) {
      alert(stringify(e));
    }
  }
</script>

<style lang="scss">
  @import "../style/common";
  .form {
    @include ThemedBackground;
    & {
      max-width: 600px;
      padding: 0 40px 20px 20px;
      margin: 0 auto;
    }
    select {
      width: 100%;
    }
  }
</style>

<Modal bind:open={open}>
  <form class="form">
    <h2>Edit Datasource</h2>
    <table>
      <tr>
        <td><label for="Name">Name*:</label></td>
        <td><input id="Name" name="name" type="text" bind:value={ds.name} /></td>
      </tr>
      <tr>
        <td><label for="Type">Type*:</label></td>
        <td>
          <select id="Type" name="type" bind:value={ds.type}>
            {#each Types as type (type.value)}
              <option value={type.value}>{type.label}</option>
            {/each}
          </select>
        </td>
      </tr>
      <tr>
        <td><label for="CWD">CWD*:</label></td>
        <td><input id="CWD" name="cwd" type="text" bind:value={ds.cwd} /></td>
      </tr>
    </table>
    <div>
      <button type="button" on:click={() => open = false}>Close</button>
      <button type="button" on:click={handleSubmit}>Save</button>
    </div>
  </form>
</Modal>
