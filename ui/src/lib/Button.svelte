<script lang="ts">
  import { onMount } from 'svelte';

  export let loading: boolean = false;
  let step: number = 0;

  let timerId = -1;

  const LoadingChars = ['-', '\\', '|', '/'];

  export let onClick: () => undefined | Promise<unknown> = () => undefined;

  async function handleClick() {
    if (loading) {
      return;
    }
    loading = true;
    try {
      await onClick();
    } finally {
      loading = false;
    }
  }

  function startLoadingText() {
    clearInterval(timerId);
    timerId = setInterval(() => {
      step = (step + 1) % LoadingChars.length;
    }, 100);
  }

  function stopLoadingText() {
    clearInterval(timerId);
  }

  $: {
    if (loading) {
      startLoadingText();
    } else {
      stopLoadingText();
    }
  }

  onMount(() => {
    return () => {
      stopLoadingText();
    };
  });

  console.log($$props);
</script>

<style lang="scss">
  .button {
    position: relative;
    .loadingText {
      width: 100%;
      height: 100%;
      overflow: hidden;
      white-space: nowrap;
      position: absolute;
      top: 0;
      left: 0;
      font-family: monospace;
    }
  }
</style>

<button {...$$props} class:button={true} disabled={loading || $$props.disabled} on:click={handleClick}>
  <span class="loadingText" style:opacity={loading ? 1 : 0}>
    {LoadingChars[step]}
  </span>
  <span style:opacity={loading ? 0 : 1}>
    <slot></slot>
  </span>
</button>
