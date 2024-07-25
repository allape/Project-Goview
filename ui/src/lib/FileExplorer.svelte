<script lang="ts">
  import { type IDatasource, type IFile, ls, stat } from '../api/datasource';
  import { gen } from '../api/preview';
  import CDDotDot from '../asset/i_v_cd...jpg';
  import Folder from '../asset/i_v_folder.jpg';
  import Button from './Button.svelte';

  const SEP = '/';

  let loading: boolean = false;

  export let datasource: IDatasource | undefined;
  export let cwd: string = '';

  export let files: IFile[] = [];

  async function render() {
    if (!datasource || loading) {
      return;
    }

    loading = true;

    files = [];

    try {
      const s = await stat(datasource.id, cwd);
      if (!s.isDir) {
        files = [s];
      } else {
        files = await ls(datasource.id, cwd);
        files.sort((a, b) => {
          if (a.isDir && !b.isDir) {
            return -1;
          } else if (!a.isDir && b.isDir) {
            return 1;
          } else {
            return a.name.localeCompare(b.name);
          }
        });
      }
    } finally {
      loading = false;
    }
  }

  function onClick(file: IFile) {
    if (file.isDir) {
      cwd += file.name + SEP;
      render();
    } else {
      // TODO open file?
    }
  }

  async function genPreview(file: IFile) {
    if (file.isDir) {
      return;
    }
    await gen(datasource!.id, cwd + file.name);
    alert('Done Gen Preview');
    render().then();
  }

  function onCdDotDot() {
    const parts = cwd.split(SEP);
    parts.pop();
    parts.pop();
    cwd = parts.join(SEP);
    if (cwd) {
      cwd += SEP;
    }
    render();
  }

  $: {
    if (datasource) {
      render();
    } else {
      cwd = '';
    }
  }
</script>

<style lang="scss">
  @import "../style/common";

  .wrapper {
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
    height: 100%;
    overflow: hidden;

    .nav {
      display: flex;
      padding: 0 0 10px 0;
      flex-wrap: nowrap;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      gap: 10px;

      .cwd {
        display: flex;
        flex: 1;

        input {
          flex: 1;
        }
      }

      .buttons {
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }

    .files {
      flex: 1;
      display: flex;
      justify-content: flex-start;
      align-items: stretch;
      gap: 10px;
      flex-wrap: wrap;
      overflow-y: auto;

      .empty {
        width: 100%;
        padding: 40px;
        font-size: 30px;
        font-weight: bold;
        text-align: center;
      }

      .file {
        width: 180px;
        display: flex;
        flex-direction: column;
        border: 1px solid lightgray;
        position: relative;

        &.folder {
          cursor: pointer;
        }

        &.placeholder {
          .name {
            text-align: center;
            user-select: none;
          }
        }

        &:not(.placeholder):hover {
          .name {
            opacity: 0;
          }

          .fullname {
            display: block;
          }
        }

        .preview {
          width: 100%;
          flex: 1;
          font-size: 0;

          img {
            width: 100%;
            height: 100%;
            object-fit: contain;
          }
        }

        .name {
          border-top: 1px solid lightgray;
          padding: 3px 5px;
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
        }

        .fullname {
          @include ThemedBackground;

          & {
            padding: 3px 5px;
            position: absolute;
            bottom: 0;
            left: 0;
            width: 100%;
            white-space: pre-wrap;
            display: none;
          }
        }

        .buttons {
          position: absolute;
          right: 0;
          top: 0;
          padding: 10px;
          transition: 250ms;
        }
      }
    }
  }
</style>

<div class="wrapper">
  <div class="nav">
    <div class="cwd">
      <input type="text" placeholder="cwd" bind:value={cwd}>
    </div>
    <div class="buttons">
      <Button onClick={render}>Refresh</Button>
    </div>
  </div>
  <div class="files">
    {#if loading}
      Loading...
    {:else if cwd}
      <div class="file folder placeholder" on:click={onCdDotDot} role="none">
        <div class="preview">
          <img src={CDDotDot} alt="cd ..">
        </div>
        <div class="name">cd ..</div>
      </div>
    {:else}
      {#if !files.length}
        <div class="empty">EMPTY</div>
      {/if}
    {/if}

    {#each files as file (file.name)}
      <div class="file" data-path={file._path} class:folder={file.isDir} on:click={() => onClick(file)} role="none">
        <div class="preview">
          {#if file.isDir}
            <img src={Folder} alt={file.name}>
          {:else}
            <img src={file._preview} alt={file.name}>
          {/if}
        </div>
        <div class="name">{file._displayName}</div>
        <div class="fullname">
          {file._displayName}
        </div>
        <div class="buttons">
          {#if !file.isDir}
            <Button onClick={() => genPreview(file)}>Gen Preview</Button>
            <button on:click={() => window.open(file._preview)}>Open</button>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>
